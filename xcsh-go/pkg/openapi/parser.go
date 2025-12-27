// Package openapi provides utilities for parsing OpenAPI 3.0 specifications
// and generating example JSON for F5 XC resources.
package openapi

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// Spec represents an OpenAPI 3.0 specification
type Spec struct {
	Components Components `json:"components"`
	Info       Info       `json:"info"`
}

// Info contains OpenAPI info section
type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

// Components contains the OpenAPI components section
type Components struct {
	Schemas map[string]*Schema `json:"schemas"`
}

// Schema represents an OpenAPI schema definition
type Schema struct {
	Type        string             `json:"type"`
	Description string             `json:"description"`
	Title       string             `json:"title"`
	Properties  map[string]*Schema `json:"properties"`
	Items       *Schema            `json:"items"`
	Ref         string             `json:"$ref"`
	Required    []string           `json:"required"`
	MinItems    int                `json:"minItems"`
	MaxItems    int                `json:"maxItems"`
	MinLength   int                `json:"minLength"`
	MaxLength   int                `json:"maxLength"`
	Format      string             `json:"format"`
	Enum        []interface{}      `json:"enum"`
	Default     interface{}        `json:"default"`

	// F5 XC specific extensions
	XVesExample         string            `json:"x-ves-example"`
	XVesRequired        string            `json:"x-ves-required"`
	XDisplayName        string            `json:"x-displayname"`
	XVesProtoMessage    string            `json:"x-ves-proto-message"`
	XVesProtoPackage    string            `json:"x-ves-proto-package"`
	XVesValidationRules map[string]string `json:"x-ves-validation-rules"`

	// OneOf field indicators
	XVesOneOfFields map[string]string `json:"-"` // Parsed from x-ves-oneof-field-* keys

	// New F5 XC minimum configuration extensions (from issue #152)
	XVesMinimumConfiguration *MinimumConfigExtension `json:"x-ves-minimum-configuration,omitempty"`
	XVesRequiredFor          *RequiredForExtension   `json:"x-ves-required-for,omitempty"`
	XVesCLIDomain            string                  `json:"x-ves-cli-domain,omitempty"`
	XVesCLIAliases           []string                `json:"x-ves-cli-aliases,omitempty"`
}

// MinimumConfigExtension represents the x-ves-minimum-configuration extension
type MinimumConfigExtension struct {
	Description             string                   `json:"description"`
	RequiredFields          []string                 `json:"required_fields"`
	MutuallyExclusiveGroups []MutuallyExclusiveGroup `json:"mutually_exclusive_groups,omitempty"`
	ExampleYAML             string                   `json:"example_yaml"`
	ExampleCommand          string                   `json:"example_command"`
}

// MutuallyExclusiveGroup represents a group of mutually exclusive fields
type MutuallyExclusiveGroup struct {
	Fields []string `json:"fields"`
	Reason string   `json:"reason"`
}

// RequiredForExtension represents an x-ves-required-for entry with context-specific requirement flags
type RequiredForExtension struct {
	MinimumConfig bool `json:"minimum_config"`
	Create        bool `json:"create"`
	Update        bool `json:"update"`
	Read          bool `json:"read"`
}

// ParseSpec parses an OpenAPI specification from a JSON file
func ParseSpec(filePath string) (*Spec, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec file: %w", err)
	}

	var spec Spec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse spec JSON: %w", err)
	}

	// Parse x-ves-oneof-field-* extensions
	parseOneOfExtensions(&spec, data)

	return &spec, nil
}

// parseOneOfExtensions extracts x-ves-oneof-field-* extensions from schemas
func parseOneOfExtensions(spec *Spec, rawData []byte) {
	// Parse into generic map to extract x-ves-oneof-field-* keys
	var raw map[string]interface{}
	if err := json.Unmarshal(rawData, &raw); err != nil {
		return
	}

	components, ok := raw["components"].(map[string]interface{})
	if !ok {
		return
	}

	schemas, ok := components["schemas"].(map[string]interface{})
	if !ok {
		return
	}

	for schemaName, schemaData := range schemas {
		schemaMap, ok := schemaData.(map[string]interface{})
		if !ok {
			continue
		}

		if spec.Components.Schemas[schemaName] == nil {
			continue
		}

		spec.Components.Schemas[schemaName].XVesOneOfFields = make(map[string]string)
		for key, value := range schemaMap {
			if strings.HasPrefix(key, "x-ves-oneof-field-") {
				fieldName := strings.TrimPrefix(key, "x-ves-oneof-field-")
				if strValue, ok := value.(string); ok {
					spec.Components.Schemas[schemaName].XVesOneOfFields[fieldName] = strValue
				}
			}
		}
	}
}

// GetSchema retrieves a schema by name from the specification
func (s *Spec) GetSchema(name string) *Schema {
	if s.Components.Schemas == nil {
		return nil
	}
	return s.Components.Schemas[name]
}

// ResolveRef resolves a $ref string to its schema
// Format: "#/components/schemas/SchemaName"
func (s *Spec) ResolveRef(ref string) *Schema {
	if !strings.HasPrefix(ref, "#/components/schemas/") {
		return nil
	}
	name := strings.TrimPrefix(ref, "#/components/schemas/")
	return s.GetSchema(name)
}

// FindCreateRequestSchema finds the CreateRequest schema for a resource type
func (s *Spec) FindCreateRequestSchema(resourceName string) *Schema {
	// Try common patterns
	patterns := []string{
		resourceName + "CreateRequest",
		strings.ReplaceAll(resourceName, "_", "") + "CreateRequest",
	}

	for _, pattern := range patterns {
		for name := range s.Components.Schemas {
			if strings.EqualFold(name, pattern) {
				return s.Components.Schemas[name]
			}
		}
	}

	return nil
}

// FindCreateSpecTypeSchema finds the CreateSpecType schema for a resource type
func (s *Spec) FindCreateSpecTypeSchema(resourceName string) *Schema {
	// Try common patterns - the spec type is usually viewsXXXCreateSpecType
	patterns := []string{
		"views" + resourceName + "CreateSpecType",
		resourceName + "CreateSpecType",
		strings.ReplaceAll(resourceName, "_", "") + "CreateSpecType",
	}

	for _, pattern := range patterns {
		for name := range s.Components.Schemas {
			if strings.EqualFold(name, pattern) {
				return s.Components.Schemas[name]
			}
		}
	}

	return nil
}

// IsRequired checks if a field is required in the schema
func (s *Schema) IsRequired(fieldName string) bool {
	for _, req := range s.Required {
		if req == fieldName {
			return true
		}
	}
	// Also check x-ves-required extension
	if s.Properties != nil {
		if prop, ok := s.Properties[fieldName]; ok {
			return prop.XVesRequired == "true"
		}
	}
	return false
}

// GetSchemaNames returns all schema names in the specification
func (s *Spec) GetSchemaNames() []string {
	names := make([]string, 0, len(s.Components.Schemas))
	for name := range s.Components.Schemas {
		names = append(names, name)
	}
	return names
}

// LoadAllSpecs loads all OpenAPI specifications from a directory
// Supports domain-organized enriched specs (load_balancer.json, networking.json, etc.)
// Skips metadata files like index.json
func LoadAllSpecs(dir string) (map[string]*Spec, error) {
	specs := make(map[string]*Spec)

	files, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to glob spec files: %w", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no spec files found in %s", dir)
	}

	// Sort files for deterministic loading order (filepath.Glob returns undefined order)
	sort.Strings(files)

	for _, file := range files {
		// Skip metadata files
		basename := filepath.Base(file)
		if basename == "index.json" {
			continue
		}

		spec, err := ParseSpec(file)
		if err != nil {
			// Log but continue with other files
			fmt.Fprintf(os.Stderr, "Warning: failed to parse %s: %v\n", file, err)
			continue
		}
		specs[file] = spec
	}

	if len(specs) == 0 {
		return nil, fmt.Errorf("no valid spec files loaded from %s", dir)
	}

	return specs, nil
}

// FindAllResourceSchemas discovers all resources in a spec by scanning for CreateRequest schemas
// Returns map of resource_name → CreateRequest schema
func (s *Spec) FindAllResourceSchemas() map[string]*Schema {
	resources := make(map[string]*Schema)

	for schemaName, schema := range s.Components.Schemas {
		if strings.HasSuffix(schemaName, "CreateRequest") {
			resourceName := extractResourceName(schemaName)
			if resourceName != "" {
				resources[resourceName] = schema
			}
		}
	}

	return resources
}

// extractResourceName converts schema name to resource name
// Example: "viewshttploadbalancerCreateRequest" → "http_loadbalancer"
func extractResourceName(schemaName string) string {
	// Remove "CreateRequest" suffix
	name := strings.TrimSuffix(schemaName, "CreateRequest")

	// Remove common prefixes
	name = strings.TrimPrefix(name, "views")
	name = strings.TrimPrefix(name, "public")

	if name == "" {
		return ""
	}

	// Convert camelCase to snake_case
	name = camelToSnake(name)

	return name
}

// camelToSnake converts camelCase to snake_case
func camelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

// ExtractDomainFromFile extracts the domain name from a spec file path
// Example: "load_balancer.json" → "load_balancer"
// Returns empty string for other filenames
func ExtractDomainFromFile(filePath string) string {
	basename := filepath.Base(filePath)

	// Domain files: load_balancer.json, networking.json
	if strings.HasSuffix(basename, ".json") && !strings.HasPrefix(basename, "docs-cloud-f5-com") {
		return strings.TrimSuffix(basename, ".json")
	}

	return ""
}

// FuzzyMatchResourceName tries to match a resource name with alternative formats
// Handles underscore variations: http_loadbalancer = httploadbalancer
func FuzzyMatchResourceName(query, target string) bool {
	// Exact match (case-insensitive)
	if strings.EqualFold(query, target) {
		return true
	}

	// Normalize underscores
	normalized := regexp.MustCompile(`_+`).ReplaceAllString
	queryNorm := normalized(query, "")
	targetNorm := normalized(target, "")

	return strings.EqualFold(queryNorm, targetNorm)
}
