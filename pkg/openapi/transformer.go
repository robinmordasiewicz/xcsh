// Package openapi provides utilities for parsing OpenAPI 3.0 specifications
// and generating example JSON for F5 XC resources.
package openapi

import (
	"regexp"
	"strings"
)

// TransformConfig holds configuration for spec transformation
type TransformConfig struct {
	// VesctlToBinary transforms "f5xcctl" references to the current binary name
	VesctlToBinary string
	// TransformEnvVars transforms VES_* environment variables to F5XC_*
	TransformEnvVars bool
}

// DefaultTransformConfig returns the default transformation configuration
func DefaultTransformConfig() *TransformConfig {
	return &TransformConfig{
		VesctlToBinary:   "f5xcctl",
		TransformEnvVars: true,
	}
}

// envVarPattern matches VES_* environment variable patterns
var envVarPattern = regexp.MustCompile(`\bVES_([A-Z_]+)\b`)

// TransformSpecReferences transforms legacy references in the spec to current branding.
// This includes:
// - "f5xcctl" → "f5xcctl" (or configured binary name)
// - VES_* environment variables → F5XC_*
//
// This function modifies the spec in place.
func TransformSpecReferences(spec *Spec, config *TransformConfig) {
	if spec == nil {
		return
	}
	if config == nil {
		config = DefaultTransformConfig()
	}

	// Transform info section
	spec.Info.Title = transformText(spec.Info.Title, config)
	spec.Info.Description = transformText(spec.Info.Description, config)

	// Transform all schemas
	for name, schema := range spec.Components.Schemas {
		transformSchema(schema, config)
		// Update the map in case we need to track changes
		spec.Components.Schemas[name] = schema
	}
}

// transformSchema recursively transforms text fields in a schema
func transformSchema(schema *Schema, config *TransformConfig) {
	if schema == nil {
		return
	}

	// Transform description and title
	schema.Description = transformText(schema.Description, config)
	schema.Title = transformText(schema.Title, config)

	// Transform F5 XC specific extensions
	schema.XVesExample = transformText(schema.XVesExample, config)
	schema.XDisplayName = transformText(schema.XDisplayName, config)

	// Transform nested properties
	for _, propSchema := range schema.Properties {
		transformSchema(propSchema, config)
	}

	// Transform array items
	if schema.Items != nil {
		transformSchema(schema.Items, config)
	}
}

// transformText applies all text transformations
func transformText(text string, config *TransformConfig) string {
	if text == "" {
		return text
	}

	// Transform f5xcctl references
	if config.VesctlToBinary != "" {
		// Case-insensitive replacement for various f5xcctl patterns
		text = strings.ReplaceAll(text, "f5xcctl", config.VesctlToBinary)
		text = strings.ReplaceAll(text, "Vesctl", capitalizeFirst(config.VesctlToBinary))
		text = strings.ReplaceAll(text, "VESCTL", strings.ToUpper(config.VesctlToBinary))
	}

	// Transform environment variables
	if config.TransformEnvVars {
		text = transformEnvVars(text)
	}

	return text
}

// transformEnvVars transforms VES_* environment variables to F5XC_*
func transformEnvVars(text string) string {
	return envVarPattern.ReplaceAllStringFunc(text, func(match string) string {
		// Extract the suffix after VES_
		suffix := strings.TrimPrefix(match, "VES_")
		return "F5XC_" + suffix
	})
}

// capitalizeFirst returns the string with the first letter capitalized
func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// LoadAllSpecsWithTransform loads all OpenAPI specifications from a directory
// and applies transformations to normalize legacy references.
func LoadAllSpecsWithTransform(dir string, config *TransformConfig) (map[string]*Spec, error) {
	specs, err := LoadAllSpecs(dir)
	if err != nil {
		return nil, err
	}

	// Apply transformations to all specs
	for _, spec := range specs {
		TransformSpecReferences(spec, config)
	}

	return specs, nil
}

// ParseSpecWithTransform parses an OpenAPI specification and applies transformations
func ParseSpecWithTransform(filePath string, config *TransformConfig) (*Spec, error) {
	spec, err := ParseSpec(filePath)
	if err != nil {
		return nil, err
	}

	TransformSpecReferences(spec, config)
	return spec, nil
}
