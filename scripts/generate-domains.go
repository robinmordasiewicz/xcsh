// Code generated from plan. Run `go run scripts/generate-domains.go` to regenerate.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

// SpecIndex represents the structure of .specs/index.json
type SpecIndex struct {
	Version        string           `json:"version"`
	Timestamp      string           `json:"timestamp"`
	Specifications []SpecIndexEntry `json:"specifications"`
}

// SpecIndexEntry represents a single domain in index.json
type SpecIndexEntry struct {
	Domain         string                 `json:"domain"`
	Title          string                 `json:"title"`
	Description    string                 `json:"description"`
	File           string                 `json:"file"`
	PathCount      int                    `json:"path_count"`
	SchemaCount    int                    `json:"schema_count"`
	Complexity     string                 `json:"complexity"`
	IsPreview      bool                   `json:"is_preview"`
	RequiresTier   string                 `json:"requires_tier"`
	DomainCategory string                 `json:"domain_category"`
	UseCases       []string               `json:"use_cases"`
	RelatedDomains []string               `json:"related_domains"`
	CLIMetadata    map[string]interface{} `json:"cli_metadata"`
}

// DomainConfig represents the structure of domain_config.yaml
type DomainConfig struct {
	Version           string                    `yaml:"version"`
	Aliases           map[string][]string       `yaml:"aliases"`
	DeprecatedDomains map[string]DeprecatedInfo `yaml:"deprecated_domains"`
	MissingMetadata   []MetadataIssue           `yaml:"missing_metadata"`
}

// DeprecatedInfo tracks deprecated domain information
type DeprecatedInfo struct {
	MapsTo          string `yaml:"maps_to"`
	Reason          string `yaml:"reason"`
	DeprecatedSince string `yaml:"deprecated_since"`
}

// MetadataIssue tracks missing metadata in upstream specs
type MetadataIssue struct {
	Domain       string `yaml:"domain"`
	MissingField string `yaml:"missing_field"`
	Reason       string `yaml:"reason"`
	GitHubIssue  string `yaml:"github_issue"`
}

// DomainInfo represents domain metadata for generated code
type DomainInfo struct {
	Name           string
	DisplayName    string
	Description    string
	Aliases        []string
	Deprecated     bool
	MapsTo         string
	AddedInVersion string
	Complexity     string
	IsPreview      bool
	RequiresTier   string
	Category       string
	UseCases       []string
	RelatedDomains []string
	CLIMetadata    map[string]interface{}
}

// GeneratedDomainRegistry holds all domain info for template rendering
type GeneratedDomainRegistry struct {
	Version string
	Domains map[string]*DomainInfo
}

func main() {
	log.Println("🏗️  Generating domains from upstream specs...")

	// Step 1: Read .specs/index.json
	specIndex, err := readSpecIndex(".specs/index.json")
	if err != nil {
		log.Fatalf("Failed to read spec index: %v", err)
	}
	log.Printf("✓ Loaded spec index v%s with %d domains", specIndex.Version, len(specIndex.Specifications))

	// Step 2: Load domain_config.yaml for overrides
	config, err := readDomainConfig(".specs/domain_config.yaml")
	if err != nil {
		log.Printf("⚠️  Could not load domain config (will use defaults): %v", err)
		config = &DomainConfig{
			Aliases:           make(map[string][]string),
			DeprecatedDomains: make(map[string]DeprecatedInfo),
		}
	}

	// Step 3: Build domain registry from specs
	registry := &GeneratedDomainRegistry{
		Version: specIndex.Version,
		Domains: make(map[string]*DomainInfo),
	}

	for _, spec := range specIndex.Specifications {
		// Skip empty domains
		if spec.PathCount == 0 && spec.SchemaCount == 0 {
			log.Printf("⊘ Skipping empty domain: %s", spec.Domain)
			continue
		}

		domainInfo := &DomainInfo{
			Name:           spec.Domain,
			DisplayName:    titleCase(spec.Domain),
			Description:    spec.Description,
			Aliases:        config.Aliases[spec.Domain],
			Deprecated:     false,
			MapsTo:         "",
			Complexity:     spec.Complexity,
			IsPreview:      spec.IsPreview,
			RequiresTier:   spec.RequiresTier,
			Category:       spec.DomainCategory,
			UseCases:       spec.UseCases,
			RelatedDomains: spec.RelatedDomains,
			CLIMetadata:    spec.CLIMetadata,
		}

		// Apply deprecated domain mappings
		if deprecated, exists := config.DeprecatedDomains[spec.Domain]; exists {
			domainInfo.Deprecated = true
			domainInfo.MapsTo = deprecated.MapsTo
		}

		registry.Domains[spec.Domain] = domainInfo
	}

	log.Printf("✓ Generated registry with %d active domains", len(registry.Domains))

	// Step 4: Generate domains_generated.go from template
	err = generateDomainsFile(registry)
	if err != nil {
		log.Fatalf("Failed to generate domains file: %v", err)
	}

	log.Println("✓ Generated: pkg/types/domains_generated.go")
	log.Println("✅ Domain generation complete!")
}

// readSpecIndex reads and parses .specs/index.json
func readSpecIndex(path string) (*SpecIndex, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var index SpecIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &index, nil
}

// readDomainConfig reads and parses domain_config.yaml
func readDomainConfig(path string) (*DomainConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var config DomainConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &config, nil
}

// generateDomainsFile generates pkg/types/domains_generated.go from registry
func generateDomainsFile(registry *GeneratedDomainRegistry) error {
	// Load template
	templatePath := "scripts/templates/domains.go.tmpl"
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create output directory if needed
	outputDir := "pkg/types"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Open output file
	outputPath := filepath.Join(outputDir, "domains_generated.go")
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() {
		_ = outFile.Close()
	}()

	// Sort domains for consistent output (idempotent generation)
	sortedDomains := make(map[string]*DomainInfo)
	keys := make([]string, 0, len(registry.Domains))
	for k := range registry.Domains {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Rebuild map in sorted order (for template iteration)
	data := struct {
		Version string
		Domains []struct {
			Name string
			Info *DomainInfo
		}
	}{
		Version: registry.Version,
		Domains: make([]struct {
			Name string
			Info *DomainInfo
		}, len(keys)),
	}

	for i, key := range keys {
		data.Domains[i].Name = key
		data.Domains[i].Info = registry.Domains[key]
		sortedDomains[key] = registry.Domains[key]
	}

	// Execute template
	if err := tmpl.Execute(outFile, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// titleCase converts snake_case to Title Case
func titleCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(string(part[0])) + part[1:]
		}
	}
	return strings.Join(parts, " ")
}
