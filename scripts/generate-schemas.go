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

	"github.com/robinmordasiewicz/xcsh/pkg/openapi"
	"github.com/robinmordasiewicz/xcsh/pkg/types"
)

var (
	outputFile      = flag.String("output", "pkg/types/schemas_generated.go", "Output file path")
	resourcesFile   = flag.String("resources-output", "pkg/types/resources_generated.go", "Resources output file path")
	specsDir        = flag.String("specs", ".specs/domains", "Directory containing OpenAPI specs")
	verbose         = flag.Bool("v", false, "Verbose output")
	strict          = flag.Bool("strict", false, "Fail on critical resource missing specs")
	validateOnly    = flag.Bool("validate", false, "Validate only, don't write output")
	reportMissing   = flag.Bool("report", false, "Report all missing specs and exit")
	updateResources = flag.Bool("update-resources", false, "Also regenerate resources_generated.go with full descriptions")
)

// criticalResources are resources that MUST have schemas generated
// These are core F5 XC resources that AI assistants commonly work with
// NOTE: This list can be overridden by x-ves-critical extension in upstream specs
// When upstream adds x-ves-critical markers, loadCriticalResourcesFromIndex() will use those instead
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

// criticalResourcesSource tracks where the critical resources list came from
var criticalResourcesSource = "hardcoded fallback"

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

// computedFieldPatterns maps field names to reasons why they are API-computed
// These fields should not be provided in create/update requests as the API sets them
var computedFieldPatterns = map[string]string{
	"tenant":                 "Set by API from authentication context",
	"uid":                    "Generated unique identifier by API",
	"kind":                   "Set by API based on object type",
	"creation_timestamp":     "Set by server on object creation",
	"modification_timestamp": "Updated by server on each modification",
	"creator_id":             "Set by API from authentication context",
	"creator_class":          "Set by API from authentication context",
	"object_index":           "Internal index maintained by API",
	"owner_view":             "Set by API based on permissions",
}

// immutableFieldPatterns maps field names to reasons why they cannot be changed after creation
var immutableFieldPatterns = map[string]string{
	"name":      "Resource identifier cannot be changed after creation",
	"namespace": "Resource namespace is immutable - requires re-creation to change",
}

// objectRefComputedFields are fields in ObjectRef types that are API-computed
var objectRefComputedFields = map[string]string{
	"tenant": "Auto-populated for object references from context",
	"uid":    "Auto-populated for object references by API",
	"kind":   "Auto-populated for object references based on target type",
}

// isComputedField checks if a field is computed by the API
func isComputedField(fieldName string, parentSchemaName string) (bool, string) {
	// Check direct patterns
	if reason, ok := computedFieldPatterns[fieldName]; ok {
		return true, reason
	}

	// Check ObjectRef computed fields (common in F5 XC)
	if strings.Contains(strings.ToLower(parentSchemaName), "objectref") ||
		strings.Contains(strings.ToLower(parentSchemaName), "object_ref") {
		if reason, ok := objectRefComputedFields[fieldName]; ok {
			return true, reason
		}
	}

	return false, ""
}

// isImmutableField checks if a field cannot be changed after resource creation
func isImmutableField(fieldName string) (bool, string) {
	if reason, ok := immutableFieldPatterns[fieldName]; ok {
		return true, reason
	}
	return false, ""
}

// indexSpec represents the structure of .specs/index.json
type indexSpec struct {
	Specifications []struct {
		Domain            string   `json:"domain"`
		CriticalResources []string `json:"x-ves-critical-resources,omitempty"`
	} `json:"specifications"`
	// Future: x-ves-critical-resources at top level for global critical list
	CriticalResources []string `json:"x-ves-critical-resources,omitempty"`
}

// loadCriticalResourcesFromIndex attempts to load critical resources from upstream spec index.json
// Returns true if loaded from index, false if using hardcoded fallback
func loadCriticalResourcesFromIndex(indexPath string, verbose bool) bool {
	data, err := os.ReadFile(indexPath)
	if err != nil {
		if verbose {
			fmt.Printf("Note: Could not read index.json (%v), using hardcoded critical resources\n", err)
		}
		return false
	}

	var index indexSpec
	if err := json.Unmarshal(data, &index); err != nil {
		if verbose {
			fmt.Printf("Note: Could not parse index.json (%v), using hardcoded critical resources\n", err)
		}
		return false
	}

	// Check for global x-ves-critical-resources (future upstream enhancement)
	if len(index.CriticalResources) > 0 {
		criticalResources = index.CriticalResources
		criticalResourcesSource = "upstream index.json (x-ves-critical-resources)"
		return true
	}

	// Future: Could also aggregate per-domain critical resources
	// For now, if upstream doesn't have the extension, use hardcoded fallback
	return false
}

func main() {
	flag.Parse()

	// Try to load critical resources from upstream spec metadata
	// Falls back to hardcoded list if x-ves-critical-resources not present in specs
	indexPath := filepath.Join(filepath.Dir(*specsDir), "index.json")
	loadCriticalResourcesFromIndex(indexPath, *verbose)

	if *verbose {
		fmt.Printf("Critical resources source: %s (%d resources)\n", criticalResourcesSource, len(criticalResources))
		fmt.Println("Loading OpenAPI specifications...")
	}

	// Load all specs from enriched source (downloaded at build time)
	specs, err := openapi.LoadAllSpecs(*specsDir)
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

	// Optionally update resources_generated.go with full descriptions and domain mappings
	if *updateResources {
		descriptionMap := buildDescriptionMap(mapper, allResources)
		domainMap := buildResourceDomainMap(specs, mapper, allResources)
		if err := writeResourcesFile(*resourcesFile, allResources, descriptionMap, domainMap); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing resources file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Updated %s with full descriptions and domain mappings\n", *resourcesFile)
	}
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

// getResourceDescription returns the best available description for a resource.
// Prefers the OpenAPI info.description (file-level) over schema.description.
func getResourceDescription(resourceName string, spec *openapi.Spec, schema *openapi.Schema) string {
	// Prefer info-level description (richer content)
	if spec.Info.Description != "" {
		return spec.Info.Description
	}
	// Fall back to schema-level description
	if schema != nil && schema.Description != "" {
		return schema.Description
	}
	return ""
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
		Description:    getResourceDescription(resourceName, spec, schema),
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

	// Extract MinimumConfiguration from x-ves-minimum-configuration extension (issue #152)
	if schema.XVesMinimumConfiguration != nil {
		info.MinimumConfiguration = &types.MinimumConfigSpec{
			Description:    schema.XVesMinimumConfiguration.Description,
			RequiredFields: schema.XVesMinimumConfiguration.RequiredFields,
			ExampleYAML:    schema.XVesMinimumConfiguration.ExampleYAML,
			ExampleCommand: schema.XVesMinimumConfiguration.ExampleCommand,
			Domain:         schema.XVesCLIDomain, // from x-ves-cli-domain
		}
	}

	return info
}

// extractFieldInfo extracts metadata for a single field
func extractFieldInfo(name string, schema *openapi.Schema, spec *openapi.Spec) types.FieldInfo {
	// Get the schema name for context (used in computed field detection)
	schemaName := ""
	if schema.Ref != "" {
		// Extract schema name from reference
		parts := strings.Split(schema.Ref, "/")
		if len(parts) > 0 {
			schemaName = parts[len(parts)-1]
		}
	}

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

	// Check if field is computed (API-populated)
	if computed, reason := isComputedField(name, schemaName); computed {
		info.Computed = true
		info.ComputedReason = reason
	}

	// Check if field is immutable (cannot be changed after creation)
	if immutable, reason := isImmutableField(name); immutable {
		info.Immutable = true
		info.ImmutableReason = reason
	}

	// Note: Deprecated field detection is available in FieldInfo but requires
	// x-ves-deprecated extension support in OpenAPI Schema struct. Fields can be
	// manually marked as deprecated when specific deprecation patterns are identified.

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

	// Marshal MinimumConfiguration if present (from x-ves-minimum-configuration extension)
	minConfigJSON := "nil"
	if schema.MinimumConfiguration != nil {
		bytes, err := json.Marshal(schema.MinimumConfiguration)
		if err != nil {
			return err
		}
		minConfigJSON = fmt.Sprintf("unmarshalMinimumConfig(%s)", escapeForGoString(string(bytes)))
	}

	entry := fmt.Sprintf(`	%q: {
		ResourceName:          %q,
		Description:           %q,
		Fields:                unmarshalFields(%s),
		OneOfGroups:           unmarshalOneOfGroups(%s),
		DecisionTree:          %s,
		RequiredFields:        unmarshalStringSlice(%s),
		MinimumConfiguration:  %s,
	},
`,
		name,
		schema.ResourceName,
		schema.Description,
		escapeForGoString(string(fieldsJSON)),
		escapeForGoString(string(oneOfJSON)),
		decisionJSON,
		escapeForGoString(string(requiredJSON)),
		minConfigJSON,
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

// buildDescriptionMap creates a map of resource names to their full descriptions from OpenAPI specs
func buildDescriptionMap(mapper *openapi.SpecMapper, resources []*types.ResourceType) map[string]string {
	descriptions := make(map[string]string)
	for _, rt := range resources {
		// Use OpenAPI spec description
		spec := mapper.FindSpec(rt.Name)
		if spec != nil && spec.Info.Description != "" {
			descriptions[rt.Name] = spec.Info.Description
		}
	}
	return descriptions
}

// domainIndexEntry represents a single entry in the domain index
type domainIndexEntry struct {
	Domain string `json:"domain"`
	Title  string `json:"title"`
	File   string `json:"file"`
}

// domainIndexFile represents the complete domain index
type domainIndexFile struct {
	Version        string             `json:"version"`
	Specifications []domainIndexEntry `json:"specifications"`
}

// loadDomainIndex loads the domain mappings from .specs/index.json
func loadDomainIndex(indexPath string) (map[string]string, error) {
	// Map from filename (e.g., "load_balancer.json") to domain name
	fileToDomain := make(map[string]string)

	data, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read domain index: %w", err)
	}

	var index domainIndexFile
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("failed to parse domain index: %w", err)
	}

	for _, spec := range index.Specifications {
		fileToDomain[spec.File] = spec.Domain
	}

	return fileToDomain, nil
}

// ResourceDomainInfo contains domain information for a resource
type ResourceDomainInfo struct {
	PrimaryDomain string   // Domain containing CreateRequest schema (authoritative)
	Domains       []string // All domains this resource appears in
}

// buildResourceDomainMap builds a mapping of resources to their domains from OpenAPI specs
func buildResourceDomainMap(specs map[string]*openapi.Spec, mapper *openapi.SpecMapper, resources []*types.ResourceType) map[string]*ResourceDomainInfo {
	domainMap := make(map[string]*ResourceDomainInfo)

	// Load domain index to extract domain from filenames
	fileToDomain, err := loadDomainIndex(".specs/index.json")
	if err != nil {
		if *verbose {
			fmt.Fprintf(os.Stderr, "Warning: could not load domain index: %v\n", err)
		}
		// Fall back to extracting domain from filename
	}

	// For each resource, find all specs that contain it
	for _, rt := range resources {
		info := &ResourceDomainInfo{
			Domains: []string{},
		}

		// Sort spec filePaths for deterministic iteration (maps have random order in Go)
		var sortedFilePaths []string
		for filePath := range specs {
			sortedFilePaths = append(sortedFilePaths, filePath)
		}
		sort.Strings(sortedFilePaths)

		// Check each spec file to see if it contains this resource
		for _, filePath := range sortedFilePaths {
			spec := specs[filePath]
			// Extract domain from filename (e.g., "load_balancer.json" -> "load_balancer")
			// The filepath from LoadAllSpecs is the full path, so extract just the filename
			filename := filepath.Base(filePath)
			domain := ""
			if d, ok := fileToDomain[filename]; ok {
				domain = d
			} else {
				// Fallback: extract from filename
				domain = strings.TrimSuffix(filename, ".json")
			}

			if domain == "other" || domain == "" {
				// Skip the "other" domain or invalid domains
				continue
			}

			// Check if resource is in this spec
			// First check for CreateRequest (primary domain indicator)
			if spec.FindCreateRequestSchema(rt.Name) != nil {
				// This domain has the CreateRequest - it's the primary domain
				if info.PrimaryDomain == "" {
					info.PrimaryDomain = domain
				}
				// Add to domains list if not already there
				if !contains(info.Domains, domain) {
					info.Domains = append(info.Domains, domain)
				}
			} else if spec.FindCreateSpecTypeSchema(rt.Name) != nil {
				// Has CreateSpecType schema but no CreateRequest - it's a secondary domain
				if !contains(info.Domains, domain) {
					info.Domains = append(info.Domains, domain)
				}
				// Only set PrimaryDomain if we haven't found one yet
				if info.PrimaryDomain == "" {
					info.PrimaryDomain = domain
				}
			}
		}

		// Sort domains for deterministic output
		sort.Strings(info.Domains)

		domainMap[rt.Name] = info
	}

	return domainMap
}

// writeResourcesFile generates resources_generated.go with full descriptions and domain mappings
func writeResourcesFile(outputPath string, resources []*types.ResourceType, descriptions map[string]string, domainMap map[string]*ResourceDomainInfo) error {
	// Ensure output directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Sort resources by name for deterministic output
	sorted := make([]*types.ResourceType, len(resources))
	copy(sorted, resources)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Name < sorted[j].Name
	})

	// Create output file
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	// Write header
	header := fmt.Sprintf(`package types

// Code generated from OpenAPI specifications. DO NOT EDIT.
// This file contains %d resource types parsed from F5 XC API specs

func init() {
	registerGeneratedResources()
}

func registerGeneratedResources() {
`, len(resources))
	if _, err := f.WriteString(header); err != nil {
		return err
	}

	// Write each resource
	for _, rt := range sorted {
		desc := rt.Description
		// Use full description from OpenAPI if available
		if fullDesc, ok := descriptions[rt.Name]; ok && fullDesc != "" {
			desc = fullDesc
		}

		domainInfo := domainMap[rt.Name]
		if err := writeResourceEntry(f, rt, desc, domainInfo); err != nil {
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

// writeResourceEntry writes a single resource registration to the file
func writeResourceEntry(f *os.File, rt *types.ResourceType, description string, domainInfo *ResourceDomainInfo) error {
	// Start the registration
	entry := fmt.Sprintf(`	Register(&ResourceType{
		Name:              %q,
		CLIName:           %q,
		Description:       %q,
		APIPath:           %q,
		SupportsNamespace: %t,
		Operations:        %s,`,
		rt.Name,
		rt.CLIName,
		description,
		rt.APIPath,
		rt.SupportsNamespace,
		formatOperations(rt.Operations),
	)

	if _, err := f.WriteString(entry); err != nil {
		return err
	}

	// Add DeleteConfig if present
	if rt.DeleteConfig != nil {
		deleteEntry := fmt.Sprintf(`
		DeleteConfig: &DeleteConfig{
			PathSuffix:  %q,
			Method:      %q,
			IncludeBody: %t,
		},`,
			rt.DeleteConfig.PathSuffix,
			rt.DeleteConfig.Method,
			rt.DeleteConfig.IncludeBody,
		)
		if _, err := f.WriteString(deleteEntry); err != nil {
			return err
		}
	}

	// Add domain information if available
	if domainInfo != nil && (domainInfo.PrimaryDomain != "" || len(domainInfo.Domains) > 0) {
		domainEntry := fmt.Sprintf(`
		PrimaryDomain:    %q,
		Domains:          []string{%s},`,
			domainInfo.PrimaryDomain,
			formatDomainSlice(domainInfo.Domains),
		)
		if _, err := f.WriteString(domainEntry); err != nil {
			return err
		}
	}

	// Close the registration
	if _, err := f.WriteString("\n\t})\n\n"); err != nil {
		return err
	}

	return nil
}

// formatDomainSlice formats a slice of domains as Go code
func formatDomainSlice(domains []string) string {
	if len(domains) == 0 {
		return ""
	}
	var parts []string
	for _, d := range domains {
		parts = append(parts, fmt.Sprintf("%q", d))
	}
	return strings.Join(parts, ", ")
}

// formatOperations formats ResourceOperations as Go code
func formatOperations(ops types.ResourceOperations) string {
	if ops.Create && ops.Get && ops.List && ops.Update && ops.Delete && ops.Status {
		return "AllOperations()"
	}
	if !ops.Create && ops.Get && ops.List && !ops.Update && !ops.Delete && ops.Status {
		return "ReadOnlyOperations()"
	}
	// Custom operations
	return fmt.Sprintf("ResourceOperations{Create: %t, Get: %t, List: %t, Update: %t, Delete: %t, Status: %t}",
		ops.Create, ops.Get, ops.List, ops.Update, ops.Delete, ops.Status)
}
