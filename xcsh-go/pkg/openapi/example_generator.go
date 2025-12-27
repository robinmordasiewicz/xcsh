package openapi

import (
	"encoding/json"
	"sort"
	"strings"
)

// ExampleGenerator generates example JSON from OpenAPI schemas
type ExampleGenerator struct {
	spec       *Spec
	maxDepth   int
	visitStack map[string]bool // Track visited refs to prevent infinite recursion
}

// NewExampleGenerator creates a new example generator for the given spec
func NewExampleGenerator(spec *Spec) *ExampleGenerator {
	return &ExampleGenerator{
		spec:       spec,
		maxDepth:   2, // Allow some nesting for useful examples
		visitStack: make(map[string]bool),
	}
}

// GenerateExample generates a JSON example for a CreateRequest schema
func (g *ExampleGenerator) GenerateExample(resourceName string) (string, error) {
	// Reset visit stack
	g.visitStack = make(map[string]bool)

	// Find the CreateRequest schema
	createReq := g.spec.FindCreateRequestSchema(resourceName)
	if createReq == nil {
		return "", nil // No schema found
	}

	// Generate the example
	example := g.generateFromSchema(createReq, 0)
	if example == nil {
		return "", nil
	}

	// Pretty print with indentation
	jsonBytes, err := json.MarshalIndent(example, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// generateFromSchema generates an example value from a schema
func (g *ExampleGenerator) generateFromSchema(schema *Schema, depth int) interface{} {
	if schema == nil || depth > g.maxDepth {
		return nil
	}

	// Handle $ref
	if schema.Ref != "" {
		refName := strings.TrimPrefix(schema.Ref, "#/components/schemas/")
		if g.visitStack[refName] {
			return nil // Prevent infinite recursion
		}
		g.visitStack[refName] = true
		defer func() { g.visitStack[refName] = false }()

		resolved := g.spec.ResolveRef(schema.Ref)
		if resolved == nil {
			return nil
		}
		return g.generateFromSchema(resolved, depth)
	}

	switch schema.Type {
	case "object":
		return g.generateObject(schema, depth)
	case "array":
		return g.generateArray(schema, depth)
	case "string":
		return g.generateString(schema)
	case "integer", "number":
		return g.generateNumber(schema)
	case "boolean":
		return g.generateBoolean(schema)
	default:
		// If type is not set but has properties, treat as object
		if len(schema.Properties) > 0 {
			return g.generateObject(schema, depth)
		}
		return nil
	}
}

// generateObject generates an example object from a schema
func (g *ExampleGenerator) generateObject(schema *Schema, depth int) map[string]interface{} {
	if schema.Properties == nil {
		return nil
	}

	result := make(map[string]interface{})

	// Collect fields to include
	fieldsToInclude := g.selectFieldsToInclude(schema, depth)

	// Sort keys for deterministic output
	keys := make([]string, 0, len(fieldsToInclude))
	for k := range fieldsToInclude {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		propSchema := schema.Properties[key]
		if propSchema == nil {
			continue
		}

		value := g.generateFromSchema(propSchema, depth+1)
		// Skip nil values to keep output clean
		if value == nil {
			continue
		}

		// Skip empty arrays and objects
		if arr, ok := value.([]interface{}); ok && len(arr) == 0 {
			continue
		}
		if obj, ok := value.(map[string]interface{}); ok && len(obj) == 0 {
			continue
		}

		result[key] = value
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

// selectFieldsToInclude determines which fields to include in the example
// Balance between useful examples and avoiding noise
func (g *ExampleGenerator) selectFieldsToInclude(schema *Schema, depth int) map[string]bool {
	fields := make(map[string]bool)

	// Include fields based on their usefulness
	for name, prop := range schema.Properties {
		// Skip disable_* fields - these are rarely useful in examples
		if strings.HasPrefix(name, "disable_") {
			continue
		}

		// Skip "gc_spec" internal fields
		if name == "gc_spec" {
			continue
		}

		// Include required fields (from schema.Required array)
		for _, req := range schema.Required {
			if req == name {
				fields[name] = true
				break
			}
		}

		// Include x-ves-required="true" fields
		if prop.XVesRequired == "true" {
			fields[name] = true
			continue
		}

		// Include fields with explicit examples - these are documented and useful
		if prop.XVesExample != "" {
			fields[name] = true
			continue
		}

		// Include simple primitive types with defaults or enums
		if prop.Type == "string" || prop.Type == "integer" || prop.Type == "number" || prop.Type == "boolean" {
			if prop.Default != nil || len(prop.Enum) > 0 {
				fields[name] = true
			}
		}
	}

	// Add commonly-used fields at all levels
	commonFields := []string{
		// Identifiers
		"name", "namespace", "tenant",
		// Networking
		"domains", "http", "https", "port", "ports",
		// Load balancing
		"origin_servers", "origin_pool", "default_route_pools", "default_pool",
		"weight", "priority",
		// Security
		"app_firewall", "waf_exclusion_rules",
		// TLS
		"tls_config", "use_tls",
	}
	for _, name := range commonFields {
		if _, exists := schema.Properties[name]; exists {
			fields[name] = true
		}
	}

	// Handle oneOf fields - keep only the first/simplest option
	for _, optionsJSON := range schema.XVesOneOfFields {
		var options []string
		if err := json.Unmarshal([]byte(optionsJSON), &options); err != nil {
			continue
		}
		// Remove all oneOf options except the first useful one
		foundFirst := false
		for _, opt := range options {
			if _, included := fields[opt]; included {
				if foundFirst {
					delete(fields, opt)
				} else {
					foundFirst = true
				}
			}
		}
	}

	return fields
}

// generateArray generates an example array from a schema
func (g *ExampleGenerator) generateArray(schema *Schema, depth int) []interface{} {
	if schema.Items == nil {
		return nil
	}

	// Generate one item for the array
	item := g.generateFromSchema(schema.Items, depth+1)
	if item == nil {
		return nil
	}

	// Skip if the item is an empty object
	if obj, ok := item.(map[string]interface{}); ok && len(obj) == 0 {
		return nil
	}

	return []interface{}{item}
}

// generateString generates an example string value
func (g *ExampleGenerator) generateString(schema *Schema) string {
	// Use x-ves-example if available
	if schema.XVesExample != "" {
		return schema.XVesExample
	}

	// Use default if available
	if schema.Default != nil {
		if s, ok := schema.Default.(string); ok {
			return s
		}
	}

	// Use first enum value if available
	if len(schema.Enum) > 0 {
		if s, ok := schema.Enum[0].(string); ok {
			return s
		}
	}

	// Generate placeholder based on field type hints
	title := strings.ToLower(schema.Title)
	desc := strings.ToLower(schema.Description)

	if strings.Contains(title, "name") || strings.Contains(desc, "name of") {
		return "example-resource"
	}
	if strings.Contains(title, "namespace") {
		return "example-namespace"
	}
	if strings.Contains(title, "domain") || strings.Contains(desc, "domain") {
		return "www.example.com"
	}

	return "example-value"
}

// generateNumber generates an example number value
func (g *ExampleGenerator) generateNumber(schema *Schema) interface{} {
	// Use x-ves-example if available
	if schema.XVesExample != "" {
		// Try to parse as integer first
		var intVal int
		if err := json.Unmarshal([]byte(schema.XVesExample), &intVal); err == nil {
			return intVal
		}
	}

	// Use default if available
	if schema.Default != nil {
		return schema.Default
	}

	// Return sensible defaults based on field hints
	title := strings.ToLower(schema.Title)
	if strings.Contains(title, "port") {
		return 80
	}
	if strings.Contains(title, "weight") {
		return 1
	}
	if strings.Contains(title, "priority") {
		return 1
	}

	return 1
}

// generateBoolean generates an example boolean value
func (g *ExampleGenerator) generateBoolean(schema *Schema) bool {
	// Use default if available
	if schema.Default != nil {
		if b, ok := schema.Default.(bool); ok {
			return b
		}
	}
	return false
}

// GenerateMetadataExample generates a metadata example object
func (g *ExampleGenerator) GenerateMetadataExample(resourceName string) map[string]interface{} {
	// Convert resource name to a more readable form
	displayName := strings.ReplaceAll(resourceName, "_", "-")

	return map[string]interface{}{
		"name":      "example-" + displayName,
		"namespace": "example-namespace",
	}
}

// GenerateCreateRequestExample generates a complete CreateRequest example
func (g *ExampleGenerator) GenerateCreateRequestExample(resourceName string) (string, error) {
	// Reset visit stack
	g.visitStack = make(map[string]bool)

	// Find the spec schema
	specSchema := g.spec.FindCreateSpecTypeSchema(resourceName)
	if specSchema == nil {
		return "", nil
	}

	// Generate spec example
	specExample := g.generateFromSchema(specSchema, 0)

	// For http_loadbalancer specifically, add domains if not present
	if resourceName == "http_loadbalancer" {
		if specMap, ok := specExample.(map[string]interface{}); ok {
			if _, hasDomains := specMap["domains"]; !hasDomains {
				specMap["domains"] = []interface{}{"www.example.com"}
			}
		}
	}

	// Clean up the spec - remove empty/nil values
	specExample = g.cleanupExample(specExample)

	if specExample == nil {
		specExample = map[string]interface{}{}
	}

	// Build complete CreateRequest
	createRequest := map[string]interface{}{
		"metadata": g.GenerateMetadataExample(resourceName),
		"spec":     specExample,
	}

	// Pretty print with indentation
	jsonBytes, err := json.MarshalIndent(createRequest, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// cleanupExample removes nil values and empty collections from the example
func (g *ExampleGenerator) cleanupExample(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	switch val := v.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, vv := range val {
			cleaned := g.cleanupExample(vv)
			if cleaned != nil {
				result[k] = cleaned
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result

	case []interface{}:
		var result []interface{}
		for _, item := range val {
			cleaned := g.cleanupExample(item)
			if cleaned != nil {
				result = append(result, cleaned)
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result

	default:
		return v
	}
}
