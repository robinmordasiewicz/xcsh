package types

import "github.com/robinmordasiewicz/f5xcctl/pkg/naming"

// ResourceType defines a F5 XC resource type
type ResourceType struct {
	// Name is the resource type name (e.g., "http_loadbalancer")
	Name string

	// CLIName is the CLI command name (e.g., "http-loadbalancer")
	CLIName string

	// Description is a short description of the resource
	Description string

	// APIPath is the base API path for this resource
	// e.g., "/api/config/namespaces/{namespace}/http_loadbalancers"
	APIPath string

	// SupportsNamespace indicates if the resource is namespace-scoped
	SupportsNamespace bool

	// Operations supported by this resource type
	Operations ResourceOperations

	// DeleteConfig contains custom delete configuration
	// If nil, standard DELETE method is used
	DeleteConfig *DeleteConfig

	// RequiredTier is the minimum subscription tier required for this resource
	// Values: "STANDARD", "ADVANCED" (empty means no restriction)
	RequiredTier string

	// RequiredAddons lists addon services that must be subscribed for this resource
	// e.g., ["bot-defense", "api-security"]
	RequiredAddons []string

	// HelpAnnotation is shown in --help output for tier-restricted resources
	// e.g., "[Requires Advanced]"
	HelpAnnotation string
}

// HumanReadableName returns the human-readable name of the resource type
// with proper acronym casing (e.g., "http_loadbalancer" -> "HTTP Load Balancer")
func (r *ResourceType) HumanReadableName() string {
	return naming.ToHumanReadable(r.Name)
}

// DeleteConfig defines custom delete behavior for a resource type
type DeleteConfig struct {
	// PathSuffix is appended to the resource path (e.g., "/cascade_delete")
	PathSuffix string

	// Method is the HTTP method to use (e.g., "POST" instead of "DELETE")
	Method string

	// IncludeBody indicates if the request should include a body with the name
	IncludeBody bool
}

// ResourceOperations defines which operations are supported
type ResourceOperations struct {
	Create bool
	Get    bool
	List   bool
	Update bool
	Delete bool
	Status bool
}

// AllOperations returns ResourceOperations with all operations enabled
func AllOperations() ResourceOperations {
	return ResourceOperations{
		Create: true,
		Get:    true,
		List:   true,
		Update: true,
		Delete: true,
		Status: true,
	}
}

// ReadOnlyOperations returns ResourceOperations with only read operations
func ReadOnlyOperations() ResourceOperations {
	return ResourceOperations{
		Get:    true,
		List:   true,
		Status: true,
	}
}

// ResourceMetadata contains common metadata for F5 XC resources
type ResourceMetadata struct {
	Name        string            `json:"name" yaml:"name"`
	Namespace   string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
}

// ResourceSpec is the generic specification for a resource
type ResourceSpec map[string]interface{}

// Resource represents a generic F5 XC resource
type Resource struct {
	Metadata ResourceMetadata `json:"metadata" yaml:"metadata"`
	Spec     ResourceSpec     `json:"spec" yaml:"spec"`
}

// ResourceList represents a list of resources
type ResourceList struct {
	Items []Resource `json:"items" yaml:"items"`
}

// CreateRequest represents a resource creation request
type CreateRequest struct {
	Metadata ResourceMetadata `json:"metadata" yaml:"metadata"`
	Spec     ResourceSpec     `json:"spec" yaml:"spec"`
}

// UpdateRequest represents a resource update request
type UpdateRequest struct {
	Metadata ResourceMetadata `json:"metadata" yaml:"metadata"`
	Spec     ResourceSpec     `json:"spec" yaml:"spec"`
}

// StatusResponse represents a resource status response
type StatusResponse struct {
	Metadata ResourceMetadata       `json:"metadata" yaml:"metadata"`
	Status   map[string]interface{} `json:"status" yaml:"status"`
}
