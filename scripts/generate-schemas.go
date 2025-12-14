//go:build ignore

// generate-schemas.go generates schemas_generated.go from OpenAPI specifications.
// Run with: go run scripts/generate-schemas.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/robinmordasiewicz/f5xcctl/pkg/openapi"
	"github.com/robinmordasiewicz/f5xcctl/pkg/types"
)

var (
	outputFile    = flag.String("output", "pkg/types/schemas_generated.go", "Output file path")
	specsDir      = flag.String("specs", "docs/specifications/api", "Directory containing OpenAPI specs")
	verbose       = flag.Bool("v", false, "Verbose output")
	strict        = flag.Bool("strict", false, "Fail on critical resource missing specs")
	validateOnly  = flag.Bool("validate", false, "Validate only, don't write output")
	reportMissing = flag.Bool("report", false, "Report all missing specs and exit")
)

// criticalResources are resources that MUST have schemas generated
// These are core F5 XC resources that AI assistants commonly work with
var criticalResources = []string{
	"http_loadbalancer",
	"tcp_loadbalancer",
	"origin_pool",
	"healthcheck",
	"app_firewall",
	"service_policy",
	"dns_zone",
	"dns_load_balancer",
	"certificate",
	"namespace",
	"virtual_site",
	"network_policy",
	"aws_vpc_site",
	"azure_vnet_site",
	"gcp_vpc_site",
}

// metadataOnlyResources are resources that intentionally have empty CreateSpecType schemas
// These resources only require metadata (name, namespace, labels) for creation - no additional spec fields
// The empty schema is expected behavior from the upstream F5 XC API, not a bug
var metadataOnlyResources = []string{
	"token",       // Site admission tokens - creation generates a token, no config needed
	"role",        // RBAC roles - standard CRUD uses empty spec (custom endpoint has fields)
	"tpm_manager", // TPM management - creation only needs metadata
}

// exclusiveWithRegex matches "Exclusive with [field1 field2 ...]" in descriptions
var exclusiveWithRegex = regexp.MustCompile(`[Ee]xclusive with \[([^\]]+)\]`)

func main() {
	flag.Parse()

	if *verbose {
		fmt.Println("Loading OpenAPI specifications...")
	}

	// Load all specs with transformation to normalize legacy references
	specs, err := openapi.LoadAllSpecsWithTransform(*specsDir, openapi.DefaultTransformConfig())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading specs: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("Loaded %d OpenAPI specifications\n", len(specs))
	}

	// Create spec mapper
	mapper := openapi.NewSpecMapper(specs)

	// Get all registered resource types
	allResources := types.All()
	if *verbose {
		fmt.Printf("Found %d registered resource types\n", len(allResources))
	}

	// Track missing specs for reporting
	var missingSpecs []string
	var missingCritical []string

	// Generate schemas for each resource
	schemas := make(map[string]types.ResourceSchemaInfo)
	generated := 0
	skipped := 0
	noCreate := 0

	for _, rt := range allResources {
		// Only process resources that support create (they have meaningful schemas)
		if !rt.Operations.Create {
			noCreate++
			continue
		}

		spec := mapper.FindSpec(rt.Name)
		if spec == nil {
			if *verbose {
				fmt.Printf("  Skipped %s (no spec found)\n", rt.Name)
			}
			missingSpecs = append(missingSpecs, rt.Name)
			if isCriticalResource(rt.Name) {
				missingCritical = append(missingCritical, rt.Name)
			}
			skipped++
			continue
		}

		schemaInfo := extractSchemaInfo(rt.Name, spec)
		if schemaInfo != nil {
			schemas[rt.Name] = *schemaInfo
			generated++
			if *verbose {
				fmt.Printf("  Generated schema for %s (%d fields, %d oneOf groups)\n",
					rt.Name, len(schemaInfo.Fields), len(schemaInfo.OneOfGroups))
			}
		} else {
			if *verbose {
				fmt.Printf("  Skipped %s (no schema extracted)\n", rt.Name)
			}
			missingSpecs = append(missingSpecs, rt.Name)
			if isCriticalResource(rt.Name) {
				missingCritical = append(missingCritical, rt.Name)
			}
			skipped++
		}
	}

	// Print summary
	fmt.Printf("\nGenerated %d schemas, skipped %d (no create: %d)\n", generated, skipped, noCreate)

	// Report mode: just report and exit
	if *reportMissing {
		printMissingReport(missingSpecs, missingCritical)
		if len(missingCritical) > 0 {
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Check for critical missing resources
	if len(missingCritical) > 0 {
		fmt.Fprintf(os.Stderr, "\n⚠️  WARNING: %d critical resources missing schemas:\n", len(missingCritical))
		for _, name := range missingCritical {
			fmt.Fprintf(os.Stderr, "   - %s\n", name)
		}

		if *strict {
			fmt.Fprintf(os.Stderr, "\nFailed: --strict mode requires all critical resources to have schemas\n")
			os.Exit(1)
		}
	}

	// Validate mode: don't write output
	if *validateOnly {
		fmt.Println("\nValidation complete (--validate mode, no output written)")
		if len(missingCritical) > 0 {
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Generate the output file
	if err := writeGeneratedFile(*outputFile, schemas); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s with %d resource schemas\n", *outputFile, generated)

	// Final validation
	validateGeneratedSchemas(schemas)
}

// isCriticalResource checks if a resource is in the critical list
func isCriticalResource(name string) bool {
	for _, critical := range criticalResources {
		if critical == name {
			return true
		}
	}
	return false
}

// isMetadataOnlyResource returns true if the resource is known to have an intentionally empty schema
func isMetadataOnlyResource(name string) bool {
	for _, r := range metadataOnlyResources {
		if r == name {
			return true
		}
	}
	return false
}

// printMissingReport prints a detailed report of missing specs
func printMissingReport(missing, critical []string) {
	fmt.Println("\n=== Schema Generation Report ===")

	if len(critical) > 0 {
		fmt.Println("\n❌ CRITICAL RESOURCES MISSING SCHEMAS:")
		for _, name := range critical {
			fmt.Printf("   - %s\n", name)
		}
	} else {
		fmt.Println("\n✅ All critical resources have schemas")
	}

	if len(missing) > 0 {
		fmt.Printf("\n📋 All resources missing specs (%d total):\n", len(missing))
		sort.Strings(missing)
		for _, name := range missing {
			marker := "  "
			if isCriticalResource(name) {
				marker = "❌"
			}
			fmt.Printf("   %s %s\n", marker, name)
		}
	}
}

// validateGeneratedSchemas performs post-generation validation
func validateGeneratedSchemas(schemas map[string]types.ResourceSchemaInfo) {
	var issues []string

	for name, schema := range schemas {
		// Check for empty schemas (skip known metadata-only resources)
		if len(schema.Fields) == 0 && len(schema.OneOfGroups) == 0 {
			if !isMetadataOnlyResource(name) {
				issues = append(issues, fmt.Sprintf("%s: empty schema (no fields or oneOf groups)", name))
			}
		}

		// Check for missing descriptions on critical resources
		if isCriticalResource(name) && schema.Description == "" {
			issues = append(issues, fmt.Sprintf("%s: missing description", name))
		}
	}

	if len(issues) > 0 {
		fmt.Println("\n⚠️  Schema quality warnings:")
		for _, issue := range issues {
			fmt.Printf("   - %s\n", issue)
		}
	}
}

// extractSchemaInfo extracts schema intelligence from an OpenAPI spec
func extractSchemaInfo(resourceName string, spec *openapi.Spec) *types.ResourceSchemaInfo {
	// Find the CreateSpecType schema (most complete definition)
	schema := spec.FindCreateSpecTypeSchema(resourceName)
	if schema == nil {
		// Fall back to CreateRequest schema
		schema = spec.FindCreateRequestSchema(resourceName)
	}
	if schema == nil {
		return nil
	}

	info := &types.ResourceSchemaInfo{
		ResourceName:   resourceName,
		Description:    schema.Description,
		Fields:         make(map[string]types.FieldInfo),
		OneOfGroups:    []types.OneOfGroup{},
		RequiredFields: schema.Required,
	}

	// Extract fields
	if schema.Properties != nil {
		for fieldName, fieldSchema := range schema.Properties {
			fieldInfo := extractFieldInfo(fieldName, fieldSchema, spec)
			info.Fields[fieldName] = fieldInfo
		}
	}

	// Extract oneOf groups from x-ves-oneof-field-* extensions
	for groupName, choicesJSON := range schema.XVesOneOfFields {
		choices := parseOneOfChoices(choicesJSON)
		if len(choices) > 0 {
			group := types.OneOfGroup{
				GroupName:   groupName,
				Description: fmt.Sprintf("Choose one of: %s", strings.Join(choices, ", ")),
				Choices:     choices,
				Required:    isOneOfRequired(groupName, schema),
			}
			info.OneOfGroups = append(info.OneOfGroups, group)
		}
	}

	// Sort oneOf groups for deterministic output
	sort.Slice(info.OneOfGroups, func(i, j int) bool {
		return info.OneOfGroups[i].GroupName < info.OneOfGroups[j].GroupName
	})

	// Build decision tree from oneOf groups
	info.DecisionTree = buildDecisionTree(info.OneOfGroups, info.Fields, spec, schema)

	return info
}

// extractFieldInfo extracts metadata for a single field
func extractFieldInfo(name string, schema *openapi.Schema, spec *openapi.Spec) types.FieldInfo {
	// Resolve reference if needed
	if schema.Ref != "" {
		resolved := spec.ResolveRef(schema.Ref)
		if resolved != nil {
			schema = resolved
		}
	}

	info := types.FieldInfo{
		Name:        name,
		Type:        determineFieldType(schema),
		Description: schema.Description,
		Required:    schema.XVesRequired == "true",
		Default:     schema.Default,
	}

	// Extract enum values
	if len(schema.Enum) > 0 {
		for _, v := range schema.Enum {
			if s, ok := v.(string); ok {
				info.Enum = append(info.Enum, s)
			}
		}
	}

	// Extract validation rules
	if len(schema.XVesValidationRules) > 0 {
		info.ValidationRules = make(map[string]interface{})
		for k, v := range schema.XVesValidationRules {
			info.ValidationRules[k] = v
		}
	}

	// Extract "Exclusive with [...]" from description
	info.ExclusiveWith = parseExclusiveWith(schema.Description)

	// Extract nested fields for object types
	if schema.Type == "object" && schema.Properties != nil {
		info.NestedFields = make(map[string]types.FieldInfo)
		for nestedName, nestedSchema := range schema.Properties {
			info.NestedFields[nestedName] = extractFieldInfo(nestedName, nestedSchema, spec)
		}
	}

	return info
}

// determineFieldType determines the JSON Schema type of a field
func determineFieldType(schema *openapi.Schema) string {
	if schema.Type != "" {
		return schema.Type
	}
	if schema.Properties != nil {
		return "object"
	}
	if schema.Items != nil {
		return "array"
	}
	if schema.Ref != "" {
		return "object" // References are typically objects
	}
	return "unknown"
}

// parseOneOfChoices parses the JSON array of choices from x-ves-oneof-field-*
func parseOneOfChoices(choicesJSON string) []string {
	var choices []string
	if err := json.Unmarshal([]byte(choicesJSON), &choices); err != nil {
		// Try parsing as a simple comma-separated string
		return strings.Split(choicesJSON, ",")
	}
	return choices
}

// parseExclusiveWith extracts field names from "Exclusive with [field1 field2]"
func parseExclusiveWith(description string) []string {
	matches := exclusiveWithRegex.FindStringSubmatch(description)
	if len(matches) < 2 {
		return nil
	}
	fields := strings.Fields(matches[1])
	return fields
}

// isOneOfRequired determines if a oneOf group requires a selection
func isOneOfRequired(groupName string, schema *openapi.Schema) bool {
	// Check if any of the group's fields are required
	for _, req := range schema.Required {
		if req == groupName {
			return true
		}
	}
	return false
}

// buildDecisionTree builds a decision tree from oneOf groups
func buildDecisionTree(groups []types.OneOfGroup, fields map[string]types.FieldInfo, spec *openapi.Spec, schema *openapi.Schema) *types.DecisionNode {
	if len(groups) == 0 {
		return nil
	}

	// Find the primary/root decision (first required oneOf group, or first group)
	var rootGroup *types.OneOfGroup
	for i := range groups {
		if groups[i].Required {
			rootGroup = &groups[i]
			break
		}
	}
	if rootGroup == nil && len(groups) > 0 {
		rootGroup = &groups[0]
	}
	if rootGroup == nil {
		return nil
	}

	root := &types.DecisionNode{
		Field:       rootGroup.GroupName,
		Description: rootGroup.Description,
		Choices:     make(map[string]*types.DecisionBranch),
	}

	// Build branches for each choice
	for _, choice := range rootGroup.Choices {
		branch := &types.DecisionBranch{
			RequiredFields: []string{},
			OptionalFields: []string{},
		}

		// Find fields that become required when this choice is selected
		if fieldInfo, ok := fields[choice]; ok {
			branch.RequiredFields = fieldInfo.RequiresFields
		}

		// Look for nested oneOf groups that apply to this choice
		for _, group := range groups {
			if group.GroupName != rootGroup.GroupName {
				// Check if this group is related to the current choice
				for _, groupChoice := range group.Choices {
					if fieldInfo, ok := fields[groupChoice]; ok {
						// Check if this field references the current choice
						for _, excl := range fieldInfo.ExclusiveWith {
							if excl == choice || contains(rootGroup.Choices, excl) {
								// This is a nested decision
								if branch.NextDecision == nil {
									branch.NextDecision = &types.DecisionNode{
										Field:       group.GroupName,
										Description: group.Description,
										Choices:     make(map[string]*types.DecisionBranch),
									}
								}
								break
							}
						}
					}
				}
			}
		}

		root.Choices[choice] = branch
	}

	return root
}

// contains checks if a slice contains a string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// writeGeneratedFile writes the generated Go file
func writeGeneratedFile(outputPath string, schemas map[string]types.ResourceSchemaInfo) error {
	// Ensure output directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Sort schema keys for deterministic output
	var keys []string
	for k := range schemas {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Create output file
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	// Write header
	header := `// Code generated by generate-schemas.go. DO NOT EDIT.
// This file contains schema intelligence for AI-assisted CLI usage.

package types

// ResourceSchemas maps resource names to their schema intelligence.
// These schemas are generated from OpenAPI specifications.
var ResourceSchemas = map[string]ResourceSchemaInfo{
`
	if _, err := f.WriteString(header); err != nil {
		return err
	}

	// Write each schema
	for _, key := range keys {
		schema := schemas[key]
		if err := writeSchemaEntry(f, key, &schema); err != nil {
			return err
		}
	}

	// Write footer
	footer := `}
`
	if _, err := f.WriteString(footer); err != nil {
		return err
	}

	return nil
}

// writeSchemaEntry writes a single schema entry to the file
func writeSchemaEntry(f *os.File, name string, schema *types.ResourceSchemaInfo) error {
	// Convert schema to JSON for embedding
	fieldsJSON, err := json.Marshal(schema.Fields)
	if err != nil {
		return err
	}

	oneOfJSON, err := json.Marshal(schema.OneOfGroups)
	if err != nil {
		return err
	}

	decisionJSON := "nil"
	if schema.DecisionTree != nil {
		bytes, err := json.Marshal(schema.DecisionTree)
		if err != nil {
			return err
		}
		decisionJSON = fmt.Sprintf("unmarshalDecisionTree(%s)", escapeForGoString(string(bytes)))
	}

	requiredJSON, err := json.Marshal(schema.RequiredFields)
	if err != nil {
		return err
	}

	entry := fmt.Sprintf(`	%q: {
		ResourceName:   %q,
		Description:    %q,
		Fields:         unmarshalFields(%s),
		OneOfGroups:    unmarshalOneOfGroups(%s),
		DecisionTree:   %s,
		RequiredFields: unmarshalStringSlice(%s),
	},
`,
		name,
		schema.ResourceName,
		escapeString(schema.Description),
		escapeForGoString(string(fieldsJSON)),
		escapeForGoString(string(oneOfJSON)),
		decisionJSON,
		escapeForGoString(string(requiredJSON)),
	)

	_, err = f.WriteString(entry)
	return err
}

// escapeForGoString escapes a string for use in Go interpreted string literal
func escapeForGoString(s string) string {
	// Use %q format which handles all escaping properly
	return fmt.Sprintf("%q", s)
}

// escapeString escapes a string for use in Go source
func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

// init ensures the package can be built
func init() {
	_ = types.All
}
