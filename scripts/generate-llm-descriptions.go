//go:build ignore

// generate-llm-descriptions.go generates grammatically correct descriptions using a local LLM.
// Run with: go run scripts/generate-llm-descriptions.go
// Requires: Ollama running with deepseek-r1:1.5b model
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	specsDir   = flag.String("specs", "docs/specifications/api", "OpenAPI specs directory")
	ollamaURL  = flag.String("ollama-url", "http://localhost:11434", "Ollama API URL")
	model      = flag.String("model", "deepseek-r1:1.5b", "LLM model to use")
	outputFile = flag.String("output", "pkg/types/descriptions_generated.json", "Output JSON file")
	timeout    = flag.Duration("timeout", 60*time.Second, "Per-request timeout")
	verbose    = flag.Bool("v", false, "Verbose output")
	dryRun     = flag.Bool("dry-run", false, "Print what would be done without calling LLM")
)

// DescriptionOutput is the JSON output format
type DescriptionOutput struct {
	Generated    string            `json:"generated"`
	Model        string            `json:"model"`
	Descriptions map[string]string `json:"descriptions"`
	Errors       []string          `json:"errors,omitempty"`
}

// OllamaRequest is the Ollama API request format
type OllamaRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream"`
	Options map[string]interface{} `json:"options"`
}

// OllamaResponse is the Ollama API response format
type OllamaResponse struct {
	Response string `json:"response"`
}

// RawSpec is used to parse the OpenAPI spec with extensions
type RawSpec struct {
	Info  RawInfo                `json:"info"`
	Paths map[string]interface{} `json:"paths"`
}

// RawInfo contains the info section
type RawInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func main() {
	flag.Parse()

	output := DescriptionOutput{
		Generated:    time.Now().UTC().Format(time.RFC3339),
		Model:        *model,
		Descriptions: make(map[string]string),
	}

	// Check Ollama availability (unless dry-run)
	if !*dryRun {
		if !checkOllamaAvailable() {
			fmt.Fprintf(os.Stderr, "Error: Ollama not available at %s\n", *ollamaURL)
			fmt.Fprintf(os.Stderr, "Start Ollama with: ollama serve\n")
			os.Exit(1)
		}

		if !checkModelAvailable(*model) {
			fmt.Fprintf(os.Stderr, "Error: Model %s not found\n", *model)
			fmt.Fprintf(os.Stderr, "Install with: ollama pull %s\n", *model)
			os.Exit(1)
		}
		fmt.Printf("Using Ollama at %s with model %s\n", *ollamaURL, *model)
	}

	// Find all OpenAPI specs
	pattern := filepath.Join(*specsDir, "*.json")
	specs, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find specs: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d OpenAPI specs\n", len(specs))

	// Process each spec
	processed := 0
	for _, specPath := range specs {
		resourceName := extractResourceName(specPath)
		if resourceName == "" {
			continue
		}

		// Skip non-resource specs (e.g., common types)
		if !isResourceSpec(specPath) {
			continue
		}

		if *verbose {
			fmt.Printf("Processing %s...\n", resourceName)
		}

		desc, err := processSpec(specPath, resourceName)
		if err != nil {
			if *verbose {
				fmt.Printf("  Error: %v\n", err)
			}
			output.Errors = append(output.Errors, fmt.Sprintf("%s: %v", resourceName, err))
			continue
		}

		if desc != "" {
			output.Descriptions[resourceName] = desc
			processed++
			if *verbose {
				fmt.Printf("  OK: %s\n", truncate(desc, 60))
			}
		}
	}

	fmt.Printf("\nProcessed %d resources, %d errors\n", processed, len(output.Errors))

	// Write JSON output
	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal JSON: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(*outputFile, jsonData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write output: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Output written to %s\n", *outputFile)

	// Also output to stdout for CI visibility
	if *verbose {
		fmt.Println("\n--- Generated Descriptions ---")
		fmt.Println(string(jsonData))
	}
}

// extractResourceName extracts the resource name from the spec filename
func extractResourceName(path string) string {
	base := filepath.Base(path)
	// docs-cloud-f5-com.0073.public.ves.io.schema.views.http_loadbalancer.ves-swagger.json
	// -> http_loadbalancer

	// Remove prefix and suffix
	name := strings.TrimSuffix(base, ".ves-swagger.json")
	name = strings.TrimSuffix(name, ".json")

	// Split by dots and find the resource part
	parts := strings.Split(name, ".")
	if len(parts) < 2 {
		return ""
	}

	// The resource name is typically the last meaningful part
	// Look for patterns like "views.http_loadbalancer" or just "namespace"
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		// Skip common prefixes
		if part == "ves-swagger" || part == "ves" || part == "io" ||
			part == "schema" || part == "views" || part == "public" {
			continue
		}
		// Skip numeric prefixes
		if len(part) > 0 && part[0] >= '0' && part[0] <= '9' {
			continue
		}
		// Skip the domain prefix
		if strings.HasPrefix(part, "docs-cloud-f5-com") {
			continue
		}
		return part
	}

	return ""
}

// isResourceSpec checks if this is a resource spec (not a common type spec)
func isResourceSpec(path string) bool {
	base := filepath.Base(path)
	// Skip common schemas that aren't resources
	skipPatterns := []string{
		".types.", ".common.", ".object.", ".status.",
		"_types.", "_common.", "Custom", "custom_",
	}
	for _, pattern := range skipPatterns {
		if strings.Contains(base, pattern) {
			return false
		}
	}
	return true
}

// processSpec processes a single OpenAPI spec and generates a description
func processSpec(path, resourceName string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var spec RawSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return "", fmt.Errorf("failed to parse spec: %w", err)
	}

	// Extract x-displayname from paths
	displayName := extractDisplayName(spec.Paths)

	// Build prompt
	prompt := buildPrompt(resourceName, displayName, spec.Info.Title, spec.Info.Description)

	if *dryRun {
		fmt.Printf("\n--- %s ---\n", resourceName)
		fmt.Printf("Display Name: %s\n", displayName)
		fmt.Printf("Title: %s\n", spec.Info.Title)
		fmt.Printf("Description: %s\n", truncate(spec.Info.Description, 100))
		return spec.Info.Description, nil // Return original in dry-run
	}

	// Call Ollama
	return callOllama(prompt)
}

// extractDisplayName finds x-displayname from the paths
func extractDisplayName(paths map[string]interface{}) string {
	for _, pathData := range paths {
		pathMap, ok := pathData.(map[string]interface{})
		if !ok {
			continue
		}
		if displayName, ok := pathMap["x-displayname"].(string); ok {
			return displayName
		}
	}
	return ""
}

// buildPrompt creates the LLM prompt for description generation
func buildPrompt(resourceName, displayName, title, description string) string {
	return fmt.Sprintf(`You are a technical writer creating CLI documentation for F5 Distributed Cloud.

Resource: %s
Display Name: %s
Title: %s
Original Description: %s

Write a single, concise, grammatically correct description (1-2 sentences maximum).

Requirements:
- Fix any grammatical errors (e.g., "a required parameters" should be "the required parameters")
- Focus on what the resource does and its primary use case
- Use professional technical writing tone
- Do NOT start with "This resource" or "The resource"
- Do NOT include CRUD operations or API implementation details
- Do NOT include child objects or internal details

Output ONLY the corrected description text, nothing else. No explanations, no markdown, no quotes.`, resourceName, displayName, title, description)
}

// callOllama sends a request to the Ollama API
func callOllama(prompt string) (string, error) {
	reqBody := OllamaRequest{
		Model:  *model,
		Prompt: prompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.1, // Low for consistency
			"num_predict": 150, // Limit output length
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{Timeout: *timeout}
	resp, err := client.Post(*ollamaURL+"/api/generate", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama returned %d: %s", resp.StatusCode, string(body))
	}

	var result OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Clean up the response
	cleaned := cleanResponse(result.Response)
	return cleaned, nil
}

// cleanResponse removes common LLM artifacts from the response
func cleanResponse(response string) string {
	result := strings.TrimSpace(response)

	// Remove common wrapping patterns
	result = strings.Trim(result, "\"'`")

	// Remove markdown code blocks
	if strings.HasPrefix(result, "```") {
		lines := strings.Split(result, "\n")
		var cleanLines []string
		inCodeBlock := false
		for _, line := range lines {
			if strings.HasPrefix(line, "```") {
				inCodeBlock = !inCodeBlock
				continue
			}
			if !inCodeBlock {
				cleanLines = append(cleanLines, line)
			}
		}
		result = strings.Join(cleanLines, " ")
	}

	// Remove leading "Description:" or similar
	prefixes := []string{"Description:", "Output:", "Result:", "Answer:"}
	for _, prefix := range prefixes {
		result = strings.TrimPrefix(result, prefix)
	}

	result = strings.TrimSpace(result)

	// Ensure it ends with a period
	if len(result) > 0 && !strings.HasSuffix(result, ".") && !strings.HasSuffix(result, "!") && !strings.HasSuffix(result, "?") {
		result += "."
	}

	return result
}

// checkOllamaAvailable checks if Ollama is running
func checkOllamaAvailable() bool {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(*ollamaURL + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

// checkModelAvailable checks if the specified model is available
func checkModelAvailable(model string) bool {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(*ollamaURL + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false
	}

	// Sort for display
	var names []string
	for _, m := range result.Models {
		names = append(names, m.Name)
	}
	sort.Strings(names)

	if *verbose {
		fmt.Printf("Available models: %v\n", names)
	}

	for _, m := range result.Models {
		// Check for exact match or prefix match (e.g., "deepseek-r1:1.5b" matches "deepseek-r1")
		if m.Name == model || strings.HasPrefix(m.Name, strings.Split(model, ":")[0]) {
			return true
		}
	}
	return false
}

// truncate truncates a string to maxLen characters
func truncate(s string, maxLen int) string {
	// Replace newlines with spaces for display
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
