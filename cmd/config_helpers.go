package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/robinmordasiewicz/f5xcctl/pkg/output"
	"github.com/robinmordasiewicz/f5xcctl/pkg/subscription"
	"github.com/robinmordasiewicz/f5xcctl/pkg/types"
)

// configurationFlags holds flags for configuration commands
type configurationFlags struct {
	namespace      string
	name           string
	inputFile      string
	jsonData       string
	responseFormat string
	mode           string
	labelKeys      []string
	labelValues    []string
	yes            bool // Skip confirmation for destructive operations
}

// getTierAnnotation returns the tier annotation for a resource type if it requires
// a higher subscription tier than Standard. Returns empty string if no annotation needed.
func getTierAnnotation(resourceType string) string {
	registry := subscription.NewFeatureRegistry()
	features := registry.GetFeaturesForResource(resourceType)
	if len(features) == 0 {
		return ""
	}

	// Find the highest tier requirement among all features
	for _, f := range features {
		if f.HelpAnnotation != "" {
			return f.HelpAnnotation
		}
	}
	return ""
}

// formatShortWithTier creates a Short description with tier annotation if needed
func formatShortWithTier(action, displayName, resourceType string) string {
	annotation := getTierAnnotation(resourceType)
	if annotation != "" {
		return fmt.Sprintf("%s %s %s", action, displayName, annotation)
	}
	return fmt.Sprintf("%s %s", action, displayName)
}

// runConfigList executes the list operation (f5xcctl compatible)
func runConfigList(rt *types.ResourceType, flags *configurationFlags) error {
	client := GetClient()
	if client == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	namespace := flags.namespace
	if !rt.SupportsNamespace {
		namespace = ""
	}

	path := rt.BuildAPIPath(namespace, "")
	resp, err := client.Get(ctx, path, nil)
	if err != nil {
		return fmt.Errorf("error listing object: %w", err)
	}

	if resp.StatusCode >= 400 {
		return formatAPIError("listing", "GET", path, resp.StatusCode, resp.Body)
	}

	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return output.Print(result, GetOutputFormat())
}

// runConfigGet executes the get operation (f5xcctl compatible)
func runConfigGet(rt *types.ResourceType, flags *configurationFlags) error {
	// Note: We accept any response-format value including GET_RSP_FORMAT_READ
	// (original f5xcctl has a bug that rejects this valid value)

	client := GetClient()
	if client == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	namespace := flags.namespace
	if !rt.SupportsNamespace {
		namespace = ""
	}

	path := rt.BuildAPIPath(namespace, flags.name)
	resp, err := client.Get(ctx, path, nil)
	if err != nil {
		return fmt.Errorf("error getting object: %w", err)
	}

	if resp.StatusCode >= 400 {
		return formatAPIError("getting", "GET", path, resp.StatusCode, resp.Body)
	}

	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Get defaults to YAML output (matching original f5xcctl)
	return output.Print(result, GetOutputFormatWithDefault("yaml"))
}

// ensureNamespaceExists checks if a namespace exists and creates it if not.
// Returns nil if namespace exists or was successfully created.
// Reserved namespaces (system, shared, default) are always assumed to exist.
func ensureNamespaceExists(ctx context.Context, namespace string) error {
	if namespace == "" {
		return nil
	}

	// Reserved namespaces are always assumed to exist
	reservedNamespaces := map[string]bool{
		"system":  true,
		"shared":  true,
		"default": true,
	}
	if reservedNamespaces[namespace] {
		return nil
	}

	client := GetClient()
	if client == nil {
		return fmt.Errorf("client not initialized")
	}

	// Try to get the namespace
	path := fmt.Sprintf("/api/web/namespaces/%s", namespace)
	resp, err := client.Get(ctx, path, nil)
	if err == nil && resp.StatusCode == 200 {
		// Namespace exists
		return nil
	}

	// If namespace doesn't exist (404), create it
	if resp != nil && resp.StatusCode == 404 {
		output.PrintInfo(fmt.Sprintf("Namespace '%s' does not exist, creating it...", namespace))

		createPayload := map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": namespace,
			},
			"spec": map[string]interface{}{},
		}

		createResp, err := client.Post(ctx, "/api/web/namespaces", createPayload)
		if err != nil {
			return fmt.Errorf("failed to create namespace '%s': %w", namespace, err)
		}

		if createResp.StatusCode >= 400 {
			// Check if it's a conflict (namespace was created by another process)
			if createResp.StatusCode == 409 {
				output.PrintInfo(fmt.Sprintf("Namespace '%s' already exists (created concurrently)", namespace))
				return nil
			}
			return fmt.Errorf("failed to create namespace '%s': status %d, body: %s",
				namespace, createResp.StatusCode, string(createResp.Body))
		}

		output.PrintInfo(fmt.Sprintf("Namespace '%s' created successfully", namespace))
		return nil
	}

	// Other error
	if err != nil {
		return fmt.Errorf("failed to check namespace '%s': %w", namespace, err)
	}
	return fmt.Errorf("failed to check namespace '%s': status %d", namespace, resp.StatusCode)
}

// runConfigCreate executes the create operation (f5xcctl compatible)
func runConfigCreate(rt *types.ResourceType, flags *configurationFlags) error {
	client := GetClient()
	if client == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	// Load resource from input
	resource, err := loadConfigResource(flags.inputFile, flags.jsonData)
	if err != nil {
		return fmt.Errorf("failed to load resource: %w", err)
	}

	// Get namespace from resource metadata
	namespace := ""
	if meta, ok := resource["metadata"].(map[string]interface{}); ok {
		if ns, ok := meta["namespace"].(string); ok {
			namespace = ns
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Ensure namespace exists before creating resource (auto-create if needed)
	if rt.SupportsNamespace && namespace != "" {
		if err := ensureNamespaceExists(ctx, namespace); err != nil {
			return fmt.Errorf("failed to ensure namespace exists: %w", err)
		}
	}

	path := rt.BuildAPIPath(namespace, "")
	resp, err := client.Post(ctx, path, resource)
	if err != nil {
		return fmt.Errorf("error creating object: %w", err)
	}

	if resp.StatusCode >= 400 {
		return formatAPIError("creating", "POST", path, resp.StatusCode, resp.Body)
	}

	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Print "Created" header (matching original f5xcctl)
	fmt.Println("Created")
	// Create defaults to YAML output (matching original f5xcctl)
	return output.Print(result, GetOutputFormatWithDefault("yaml"))
}

// runConfigDelete executes the delete operation (f5xcctl compatible)
func runConfigDelete(rt *types.ResourceType, flags *configurationFlags) error {
	client := GetClient()
	if client == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	// Require confirmation for destructive operations
	if !flags.yes {
		if IsNonInteractive() {
			return fmt.Errorf("--yes flag is required for delete operations in non-interactive mode")
		}
		// Prompt for confirmation
		fmt.Fprintf(os.Stderr, "Are you sure you want to delete %s '%s' in namespace '%s'? [y/N]: ", rt.Name, flags.name, flags.namespace)
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read confirmation: %w", err)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			output.PrintInfo("Delete operation cancelled")
			return nil
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	namespace := flags.namespace
	if !rt.SupportsNamespace {
		namespace = ""
	}

	path := rt.BuildAPIPath(namespace, flags.name)

	// Check if resource type has custom delete configuration
	if rt.DeleteConfig != nil {
		if rt.DeleteConfig.PathSuffix != "" {
			path = path + rt.DeleteConfig.PathSuffix
		}

		if rt.DeleteConfig.Method == "POST" {
			var body interface{}
			if rt.DeleteConfig.IncludeBody {
				body = map[string]interface{}{
					"name": flags.name,
				}
			}
			resp, err := client.Post(ctx, path, body)
			if err != nil {
				return fmt.Errorf("error deleting object: %w", err)
			}
			if resp.StatusCode >= 400 {
				return formatAPIError("deleting", "POST", path, resp.StatusCode, resp.Body)
			}
			output.PrintInfo(fmt.Sprintf("Deleted %s '%s' successfully", rt.Name, flags.name))
			return nil
		}
	}

	// Standard DELETE method
	resp, err := client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("error deleting object: %w", err)
	}

	if resp.StatusCode >= 400 {
		return formatAPIError("deleting", "DELETE", path, resp.StatusCode, resp.Body)
	}

	output.PrintInfo(fmt.Sprintf("Deleted %s '%s' successfully", rt.Name, flags.name))
	return nil
}

// runConfigReplace executes the replace operation (f5xcctl compatible)
func runConfigReplace(rt *types.ResourceType, flags *configurationFlags) error {
	client := GetClient()
	if client == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	// Load resource from input
	resource, err := loadConfigResource(flags.inputFile, flags.jsonData)
	if err != nil {
		return fmt.Errorf("failed to load resource: %w", err)
	}

	// Get name and namespace from resource
	var name, namespace string
	if meta, ok := resource["metadata"].(map[string]interface{}); ok {
		if n, ok := meta["name"].(string); ok {
			name = n
		}
		if ns, ok := meta["namespace"].(string); ok {
			namespace = ns
		}
	}

	if name == "" {
		return fmt.Errorf("resource name is required in metadata")
	}

	// Require confirmation for destructive operations
	if !flags.yes {
		if IsNonInteractive() {
			return fmt.Errorf("--yes flag is required for replace operations in non-interactive mode")
		}
		// Prompt for confirmation
		nsDisplay := namespace
		if nsDisplay == "" {
			nsDisplay = "default"
		}
		fmt.Fprintf(os.Stderr, "Are you sure you want to replace %s '%s' in namespace '%s'? [y/N]: ", rt.Name, name, nsDisplay)
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read confirmation: %w", err)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			output.PrintInfo("Replace operation cancelled")
			return nil
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Ensure namespace exists before replacing resource (auto-create if needed)
	if rt.SupportsNamespace && namespace != "" {
		if err := ensureNamespaceExists(ctx, namespace); err != nil {
			return fmt.Errorf("failed to ensure namespace exists: %w", err)
		}
	}

	path := rt.BuildAPIPath(namespace, name)
	resp, err := client.Put(ctx, path, resource)
	if err != nil {
		return fmt.Errorf("error replacing object: %w", err)
	}

	if resp.StatusCode >= 400 {
		return formatAPIError("replacing", "PUT", path, resp.StatusCode, resp.Body)
	}

	// Print only "Replaced" (matching original f5xcctl - no response body output)
	fmt.Println("Replaced")
	return nil
}

// runConfigStatus executes the status operation (f5xcctl compatible)
func runConfigStatus(rt *types.ResourceType, flags *configurationFlags) error {
	client := GetClient()
	if client == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	namespace := flags.namespace
	if !rt.SupportsNamespace {
		namespace = ""
	}

	path := rt.BuildAPIPath(namespace, flags.name) + "/status"
	resp, err := client.Get(ctx, path, nil)
	if err != nil {
		return fmt.Errorf("error getting status: %w", err)
	}

	if resp.StatusCode >= 400 {
		return formatAPIError("getting status", "GET", path, resp.StatusCode, resp.Body)
	}

	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Status defaults to YAML output (matching original f5xcctl)
	return output.Print(result, GetOutputFormatWithDefault("yaml"))
}

// runConfigApply executes the apply operation (create or replace)
func runConfigApply(rt *types.ResourceType, flags *configurationFlags) error {
	client := GetClient()
	if client == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	// Load resource from input
	resource, err := loadConfigResource(flags.inputFile, flags.jsonData)
	if err != nil {
		return fmt.Errorf("failed to load resource: %w", err)
	}

	// Get name and namespace from resource
	var name, namespace string
	if meta, ok := resource["metadata"].(map[string]interface{}); ok {
		if n, ok := meta["name"].(string); ok {
			name = n
		}
		if ns, ok := meta["namespace"].(string); ok {
			namespace = ns
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Ensure namespace exists before applying resource (auto-create if needed)
	if rt.SupportsNamespace && namespace != "" {
		if err := ensureNamespaceExists(ctx, namespace); err != nil {
			return fmt.Errorf("failed to ensure namespace exists: %w", err)
		}
	}

	// If mode is "new", only create (fail if exists)
	if flags.mode == "new" {
		path := rt.BuildAPIPath(namespace, "")
		resp, err := client.Post(ctx, path, resource)
		if err != nil {
			return fmt.Errorf("error creating object: %w", err)
		}
		if resp.StatusCode >= 400 {
			return formatAPIError("creating", "POST", path, resp.StatusCode, resp.Body)
		}
		var result interface{}
		if err := json.Unmarshal(resp.Body, &result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
		fmt.Println("Created")
		return output.Print(result, GetOutputFormatWithDefault("yaml"))
	}

	// Mode "always" - try to get first, then create or replace
	if name != "" {
		getPath := rt.BuildAPIPath(namespace, name)
		getResp, _ := client.Get(ctx, getPath, nil)
		if getResp != nil && getResp.StatusCode == 200 {
			// Resource exists, replace it
			resp, err := client.Put(ctx, getPath, resource)
			if err != nil {
				return fmt.Errorf("error replacing object: %w", err)
			}
			if resp.StatusCode >= 400 {
				return formatAPIError("replacing", "PUT", getPath, resp.StatusCode, resp.Body)
			}
			// Print only "Replaced" (matching original f5xcctl - no response body output)
			fmt.Println("Replaced")
			return nil
		}
	}

	// Resource doesn't exist, create it
	createPath := rt.BuildAPIPath(namespace, "")
	resp, err := client.Post(ctx, createPath, resource)
	if err != nil {
		return fmt.Errorf("error creating object: %w", err)
	}
	if resp.StatusCode >= 400 {
		return formatAPIError("creating", "POST", createPath, resp.StatusCode, resp.Body)
	}
	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	fmt.Println("Created")
	return output.Print(result, GetOutputFormatWithDefault("yaml"))
}

// runConfigPatch executes the patch operation
func runConfigPatch(rt *types.ResourceType, flags *configurationFlags) error {
	// Patch operation requires additional implementation
	// This is a placeholder that returns not implemented
	return fmt.Errorf("patch operation not yet implemented - use replace instead")
}

// runConfigAddLabels adds labels to a resource
func runConfigAddLabels(rt *types.ResourceType, flags *configurationFlags) error {
	if len(flags.labelKeys) == 0 {
		return fmt.Errorf("at least one --label-key is required")
	}
	if len(flags.labelKeys) != len(flags.labelValues) {
		return fmt.Errorf("number of --label-key and --label-value flags must match")
	}

	client := GetClient()
	if client == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	namespace := flags.namespace
	if !rt.SupportsNamespace {
		namespace = ""
	}

	// Build labels map
	labels := make(map[string]string)
	for i, key := range flags.labelKeys {
		labels[key] = flags.labelValues[i]
	}

	// Add labels via API
	path := rt.BuildAPIPath(namespace, flags.name) + "/add_labels"
	body := map[string]interface{}{
		"labels": labels,
	}

	resp, err := client.Post(ctx, path, body)
	if err != nil {
		return fmt.Errorf("error adding labels: %w", err)
	}

	if resp.StatusCode >= 400 {
		return formatAPIError("adding labels", "POST", path, resp.StatusCode, resp.Body)
	}

	output.PrintInfo(fmt.Sprintf("Added labels to %s '%s'", rt.Name, flags.name))
	return nil
}

// runConfigRemoveLabels removes labels from a resource
func runConfigRemoveLabels(rt *types.ResourceType, flags *configurationFlags) error {
	if len(flags.labelKeys) == 0 {
		return fmt.Errorf("at least one --label-key is required")
	}

	client := GetClient()
	if client == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	namespace := flags.namespace
	if !rt.SupportsNamespace {
		namespace = ""
	}

	// Remove labels via API
	path := rt.BuildAPIPath(namespace, flags.name) + "/remove_labels"
	body := map[string]interface{}{
		"keys": flags.labelKeys,
	}

	resp, err := client.Post(ctx, path, body)
	if err != nil {
		return fmt.Errorf("error removing labels: %w", err)
	}

	if resp.StatusCode >= 400 {
		return formatAPIError("removing labels", "POST", path, resp.StatusCode, resp.Body)
	}

	output.PrintInfo(fmt.Sprintf("Removed labels from %s '%s'", rt.Name, flags.name))
	return nil
}

// formatAPIError formats an API error to match original f5xcctl error format
func formatAPIError(operation, method, path string, statusCode int, body []byte) error {
	baseURL := serverURL
	// Capitalize first letter of operation
	capOperation := strings.ToUpper(operation[:1]) + operation[1:]
	return fmt.Errorf("error %s object: %s object: unsuccessful %s at URL %s%s, status code %d, body %s, err %%!s(<nil>)",
		operation, capOperation, method, baseURL, path, statusCode, string(body))
}

// loadConfigResource loads resource from input file or JSON data
func loadConfigResource(inputFile, jsonData string) (map[string]interface{}, error) {
	var data []byte
	var err error

	if inputFile != "" {
		data, err = os.ReadFile(inputFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
	} else if jsonData != "" {
		data = []byte(jsonData)
	} else {
		return nil, fmt.Errorf("either --input-file or --json-data is required")
	}

	var resource map[string]interface{}

	// Try YAML first (YAML is a superset of JSON)
	if err := yaml.Unmarshal(data, &resource); err != nil {
		// Try JSON if YAML fails
		if err := json.Unmarshal(data, &resource); err != nil {
			return nil, fmt.Errorf("failed to parse input (not valid YAML or JSON): %w", err)
		}
	}

	return resource, nil
}

