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
	"sync"
	"sync/atomic"
	"time"
)

var (
	specsDir        = flag.String("specs", "docs/specifications/api", "OpenAPI specs directory")
	ollamaURL       = flag.String("ollama-url", "http://localhost:11434", "Ollama API URL")
	model           = flag.String("model", "deepseek-r1:1.5b", "LLM model to use")
	outputFile      = flag.String("output", "pkg/types/descriptions_generated.json", "Output JSON file")
	timeout         = flag.Duration("timeout", 120*time.Second, "Per-request timeout")
	verbose         = flag.Bool("v", false, "Verbose output")
	dryRun          = flag.Bool("dry-run", false, "Print what would be done without calling LLM")
	workers         = flag.Int("workers", 8, "Number of parallel workers (default: 8)")
	maxRetries      = flag.Int("max-retries", 3, "Maximum retries per request on timeout")
	failThreshold   = flag.Float64("fail-threshold", 0.2, "Fail if error rate exceeds this (0.0-1.0)")
	ciMode          = flag.Bool("ci", false, "CI mode: use GitHub Actions annotations and fail-fast on errors")
	failFast        = flag.Bool("fail-fast", false, "Exit immediately on first error (auto-enabled in CI mode)")
	maxConsecErrors = flag.Int("max-consec-errors", 3, "Max consecutive errors before aborting (0=disabled)")
)

// DescriptionOutput is the JSON output format
type DescriptionOutput struct {
	Generated    string            `json:"generated"`
	Model        string            `json:"model"`
	Workers      int               `json:"workers"`
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

// specJob represents a work item for the worker pool
type specJob struct {
	path         string
	resourceName string
}

// specResult represents the result of processing a spec
type specResult struct {
	resourceName string
	description  string
	err          error
}

func main() {
	flag.Parse()

	// Validate workers
	if *workers < 1 {
		*workers = 1
	}
	if *workers > 32 {
		fmt.Fprintf(os.Stderr, "Warning: %d workers is very high, capping at 32\n", *workers)
		*workers = 32
	}

	output := DescriptionOutput{
		Generated:    time.Now().UTC().Format(time.RFC3339),
		Model:        *model,
		Workers:      *workers,
		Descriptions: make(map[string]string),
	}

	// Check Ollama availability (unless dry-run)
	if !*dryRun {
		if !checkOllamaAvailable() {
			logError("Ollama not available at %s", *ollamaURL)
			logError("Start Ollama with: ollama serve")
			logError("Or ensure Ollama is running before executing this script.")
			os.Exit(1)
		}

		if !checkModelAvailable(*model) {
			logError("Model %s not found", *model)
			logError("Install with: ollama pull %s", *model)
			os.Exit(1)
		}
		fmt.Printf("Using Ollama at %s with model %s (%d workers, timeout: %v, max-retries: %d)\n",
			*ollamaURL, *model, *workers, *timeout, *maxRetries)

		// Warm up the model to ensure it's loaded into memory before workers start
		if err := warmupModel(); err != nil {
			logError("Failed to warm up model: %v", err)
			logError("The model may not be fully downloaded or Ollama may be having issues.")
			os.Exit(1)
		}
	}

	// Find all OpenAPI specs
	pattern := filepath.Join(*specsDir, "*.json")
	specs, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find specs: %v\n", err)
		os.Exit(1)
	}

	// Filter to valid resource specs
	var jobs []specJob
	for _, specPath := range specs {
		resourceName := extractResourceName(specPath)
		if resourceName == "" {
			continue
		}
		if !isResourceSpec(specPath) {
			continue
		}
		jobs = append(jobs, specJob{path: specPath, resourceName: resourceName})
	}

	fmt.Printf("Found %d OpenAPI specs, %d are resources to process\n", len(specs), len(jobs))

	// Enable fail-fast in CI mode by default
	effectiveFailFast := *failFast || *ciMode

	// Process specs with worker pool
	var (
		processed       int32
		consecutiveErrs int32
		aborted         int32
		errorsMu        sync.Mutex
		descMu          sync.Mutex
		wg              sync.WaitGroup
	)

	// Create job channel
	jobChan := make(chan specJob, len(jobs))

	// Create result channel
	resultChan := make(chan specResult, len(jobs))

	// Start workers
	for w := 0; w < *workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobChan {
				// Check if we should abort
				if atomic.LoadInt32(&aborted) == 1 {
					return
				}

				if *verbose {
					fmt.Printf("[Worker %d] Processing %s...\n", workerID, job.resourceName)
				}

				desc, err := processSpec(job.path, job.resourceName)
				resultChan <- specResult{
					resourceName: job.resourceName,
					description:  desc,
					err:          err,
				}
			}
		}(w)
	}

	// Send jobs to workers
	go func() {
		for _, job := range jobs {
			// Check if we should abort before sending more jobs
			if atomic.LoadInt32(&aborted) == 1 {
				break
			}
			jobChan <- job
		}
		close(jobChan)
	}()

	// Wait for workers and close result channel
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results with fail-fast support
	for result := range resultChan {
		if result.err != nil {
			// Always log errors with full context for troubleshooting
			logError("[%s]: %v", result.resourceName, result.err)
			errorsMu.Lock()
			output.Errors = append(output.Errors, fmt.Sprintf("%s: %v", result.resourceName, result.err))
			currentErrors := len(output.Errors)
			errorsMu.Unlock()

			// Track consecutive errors
			atomic.AddInt32(&consecutiveErrs, 1)

			// Fail fast: exit immediately on first error
			if effectiveFailFast {
				logError("Aborting due to fail-fast mode (first error encountered)")
				logError("Resource '%s' failed after %d retries", result.resourceName, *maxRetries)
				logError("Check Ollama server status: curl %s/api/tags", *ollamaURL)
				atomic.StoreInt32(&aborted, 1)
				// Drain remaining results
				go func() {
					for range resultChan {
					}
				}()
				os.Exit(1)
			}

			// Check consecutive error threshold
			if *maxConsecErrors > 0 && atomic.LoadInt32(&consecutiveErrs) >= int32(*maxConsecErrors) {
				logError("Aborting: %d consecutive errors exceeded threshold", *maxConsecErrors)
				logError("Total errors so far: %d", currentErrors)
				logError("This usually indicates Ollama is not responding correctly")
				logError("Check: curl -X POST %s/api/generate -d '{\"model\":\"%s\",\"prompt\":\"test\"}'", *ollamaURL, *model)
				atomic.StoreInt32(&aborted, 1)
				// Drain remaining results
				go func() {
					for range resultChan {
					}
				}()
				os.Exit(1)
			}
			continue
		}

		// Reset consecutive error counter on success
		atomic.StoreInt32(&consecutiveErrs, 0)

		if result.description != "" {
			descMu.Lock()
			output.Descriptions[result.resourceName] = result.description
			descMu.Unlock()
			atomic.AddInt32(&processed, 1)
			if *verbose {
				fmt.Printf("  OK [%s]: %s\n", result.resourceName, truncate(result.description, 60))
			}
		}
	}

	fmt.Printf("\nProcessed %d resources, %d errors (using %d workers)\n", processed, len(output.Errors), *workers)

	// Write JSON output
	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		logError("Failed to marshal JSON: %v", err)
		os.Exit(1)
	}

	if err := os.WriteFile(*outputFile, jsonData, 0644); err != nil {
		logError("Failed to write output: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Output written to %s\n", *outputFile)

	// Also output to stdout for CI visibility
	if *verbose {
		fmt.Println("\n--- Generated Descriptions ---")
		fmt.Println(string(jsonData))
	}

	// Validate results and exit with appropriate code
	exitCode := 0

	// Check if no resources were processed (when not in dry-run mode)
	if !*dryRun && processed == 0 {
		logError("No resources were processed. Expected to process %d resources.", len(jobs))
		logError("This usually indicates Ollama is not responding or the specs directory is wrong.")
		exitCode = 1
	}

	// Check error rate against threshold
	if len(jobs) > 0 {
		errorRate := float64(len(output.Errors)) / float64(len(jobs))
		if errorRate > *failThreshold {
			logError("Error rate %.1f%% exceeds threshold %.1f%% (%d/%d resources failed)",
				errorRate*100, *failThreshold*100, len(output.Errors), len(jobs))
			for _, e := range output.Errors {
				logError("  - %s", e)
			}
			exitCode = 1
		} else if len(output.Errors) > 0 {
			// Warn about errors even if under threshold
			logWarning("%d resources had errors (%.1f%% error rate, threshold: %.1f%%)",
				len(output.Errors), errorRate*100, *failThreshold*100)
			for _, e := range output.Errors {
				logWarning("  - %s", e)
			}
		}
	}

	if exitCode != 0 {
		os.Exit(exitCode)
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

// callOllama sends a request to the Ollama API with retry logic for timeouts
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

	var lastErr error
	for attempt := 1; attempt <= *maxRetries; attempt++ {
		client := &http.Client{Timeout: *timeout}
		resp, err := client.Post(*ollamaURL+"/api/generate", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			lastErr = fmt.Errorf("ollama request failed: %w", err)
			if isTimeoutError(err) && attempt < *maxRetries {
				if *verbose {
					fmt.Printf("    Timeout on attempt %d/%d, retrying...\n", attempt, *maxRetries)
				}
				// Exponential backoff: 2s, 4s, 8s...
				time.Sleep(time.Duration(1<<uint(attempt)) * time.Second)
				continue
			}
			return "", lastErr
		}

		// Read response body (don't use defer in loop)
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			if attempt < *maxRetries {
				time.Sleep(time.Duration(1<<uint(attempt)) * time.Second)
				continue
			}
			return "", lastErr
		}

		if resp.StatusCode != 200 {
			return "", fmt.Errorf("ollama returned %d: %s", resp.StatusCode, string(body))
		}

		var result OllamaResponse
		if err := json.Unmarshal(body, &result); err != nil {
			return "", fmt.Errorf("failed to decode response: %w", err)
		}

		// Log raw response in verbose mode for debugging
		if *verbose && len(result.Response) > 0 {
			rawPreview := result.Response
			if len(rawPreview) > 100 {
				rawPreview = rawPreview[:100] + "..."
			}
			fmt.Printf("    Raw response (%d chars): %s\n", len(result.Response), strings.ReplaceAll(rawPreview, "\n", "\\n"))
		}

		// Clean up the response
		cleaned := cleanResponse(result.Response)

		// Treat empty cleaned responses as errors (retry)
		if cleaned == "" || cleaned == "." {
			lastErr = fmt.Errorf("ollama returned empty or invalid response (raw length: %d)", len(result.Response))
			if attempt < *maxRetries {
				if *verbose {
					fmt.Printf("    Empty response on attempt %d/%d, retrying...\n", attempt, *maxRetries)
				}
				time.Sleep(time.Duration(1<<uint(attempt)) * time.Second)
				continue
			}
			return "", lastErr
		}

		return cleaned, nil
	}

	return "", lastErr
}

// cleanResponse removes common LLM artifacts from the response
func cleanResponse(response string) string {
	result := strings.TrimSpace(response)

	// Remove deepseek-r1 <think>...</think> tags (reasoning models output these)
	// The pattern can be <think>content</think> followed by the actual answer
	thinkStart := strings.Index(result, "<think>")
	thinkEnd := strings.Index(result, "</think>")
	if thinkStart != -1 && thinkEnd != -1 && thinkEnd > thinkStart {
		// Extract content after </think>
		afterThink := strings.TrimSpace(result[thinkEnd+len("</think>"):])
		if afterThink != "" {
			result = afterThink
		}
		// If nothing after think tags, check if there's content before
		if afterThink == "" && thinkStart > 0 {
			result = strings.TrimSpace(result[:thinkStart])
		}
	}

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

// warmupModel sends a simple request to ensure the model is loaded into memory
// and validates that it actually produces output
func warmupModel() error {
	fmt.Println("Warming up model (loading into memory)...")
	start := time.Now()

	// Use a simple prompt that should always produce output
	reqBody := OllamaRequest{
		Model:  *model,
		Prompt: "Respond with exactly one word: ready",
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.0,
			"num_predict": 50, // Allow enough tokens for thinking models
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal warmup request: %w", err)
	}

	// Use longer timeout for warmup (model loading can take time)
	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Post(*ollamaURL+"/api/generate", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("warmup request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read warmup response: %w", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("warmup returned %d: %s", resp.StatusCode, string(body))
	}

	// Parse response to validate model is actually generating output
	var result OllamaResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse warmup response: %w (body: %s)", err, string(body))
	}

	if result.Response == "" {
		return fmt.Errorf("warmup returned empty response - model may not be working correctly (body: %s)", string(body))
	}

	duration := time.Since(start).Round(time.Second)
	responseLen := len(result.Response)
	responsePreview := result.Response
	if len(responsePreview) > 100 {
		responsePreview = responsePreview[:100] + "..."
	}
	fmt.Printf("✅ Model warmed up in %v (response: %d chars)\n", duration, responseLen)
	if *verbose {
		fmt.Printf("   Warmup response: %s\n", strings.ReplaceAll(responsePreview, "\n", "\\n"))
	}
	return nil
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

// logError logs an error message, using GitHub Actions annotation format in CI mode
func logError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if *ciMode {
		fmt.Fprintf(os.Stderr, "::error::%s\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	}
}

// logWarning logs a warning message, using GitHub Actions annotation format in CI mode
func logWarning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if *ciMode {
		fmt.Fprintf(os.Stderr, "::warning::%s\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "Warning: %s\n", msg)
	}
}

// isTimeoutError checks if an error is a timeout error
func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "context deadline exceeded") ||
		strings.Contains(errStr, "Client.Timeout") ||
		strings.Contains(errStr, "timeout")
}
