package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/robinmordasiewicz/vesctl/pkg/output"
	"github.com/robinmordasiewicz/vesctl/pkg/types"
)

// configurationFlags holds flags for configuration commands (vesctl compatibility)
type configurationFlags struct {
	namespace      string
	name           string
	inputFile      string
	jsonData       string
	responseFormat string
	mode           string
	labelKeys      []string
	labelValues    []string
	atSite         string
}

// configurationCmd represents the configuration command (vesctl compatibility)
var configurationCmd = &cobra.Command{
	Use:     "configuration",
	Aliases: []string{"cfg", "c"},
	Short:   "Configure object",
	Long:    `Configure object`,
	Example: `vesctl configuration create virtual_host`,
}

func init() {
	rootCmd.AddCommand(configurationCmd)

	// Add subcommands
	configurationCmd.AddCommand(buildConfigListCmd())
	configurationCmd.AddCommand(buildConfigGetCmd())
	configurationCmd.AddCommand(buildConfigCreateCmd())
	configurationCmd.AddCommand(buildConfigDeleteCmd())
	configurationCmd.AddCommand(buildConfigReplaceCmd())
	configurationCmd.AddCommand(buildConfigStatusCmd())
	configurationCmd.AddCommand(buildConfigApplyCmd())
	configurationCmd.AddCommand(buildConfigPatchCmd())
	configurationCmd.AddCommand(buildConfigAddLabelsCmd())
	configurationCmd.AddCommand(buildConfigRemoveLabelsCmd())
}

// buildConfigListCmd creates the list subcommand
func buildConfigListCmd() *cobra.Command {
	var flags configurationFlags

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List configuration objects",
		Example: "vesctl configuration list virtual_host",
	}

	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Namespace in which to list objects")

	// Add resource type subcommands
	for _, rt := range types.All() {
		rtCopy := rt
		subCmd := &cobra.Command{
			Use:   rt.Name,
			Short: fmt.Sprintf("List %s", rt.Name),
			RunE: func(cmd *cobra.Command, args []string) error {
				return runConfigList(rtCopy, &flags)
			},
		}
		cmd.AddCommand(subCmd)
	}

	return cmd
}

// buildConfigGetCmd creates the get subcommand
func buildConfigGetCmd() *cobra.Command {
	var flags configurationFlags

	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get configuration object",
		Example: "vesctl configuration get virtual_host <name>",
	}

	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Namespace in which to get object")
	cmd.PersistentFlags().StringVar(&flags.responseFormat, "response-format", "read", "Format in get response (default 'read', others 'replace-request')")

	// Add resource type subcommands
	for _, rt := range types.All() {
		rtCopy := rt
		subCmd := &cobra.Command{
			Use:   fmt.Sprintf("%s <name>", rt.Name),
			Short: fmt.Sprintf("Get %s", rt.Name),
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				flags.name = args[0]
				return runConfigGet(rtCopy, &flags)
			},
		}
		cmd.AddCommand(subCmd)
	}

	return cmd
}

// buildConfigCreateCmd creates the create subcommand
func buildConfigCreateCmd() *cobra.Command {
	var flags configurationFlags

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create configuration object",
		Example: "vesctl configuration create virtual_host -i <file>",
	}

	cmd.PersistentFlags().StringVarP(&flags.inputFile, "input-file", "i", "", "File containing CreateRequest contents in yaml form")
	cmd.PersistentFlags().StringVar(&flags.jsonData, "json-data", "", "Inline CreateRequest contents in json form")

	// Add resource type subcommands
	for _, rt := range types.All() {
		if !rt.Operations.Create {
			continue
		}
		rtCopy := rt
		subCmd := &cobra.Command{
			Use:   rt.Name,
			Short: fmt.Sprintf("Create %s", rt.Name),
			RunE: func(cmd *cobra.Command, args []string) error {
				return runConfigCreate(rtCopy, &flags)
			},
		}
		cmd.AddCommand(subCmd)
	}

	return cmd
}

// buildConfigDeleteCmd creates the delete subcommand
func buildConfigDeleteCmd() *cobra.Command {
	var flags configurationFlags

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete configuration object",
		Example: "vesctl configuration delete virtual_host <name>",
	}

	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Namespace in which to delete object")

	// Add resource type subcommands
	for _, rt := range types.All() {
		if !rt.Operations.Delete {
			continue
		}
		rtCopy := rt
		subCmd := &cobra.Command{
			Use:   fmt.Sprintf("%s <name>", rt.Name),
			Short: fmt.Sprintf("Delete %s", rt.Name),
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				flags.name = args[0]
				return runConfigDelete(rtCopy, &flags)
			},
		}
		cmd.AddCommand(subCmd)
	}

	return cmd
}

// buildConfigReplaceCmd creates the replace subcommand
func buildConfigReplaceCmd() *cobra.Command {
	var flags configurationFlags

	cmd := &cobra.Command{
		Use:     "replace",
		Short:   "Replace configuration object",
		Example: "vesctl configuration replace virtual_host -i <file>",
	}

	cmd.PersistentFlags().StringVarP(&flags.inputFile, "input-file", "i", "", "File containing ReplaceRequest content in yaml form")
	cmd.PersistentFlags().StringVar(&flags.jsonData, "json-data", "", "Inline ReplaceRequest contents in json form")

	// Add resource type subcommands
	for _, rt := range types.All() {
		if !rt.Operations.Update {
			continue
		}
		rtCopy := rt
		subCmd := &cobra.Command{
			Use:   rt.Name,
			Short: fmt.Sprintf("Replace %s", rt.Name),
			RunE: func(cmd *cobra.Command, args []string) error {
				return runConfigReplace(rtCopy, &flags)
			},
		}
		cmd.AddCommand(subCmd)
	}

	return cmd
}

// buildConfigStatusCmd creates the status subcommand
func buildConfigStatusCmd() *cobra.Command {
	var flags configurationFlags

	cmd := &cobra.Command{
		Use:     "status",
		Short:   "Status of configuration object",
		Example: "vesctl configuration status virtual_host <name>",
	}

	cmd.PersistentFlags().StringVar(&flags.atSite, "at-site", "", "Site name (e.g. ce01) at which to query configuration object status")
	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Namespace of configuration object")

	// Add resource type subcommands
	for _, rt := range types.All() {
		if !rt.Operations.Status {
			continue
		}
		rtCopy := rt
		subCmd := &cobra.Command{
			Use:   fmt.Sprintf("%s <name>", rt.Name),
			Short: fmt.Sprintf("Status of %s", rt.Name),
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				flags.name = args[0]
				return runConfigStatus(rtCopy, &flags)
			},
		}
		cmd.AddCommand(subCmd)
	}

	return cmd
}

// buildConfigApplyCmd creates the apply subcommand
func buildConfigApplyCmd() *cobra.Command {
	var flags configurationFlags

	cmd := &cobra.Command{
		Use:     "apply",
		Short:   "Apply (create or replace) configuration object",
		Example: "vesctl configuration apply virtual_host -i <file>",
	}

	cmd.PersistentFlags().StringVarP(&flags.inputFile, "input-file", "i", "", "File containing CreateRequest contents")
	cmd.PersistentFlags().StringVar(&flags.jsonData, "json-data", "", "Inline CreateRequest contents in json form")
	cmd.PersistentFlags().StringVar(&flags.mode, "mode", "always", "Either new(create fails if object exists) or always(object replaced if it exists)")

	// Add resource type subcommands
	for _, rt := range types.All() {
		if !rt.Operations.Create {
			continue
		}
		rtCopy := rt
		subCmd := &cobra.Command{
			Use:   rt.Name,
			Short: fmt.Sprintf("Apply %s", rt.Name),
			RunE: func(cmd *cobra.Command, args []string) error {
				return runConfigApply(rtCopy, &flags)
			},
		}
		cmd.AddCommand(subCmd)
	}

	return cmd
}

// buildConfigPatchCmd creates the patch subcommand
func buildConfigPatchCmd() *cobra.Command {
	var flags configurationFlags

	cmd := &cobra.Command{
		Use:     "patch",
		Short:   "Patch configuration object",
		Example: "vesctl configuration replace virtual_host add /metadata/description \"desc\"",
	}

	cmd.PersistentFlags().StringVar(&flags.name, "name", "", "Name of object")
	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Namespace of object")

	// Add resource type subcommands
	for _, rt := range types.All() {
		if !rt.Operations.Update {
			continue
		}
		rtCopy := rt
		subCmd := &cobra.Command{
			Use:   rt.Name,
			Short: fmt.Sprintf("Patch %s", rt.Name),
			RunE: func(cmd *cobra.Command, args []string) error {
				return runConfigPatch(rtCopy, &flags)
			},
		}
		cmd.AddCommand(subCmd)
	}

	return cmd
}

// buildConfigAddLabelsCmd creates the add-labels subcommand
func buildConfigAddLabelsCmd() *cobra.Command {
	var flags configurationFlags

	cmd := &cobra.Command{
		Use:     "add-labels",
		Short:   "Add Labels to a configuration object",
		Example: "vesctl configuration add-labels virtual_host <name> --label-key acmecorp.com/attr-1 --label-value val-1",
	}

	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Namespace of configuration object")
	cmd.PersistentFlags().StringSliceVar(&flags.labelKeys, "label-key", nil, "Key part of label")
	cmd.PersistentFlags().StringSliceVar(&flags.labelValues, "label-value", nil, "Value part of label")

	// Add resource type subcommands
	for _, rt := range types.All() {
		rtCopy := rt
		subCmd := &cobra.Command{
			Use:   fmt.Sprintf("%s <name>", rt.Name),
			Short: fmt.Sprintf("Add Labels to %s", rt.Name),
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				flags.name = args[0]
				return runConfigAddLabels(rtCopy, &flags)
			},
		}
		cmd.AddCommand(subCmd)
	}

	return cmd
}

// buildConfigRemoveLabelsCmd creates the remove-labels subcommand
func buildConfigRemoveLabelsCmd() *cobra.Command {
	var flags configurationFlags

	cmd := &cobra.Command{
		Use:     "remove-labels",
		Short:   "Remove Labels from a configuration object",
		Example: "vesctl configuration remove-labels virtual_host <name> --label-key acmecorp.com/attr-1",
	}

	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Namespace of configuration object")
	cmd.PersistentFlags().StringSliceVar(&flags.labelKeys, "label-key", nil, "Key part of label")

	// Add resource type subcommands
	for _, rt := range types.All() {
		rtCopy := rt
		subCmd := &cobra.Command{
			Use:   fmt.Sprintf("%s <name>", rt.Name),
			Short: fmt.Sprintf("Remove Labels from %s", rt.Name),
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				flags.name = args[0]
				return runConfigRemoveLabels(rtCopy, &flags)
			},
		}
		cmd.AddCommand(subCmd)
	}

	return cmd
}

// runConfigList executes the list operation (vesctl compatible)
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
		return fmt.Errorf("Error listing object: Listing object: %w", err)
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

// runConfigGet executes the get operation (vesctl compatible)
func runConfigGet(rt *types.ResourceType, flags *configurationFlags) error {
	// Note: We accept any response-format value including GET_RSP_FORMAT_READ
	// (original vesctl has a bug that rejects this valid value)

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
		return fmt.Errorf("Error getting object: Getting object: %w", err)
	}

	if resp.StatusCode >= 400 {
		return formatAPIError("getting", "GET", path, resp.StatusCode, resp.Body)
	}

	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Get defaults to YAML output (matching original vesctl)
	return output.Print(result, GetOutputFormatWithDefault("yaml"))
}

// runConfigCreate executes the create operation (vesctl compatible)
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

	path := rt.BuildAPIPath(namespace, "")
	resp, err := client.Post(ctx, path, resource)
	if err != nil {
		return fmt.Errorf("Error creating object: Creating object: %w", err)
	}

	if resp.StatusCode >= 400 {
		return formatAPIError("creating", "POST", path, resp.StatusCode, resp.Body)
	}

	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Print "Created" header (matching original vesctl)
	fmt.Println("Created")
	// Create defaults to YAML output (matching original vesctl)
	return output.Print(result, GetOutputFormatWithDefault("yaml"))
}

// runConfigDelete executes the delete operation (vesctl compatible)
func runConfigDelete(rt *types.ResourceType, flags *configurationFlags) error {
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
				return fmt.Errorf("Error deleting object: Deleting object: %w", err)
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
		return fmt.Errorf("Error deleting object: Deleting object: %w", err)
	}

	if resp.StatusCode >= 400 {
		return formatAPIError("deleting", "DELETE", path, resp.StatusCode, resp.Body)
	}

	output.PrintInfo(fmt.Sprintf("Deleted %s '%s' successfully", rt.Name, flags.name))
	return nil
}

// runConfigReplace executes the replace operation (vesctl compatible)
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

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	path := rt.BuildAPIPath(namespace, name)
	resp, err := client.Put(ctx, path, resource)
	if err != nil {
		return fmt.Errorf("Error replacing object: Replacing object: %w", err)
	}

	if resp.StatusCode >= 400 {
		return formatAPIError("replacing", "PUT", path, resp.StatusCode, resp.Body)
	}

	// Print only "Replaced" (matching original vesctl - no response body output)
	fmt.Println("Replaced")
	return nil
}

// runConfigStatus executes the status operation (vesctl compatible)
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
		return fmt.Errorf("Error getting status: Getting status: %w", err)
	}

	if resp.StatusCode >= 400 {
		return formatAPIError("getting status", "GET", path, resp.StatusCode, resp.Body)
	}

	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Status defaults to YAML output (matching original vesctl)
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

	// If mode is "new", only create (fail if exists)
	if flags.mode == "new" {
		path := rt.BuildAPIPath(namespace, "")
		resp, err := client.Post(ctx, path, resource)
		if err != nil {
			return fmt.Errorf("Error creating object: Creating object: %w", err)
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
				return fmt.Errorf("Error replacing object: Replacing object: %w", err)
			}
			if resp.StatusCode >= 400 {
				return formatAPIError("replacing", "PUT", getPath, resp.StatusCode, resp.Body)
			}
			// Print only "Replaced" (matching original vesctl - no response body output)
			fmt.Println("Replaced")
			return nil
		}
	}

	// Resource doesn't exist, create it
	createPath := rt.BuildAPIPath(namespace, "")
	resp, err := client.Post(ctx, createPath, resource)
	if err != nil {
		return fmt.Errorf("Error creating object: Creating object: %w", err)
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
		return fmt.Errorf("Error adding labels: Adding labels: %w", err)
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
		return fmt.Errorf("Error removing labels: Removing labels: %w", err)
	}

	if resp.StatusCode >= 400 {
		return formatAPIError("removing labels", "POST", path, resp.StatusCode, resp.Body)
	}

	output.PrintInfo(fmt.Sprintf("Removed labels from %s '%s'", rt.Name, flags.name))
	return nil
}

// formatAPIError formats an API error to match original vesctl error format
func formatAPIError(operation, method, path string, statusCode int, body []byte) error {
	baseURL := serverURLs[0]
	// Capitalize first letter of operation
	capOperation := strings.ToUpper(operation[:1]) + operation[1:]
	return fmt.Errorf("Error %s object: %s object: Unsuccessful %s at URL %s%s, status code %d, body %s, err %%!s(<nil>)",
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

