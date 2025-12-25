// Package types provides schema intelligence data structures for AI-assisted CLI usage.
package types

import "encoding/json"

// ResourceSchemaInfo contains AI-friendly schema intelligence for a resource.
// This enables AI assistants to deterministically configure F5 XC resources
// by understanding field constraints, mutual exclusivity, and cascading dependencies.
type ResourceSchemaInfo struct {
	// ResourceName is the canonical name of the resource (e.g., "http_loadbalancer")
	ResourceName string `json:"resource_name"`

	// Description provides a human-readable summary of the resource
	Description string `json:"description"`

	// Fields contains metadata for all configurable fields
	Fields map[string]FieldInfo `json:"fields"`

	// OneOfGroups lists all mutually exclusive field groups
	// AI assistants should ensure only one choice per group is set
	OneOfGroups []OneOfGroup `json:"oneof_groups"`

	// DecisionTree provides a navigable graph of configuration decisions
	// AI assistants should follow this tree to gather all required fields
	DecisionTree *DecisionNode `json:"decision_tree,omitempty"`

	// RequiredFields lists fields that are always required regardless of choices
	RequiredFields []string `json:"required_fields"`

	// RequiredTier is the minimum subscription tier required for this resource
	// Values: "STANDARD", "ADVANCED" (empty means no restriction)
	RequiredTier string `json:"required_tier,omitempty"`

	// RequiredAddons lists addon services that must be subscribed for this resource
	// e.g., ["bot-defense", "api-security"]
	RequiredAddons []string `json:"required_addons,omitempty"`

	// HelpAnnotation is shown in --help output for tier-restricted resources
	// e.g., "[Requires Advanced]"
	HelpAnnotation string `json:"help_annotation,omitempty"`

	// TierRestrictedFields maps field paths to their subscription tier requirements
	// AI assistants can use this to determine which fields are available based on tier
	TierRestrictedFields map[string]TierRequirement `json:"tier_restricted_fields,omitempty"`

	// MinimumConfiguration provides a copy-paste ready minimal configuration example
	// AI assistants use this to generate working configurations with minimum required fields
	MinimumConfiguration *MinimumConfigSpec `json:"minimum_configuration,omitempty"`
}

// TierRequirement defines the subscription requirement for a field
type TierRequirement struct {
	// MinimumTier is the minimum subscription tier required
	// Values: "STANDARD", "ADVANCED"
	MinimumTier string `json:"minimum_tier"`

	// RequiredAddons lists addon services that must be subscribed
	RequiredAddons []string `json:"required_addons,omitempty"`

	// Description explains why this tier is required
	Description string `json:"description,omitempty"`
}

// MinimumConfigSpec provides the minimum viable configuration for a resource.
// AI assistants use this to generate copy-paste ready examples.
type MinimumConfigSpec struct {
	// Description explains what this minimal configuration achieves
	Description string `json:"description"`

	// RequiredFields lists the fields that must be provided
	RequiredFields []string `json:"required_fields"`

	// ExampleYAML contains a complete, minimal YAML configuration
	ExampleYAML string `json:"example_yaml"`

	// ExampleCommand shows the CLI command to create the resource
	ExampleCommand string `json:"example_command"`

	// Domain is the xcsh domain for this resource (e.g., "cdn", "virtual")
	Domain string `json:"domain"`
}

// FieldInfo contains complete metadata for a single field.
// AI assistants use this to understand field constraints and relationships.
type FieldInfo struct {
	// Name is the field name as used in YAML/JSON configuration
	Name string `json:"name"`

	// Type is the field's data type (string, integer, boolean, object, array)
	Type string `json:"type"`

	// Description explains what the field does
	Description string `json:"description"`

	// Required indicates if this field must be set (context-independent)
	Required bool `json:"required"`

	// Computed indicates this field is set by the API and should not be provided in create requests.
	// Examples: tenant in ObjectRef, uid, creation_timestamp
	// AI assistants should omit these fields when generating configurations.
	Computed bool `json:"computed,omitempty"`

	// ComputedReason explains why the field is computed (for AI assistants)
	// Examples: "Set by API from authentication context", "Generated unique identifier"
	ComputedReason string `json:"computed_reason,omitempty"`

	// Immutable indicates the field cannot be changed after resource creation.
	// AI assistants should warn users when attempting to modify immutable fields.
	Immutable bool `json:"immutable,omitempty"`

	// ImmutableReason explains why the field is immutable
	// Examples: "Resource identifier", "Requires re-creation to change"
	ImmutableReason string `json:"immutable_reason,omitempty"`

	// Deprecated indicates the field should not be used in new configurations
	Deprecated bool `json:"deprecated,omitempty"`

	// DeprecatedMessage provides guidance for deprecated fields
	DeprecatedMessage string `json:"deprecated_message,omitempty"`

	// Enum lists valid values for constrained string fields
	// AI assistants should only use values from this list
	Enum []string `json:"enum,omitempty"`

	// Default is the default value if not specified
	Default interface{} `json:"default,omitempty"`

	// ExclusiveWith lists fields that cannot be set alongside this field
	// Parsed from "Exclusive with [field1 field2]" in descriptions
	ExclusiveWith []string `json:"exclusive_with,omitempty"`

	// RequiresFields lists fields that must also be set when this field is set
	RequiresFields []string `json:"requires_fields,omitempty"`

	// ValidationRules contains F5 XC-specific validation constraints
	// Example: {"ves.io.schema.rules.uint32.gte": 1, "ves.io.schema.rules.uint32.lte": 65535}
	ValidationRules map[string]interface{} `json:"validation_rules,omitempty"`

	// NestedFields contains child field metadata for object types
	NestedFields map[string]FieldInfo `json:"nested_fields,omitempty"`
}

// OneOfGroup represents a set of mutually exclusive field choices.
// Only one field from a group can be set in a valid configuration.
type OneOfGroup struct {
	// GroupName identifies this group (e.g., "loadbalancer_type", "tls_choice")
	GroupName string `json:"group_name"`

	// Description explains what choice this group represents
	Description string `json:"description"`

	// Choices lists the field names that are mutually exclusive
	// AI assistants must select exactly one (if required) or at most one
	Choices []string `json:"choices"`

	// Required indicates if one choice must be selected
	Required bool `json:"required"`
}

// DecisionNode represents a decision point in the configuration tree.
// AI assistants traverse this tree to determine what fields are needed.
type DecisionNode struct {
	// Field is the field name where a choice must be made
	Field string `json:"field"`

	// Description explains the decision to be made
	Description string `json:"description"`

	// Choices maps each possible value to its branch
	Choices map[string]*DecisionBranch `json:"choices"`
}

// DecisionBranch represents the outcome of selecting a specific choice.
// It defines what additional fields become required or available.
type DecisionBranch struct {
	// RequiredFields lists fields that must be set when this choice is selected
	RequiredFields []string `json:"required_fields,omitempty"`

	// OptionalFields lists fields that become available (but not required)
	OptionalFields []string `json:"optional_fields,omitempty"`

	// NextDecision points to the next decision that must be made
	// This creates cascading dependencies (e.g., https → tls_choice → mtls_choice)
	NextDecision *DecisionNode `json:"next_decision,omitempty"`

	// ValidationRules contains choice-specific validation constraints
	ValidationRules map[string]interface{} `json:"validation_rules,omitempty"`
}

// GetResourceSchema returns the schema intelligence for a resource type.
// Returns nil if no schema is available for the resource.
func GetResourceSchema(resourceName string) *ResourceSchemaInfo {
	if schema, ok := ResourceSchemas[resourceName]; ok {
		return &schema
	}
	return nil
}

// GetAllResourceSchemas returns all available resource schemas.
func GetAllResourceSchemas() map[string]ResourceSchemaInfo {
	return ResourceSchemas
}

// GetResourceSchemaNames returns the names of all resources with schema intelligence.
func GetResourceSchemaNames() []string {
	names := make([]string, 0, len(ResourceSchemas))
	for name := range ResourceSchemas {
		names = append(names, name)
	}
	return names
}

// Helper functions for generated code unmarshaling

// unmarshalFields unmarshals a JSON string into a map of FieldInfo.
func unmarshalFields(jsonStr string) map[string]FieldInfo {
	var fields map[string]FieldInfo
	if err := json.Unmarshal([]byte(jsonStr), &fields); err != nil {
		return make(map[string]FieldInfo)
	}
	return fields
}

// unmarshalOneOfGroups unmarshals a JSON string into a slice of OneOfGroup.
func unmarshalOneOfGroups(jsonStr string) []OneOfGroup {
	var groups []OneOfGroup
	if err := json.Unmarshal([]byte(jsonStr), &groups); err != nil {
		return []OneOfGroup{}
	}
	return groups
}

// unmarshalDecisionTree unmarshals a JSON string into a DecisionNode.
func unmarshalDecisionTree(jsonStr string) *DecisionNode {
	var node DecisionNode
	if err := json.Unmarshal([]byte(jsonStr), &node); err != nil {
		return nil
	}
	return &node
}

// unmarshalStringSlice unmarshals a JSON string into a string slice.
func unmarshalStringSlice(jsonStr string) []string {
	var slice []string
	if err := json.Unmarshal([]byte(jsonStr), &slice); err != nil {
		return []string{}
	}
	return slice
}
