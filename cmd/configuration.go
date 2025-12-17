package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/robinmordasiewicz/f5xcctl/pkg/errors"
	"github.com/robinmordasiewicz/f5xcctl/pkg/naming"
	"github.com/robinmordasiewicz/f5xcctl/pkg/output"
	"github.com/robinmordasiewicz/f5xcctl/pkg/subscription"
	"github.com/robinmordasiewicz/f5xcctl/pkg/types"
)

// configurationFlags holds flags for configuration commands (f5xcctl compatibility)
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

// configurationCmd represents the configuration command (f5xcctl compatibility)
var configurationCmd = &cobra.Command{
	Use:     "configuration",
	Aliases: []string{"cfg", "c"},
	Short:   "Manage F5 XC configuration objects using CRUD operations.",
	Long:    `Manage F5 XC configuration objects using CRUD operations.`,
	Example: `f5xcctl configuration create virtual_host`,
}

func init() {
	rootCmd.AddCommand(configurationCmd)

	// Enable AI-agent-friendly error handling for invalid subcommands
	configurationCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("unknown command %q for %q\n\nUsage: f5xcctl configuration <action> [resource-type] [name] [flags]\n\nAvailable actions:\n  list, get, create, replace, apply, delete, status, patch, add-labels, remove-labels\n\nRun 'f5xcctl configuration --help' for usage", args[0], cmd.CommandPath())
		}
		return cmd.Help()
	}
	configurationCmd.SuggestionsMinimumDistance = 2

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
		Use:   "list",
		Short: "List all configuration objects of a specified type.",
		Long: `List all configuration objects of a specified type in F5 Distributed Cloud.

Returns a list of configurations with names, namespaces, and metadata.
Use --namespace to filter by namespace, or --output-format to control output format.`,
		Example: `  # List all http_loadbalancers in default namespace
  f5xcctl configuration list http_loadbalancer

  # List in a specific namespace
  f5xcctl configuration list http_loadbalancer -n production

  # List with JSON output
  f5xcctl configuration list http_loadbalancer --output-format json`,
	}

	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Target namespace for the operation.")

	// Add resource type subcommands
	for _, rt := range types.All() {
		rtCopy := rt

		// Build resource-specific Long description
		displayName := naming.ToHumanReadable(rt.Name)
		longDesc := fmt.Sprintf("List all %s resources in the specified namespace.\n\n", displayName)
		if rt.Description != "" {
			longDesc += rt.Description + "\n\n"
		}
		longDesc += "Returns a list of configurations with names, namespaces, and metadata."

		// Build resource-specific examples
		exampleText := fmt.Sprintf(`  # List all %s in default namespace
  f5xcctl configuration list %s

  # List %s in a specific namespace
  f5xcctl configuration list %s -n production

  # List with JSON output
  f5xcctl configuration list %s --output-format json`, rt.Name, rt.Name, rt.Name, rt.Name, rt.Name)

		subCmd := &cobra.Command{
			Use:     rt.Name,
			Short:   formatShortWithTier("List", displayName, rt.Name),
			Long:    longDesc,
			Example: exampleText,
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
		Use:   "get",
		Short: "Retrieve a specific configuration object by name.",
		Long: `Retrieve a specific configuration object by name from F5 Distributed Cloud.

Returns the full configuration including metadata and spec.
Use --response-format replace-request to get output suitable for editing and replacing.`,
		Example: `  # Get a specific http_loadbalancer
  f5xcctl configuration get http_loadbalancer example-lb

  # Get with replace-request format for editing
  f5xcctl configuration get http_loadbalancer example-lb --response-format replace-request

  # Get from a specific namespace
  f5xcctl configuration get http_loadbalancer example-lb -n production`,
	}

	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Target namespace for the operation.")
	cmd.PersistentFlags().StringVar(&flags.responseFormat, "response-format", "read", "Response format: 'read' for display or 'replace-request' for editing.")

	// Add resource type subcommands
	for _, rt := range types.All() {
		rtCopy := rt

		// Build resource-specific Long description
		displayName := naming.ToHumanReadable(rt.Name)
		longDesc := fmt.Sprintf("Retrieve a specific %s configuration by name.\n\n", displayName)
		if rt.Description != "" {
			longDesc += rt.Description + "\n\n"
		}
		longDesc += "Returns the full configuration including metadata, spec, and system metadata."

		// Build resource-specific examples
		exampleText := fmt.Sprintf(`  # Get a specific %s
  f5xcctl configuration get %s example-resource

  # Get with replace-request format for editing
  f5xcctl configuration get %s example-resource --response-format replace-request

  # Get from a specific namespace
  f5xcctl configuration get %s example-resource -n production`, rt.Name, rt.Name, rt.Name, rt.Name)

		subCmd := &cobra.Command{
			Use:     fmt.Sprintf("%s <name>", rt.Name),
			Short:   formatShortWithTier("Get", displayName, rt.Name),
			Long:    longDesc,
			Example: exampleText,
			Args:    cobra.ExactArgs(1),
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
		Short:   "Create a new configuration object from a YAML or JSON file.",
		Example: "f5xcctl configuration create virtual_host -i <file>",
	}

	cmd.PersistentFlags().StringVarP(&flags.inputFile, "input-file", "i", "", "Path to YAML or JSON file containing the resource definition.")
	cmd.PersistentFlags().StringVar(&flags.jsonData, "json-data", "", "Inline JSON string containing the resource definition.")

	// Add resource type subcommands
	for _, rt := range types.All() {
		if !rt.Operations.Create {
			continue
		}
		rtCopy := rt
		displayName := naming.ToHumanReadable(rt.Name)

		// Build example text
		exampleText := fmt.Sprintf(`# Create from file
f5xcctl configuration create %s -i config.yaml`, rt.Name)

		// Add inline JSON example if available
		if jsonExample := types.GetResourceExample(rt.Name); jsonExample != "" {
			exampleText += fmt.Sprintf(`

# Create with inline JSON using heredoc
f5xcctl configuration create %s --json-data "$(cat <<'EOF'
%s
EOF
)"`, rt.Name, jsonExample)
		}

		subCmd := &cobra.Command{
			Use:     rt.Name,
			Short:   formatShortWithTier("Create", displayName, rt.Name),
			Example: exampleText,
			RunE: func(cmd *cobra.Command, args []string) error {
				// Pre-validate subscription before creating resource
				if err := validateSubscriptionForResource(cmd.Context(), rtCopy.Name); err != nil {
					return err
				}
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
		Use:   "delete",
		Short: "Delete a configuration object by name.",
		Long: `Delete a configuration object by name from F5 Distributed Cloud.

This is a destructive operation. Use --yes to skip confirmation prompts.
In non-interactive mode (scripts, CI/CD), --yes is required.`,
		Example: `  # Delete an http_loadbalancer (with confirmation prompt)
  f5xcctl configuration delete http_loadbalancer example-lb

  # Delete without confirmation
  f5xcctl configuration delete http_loadbalancer example-lb --yes

  # Delete from a specific namespace
  f5xcctl configuration delete http_loadbalancer example-lb -n production --yes`,
	}

	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Target namespace for the operation.")
	cmd.PersistentFlags().BoolVarP(&flags.yes, "yes", "y", false, "Skip confirmation prompts for destructive operations.")

	// Add resource type subcommands
	for _, rt := range types.All() {
		if !rt.Operations.Delete {
			continue
		}
		rtCopy := rt
		displayName := naming.ToHumanReadable(rt.Name)

		// Build resource-specific Long description
		longDesc := fmt.Sprintf("Delete a %s configuration by name.\n\n", displayName)
		if rt.Description != "" {
			longDesc += rt.Description + "\n\n"
		}
		longDesc += "This is a destructive operation. Use --yes to skip confirmation prompts."

		// Build resource-specific examples
		exampleText := fmt.Sprintf(`  # Delete a %s (with confirmation prompt)
  f5xcctl configuration delete %s example-resource

  # Delete without confirmation (for scripts/CI)
  f5xcctl configuration delete %s example-resource --yes

  # Delete from a specific namespace
  f5xcctl configuration delete %s example-resource -n production --yes`, rt.Name, rt.Name, rt.Name, rt.Name)

		subCmd := &cobra.Command{
			Use:     fmt.Sprintf("%s <name>", rt.Name),
			Short:   formatShortWithTier("Delete", displayName, rt.Name),
			Long:    longDesc,
			Example: exampleText,
			Args:    cobra.ExactArgs(1),
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
		Use:   "replace",
		Short: "Replace an existing configuration object with new content.",
		Long: `Replace an existing configuration object with new content in F5 Distributed Cloud.

The resource must already exist. Use --input-file or --json-data to provide the new configuration.
This is a destructive operation. Use --yes to skip confirmation prompts.`,
		Example: `  # Replace from YAML file
  f5xcctl configuration replace http_loadbalancer -i config.yaml --yes

  # Replace with inline JSON
  f5xcctl configuration replace http_loadbalancer --json-data '{"metadata":{"name":"example-lb"},...}' --yes`,
	}

	cmd.PersistentFlags().StringVarP(&flags.inputFile, "input-file", "i", "", "Path to YAML or JSON file containing the resource definition.")
	cmd.PersistentFlags().StringVar(&flags.jsonData, "json-data", "", "Inline JSON string containing the resource definition.")
	cmd.PersistentFlags().BoolVarP(&flags.yes, "yes", "y", false, "Skip confirmation prompts for destructive operations.")

	// Add resource type subcommands
	for _, rt := range types.All() {
		if !rt.Operations.Update {
			continue
		}
		rtCopy := rt
		displayName := naming.ToHumanReadable(rt.Name)

		// Build resource-specific Long description
		longDesc := fmt.Sprintf("Replace an existing %s configuration with new content.\n\n", displayName)
		if rt.Description != "" {
			longDesc += rt.Description + "\n\n"
		}
		longDesc += "The resource must already exist. This is a destructive operation."

		// Build example text
		exampleText := fmt.Sprintf(`  # Replace from file
  f5xcctl configuration replace %s -i config.yaml --yes`, rt.Name)

		// Add inline JSON example if available
		if jsonExample := types.GetResourceExample(rt.Name); jsonExample != "" {
			exampleText += fmt.Sprintf(`

  # Replace with inline JSON using heredoc
  f5xcctl configuration replace %s --json-data "$(cat <<'EOF'
%s
EOF
)" --yes`, rt.Name, jsonExample)
		}

		subCmd := &cobra.Command{
			Use:     rt.Name,
			Short:   formatShortWithTier("Replace", displayName, rt.Name),
			Long:    longDesc,
			Example: exampleText,
			RunE: func(cmd *cobra.Command, args []string) error {
				// Pre-validate subscription before replacing resource
				if err := validateSubscriptionForResource(cmd.Context(), rtCopy.Name); err != nil {
					return err
				}
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
		Use:   "status",
		Short: "Display the current status of a configuration object.",
		Long: `Display the current status of a configuration object in F5 Distributed Cloud.

Returns the runtime status including deployment state, validation status, and any errors.
Use --at-site to query status at a specific site.`,
		Example: `  # Get status of an http_loadbalancer
  f5xcctl configuration status http_loadbalancer example-lb

  # Get status at a specific site
  f5xcctl configuration status http_loadbalancer example-lb --at-site example-site

  # Get status from a specific namespace
  f5xcctl configuration status http_loadbalancer example-lb -n production`,
	}

	cmd.PersistentFlags().StringVar(&flags.atSite, "at-site", "", "Site name to query for object status.")
	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Target namespace for the operation.")

	// Add resource type subcommands
	for _, rt := range types.All() {
		if !rt.Operations.Status {
			continue
		}
		rtCopy := rt
		displayName := naming.ToHumanReadable(rt.Name)

		// Build resource-specific Long description
		longDesc := fmt.Sprintf("Display the current status of a %s configuration.\n\n", displayName)
		if rt.Description != "" {
			longDesc += rt.Description + "\n\n"
		}
		longDesc += "Returns the runtime status including deployment state, validation status, and any errors."

		// Build resource-specific examples
		exampleText := fmt.Sprintf(`  # Get status of a %s
  f5xcctl configuration status %s example-resource

  # Get status at a specific site
  f5xcctl configuration status %s example-resource --at-site example-site

  # Get status from a specific namespace
  f5xcctl configuration status %s example-resource -n production`, rt.Name, rt.Name, rt.Name, rt.Name)

		subCmd := &cobra.Command{
			Use:     fmt.Sprintf("%s <name>", rt.Name),
			Short:   formatShortWithTier("Status of", displayName, rt.Name),
			Long:    longDesc,
			Example: exampleText,
			Args:    cobra.ExactArgs(1),
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
		Use:   "apply",
		Short: "Create or replace a configuration object using declarative input.",
		Long: `Create or replace a configuration object using declarative input in F5 Distributed Cloud.

Apply is idempotent - it creates the resource if it doesn't exist, or replaces it if it does.
Use --mode new to fail if the resource already exists (strict create behavior).`,
		Example: `  # Apply from YAML file (create or replace)
  f5xcctl configuration apply http_loadbalancer -i config.yaml

  # Apply with strict create mode (fail if exists)
  f5xcctl configuration apply http_loadbalancer -i config.yaml --mode new`,
	}

	cmd.PersistentFlags().StringVarP(&flags.inputFile, "input-file", "i", "", "Path to YAML or JSON file containing the resource definition.")
	cmd.PersistentFlags().StringVar(&flags.jsonData, "json-data", "", "Inline JSON string containing the resource definition.")
	cmd.PersistentFlags().StringVar(&flags.mode, "mode", "always", "Apply mode: 'new' to fail if exists, 'always' to create or replace.")

	// Add resource type subcommands
	for _, rt := range types.All() {
		if !rt.Operations.Create {
			continue
		}
		rtCopy := rt
		displayName := naming.ToHumanReadable(rt.Name)

		// Build resource-specific Long description
		longDesc := fmt.Sprintf("Create or replace a %s configuration using declarative input.\n\n", displayName)
		if rt.Description != "" {
			longDesc += rt.Description + "\n\n"
		}
		longDesc += "Apply is idempotent - creates if not exists, replaces if it does."

		// Build example text
		exampleText := fmt.Sprintf(`  # Apply from file (create or replace)
  f5xcctl configuration apply %s -i config.yaml

  # Apply with strict create mode
  f5xcctl configuration apply %s -i config.yaml --mode new`, rt.Name, rt.Name)

		// Add inline JSON example if available
		if jsonExample := types.GetResourceExample(rt.Name); jsonExample != "" {
			exampleText += fmt.Sprintf(`

  # Apply with inline JSON using heredoc
  f5xcctl configuration apply %s --json-data "$(cat <<'EOF'
%s
EOF
)"`, rt.Name, jsonExample)
		}

		subCmd := &cobra.Command{
			Use:     rt.Name,
			Short:   formatShortWithTier("Apply", displayName, rt.Name),
			Long:    longDesc,
			Example: exampleText,
			RunE: func(cmd *cobra.Command, args []string) error {
				// Pre-validate subscription before applying resource
				if err := validateSubscriptionForResource(cmd.Context(), rtCopy.Name); err != nil {
					return err
				}
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
		Use:   "patch",
		Short: "Apply a partial update to a configuration object.",
		Long: `Apply a partial update to a configuration object in F5 Distributed Cloud.

Note: Patch operation is not yet fully implemented. Use replace for complete updates.`,
		Example: `  # Patch is not yet implemented - use replace instead
  f5xcctl configuration replace http_loadbalancer -i updated-config.yaml`,
	}

	cmd.PersistentFlags().StringVar(&flags.name, "name", "", "Name of the target configuration object.")
	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Target namespace for the operation.")

	// Add resource type subcommands
	for _, rt := range types.All() {
		if !rt.Operations.Update {
			continue
		}
		rtCopy := rt
		displayName := naming.ToHumanReadable(rt.Name)

		// Build resource-specific Long description
		longDesc := fmt.Sprintf("Apply a partial update to a %s configuration.\n\n", displayName)
		if rt.Description != "" {
			longDesc += rt.Description + "\n\n"
		}
		longDesc += "Note: Patch operation is not yet fully implemented. Use replace for complete updates."

		subCmd := &cobra.Command{
			Use:   rt.Name,
			Short: fmt.Sprintf("Patch %s", displayName),
			Long:  longDesc,
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
		Use:   "add-labels",
		Short: "Add metadata labels to a configuration object.",
		Long: `Add metadata labels to a configuration object in F5 Distributed Cloud.

Labels are key-value pairs used for organizing and selecting resources.
Use --label-key and --label-value flags (can be repeated for multiple labels).`,
		Example: `  # Add a single label
  f5xcctl configuration add-labels http_loadbalancer example-lb --label-key env --label-value production

  # Add multiple labels
  f5xcctl configuration add-labels http_loadbalancer example-lb \
    --label-key env --label-value production \
    --label-key team --label-value platform`,
	}

	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Target namespace for the operation.")
	cmd.PersistentFlags().StringSliceVar(&flags.labelKeys, "label-key", nil, "Label key to add or remove.")
	cmd.PersistentFlags().StringSliceVar(&flags.labelValues, "label-value", nil, "Label value to assign.")

	// Add resource type subcommands
	for _, rt := range types.All() {
		rtCopy := rt
		displayName := naming.ToHumanReadable(rt.Name)

		// Build resource-specific Long description
		longDesc := fmt.Sprintf("Add metadata labels to a %s configuration.\n\n", displayName)
		if rt.Description != "" {
			longDesc += rt.Description + "\n\n"
		}
		longDesc += "Labels are key-value pairs used for organizing and selecting resources."

		// Build resource-specific examples
		exampleText := fmt.Sprintf(`  # Add a single label
  f5xcctl configuration add-labels %s example-resource --label-key env --label-value production

  # Add multiple labels
  f5xcctl configuration add-labels %s example-resource \
    --label-key env --label-value production \
    --label-key team --label-value platform`, rt.Name, rt.Name)

		subCmd := &cobra.Command{
			Use:     fmt.Sprintf("%s <name>", rt.Name),
			Short:   fmt.Sprintf("Add Labels to %s", displayName),
			Long:    longDesc,
			Example: exampleText,
			Args:    cobra.ExactArgs(1),
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
		Use:   "remove-labels",
		Short: "Remove metadata labels from a configuration object.",
		Long: `Remove metadata labels from a configuration object in F5 Distributed Cloud.

Specify the label keys to remove using --label-key flags (can be repeated for multiple labels).`,
		Example: `  # Remove a single label
  f5xcctl configuration remove-labels http_loadbalancer example-lb --label-key env

  # Remove multiple labels
  f5xcctl configuration remove-labels http_loadbalancer example-lb \
    --label-key env \
    --label-key team`,
	}

	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Target namespace for the operation.")
	cmd.PersistentFlags().StringSliceVar(&flags.labelKeys, "label-key", nil, "Label key to add or remove.")

	// Add resource type subcommands
	for _, rt := range types.All() {
		rtCopy := rt
		displayName := naming.ToHumanReadable(rt.Name)

		// Build resource-specific Long description
		longDesc := fmt.Sprintf("Remove metadata labels from a %s configuration.\n\n", displayName)
		if rt.Description != "" {
			longDesc += rt.Description + "\n\n"
		}
		longDesc += "Specify the label keys to remove using --label-key flags."

		// Build resource-specific examples
		exampleText := fmt.Sprintf(`  # Remove a single label
  f5xcctl configuration remove-labels %s example-resource --label-key env

  # Remove multiple labels
  f5xcctl configuration remove-labels %s example-resource \
    --label-key env \
    --label-key team`, rt.Name, rt.Name)

		subCmd := &cobra.Command{
			Use:     fmt.Sprintf("%s <name>", rt.Name),
			Short:   fmt.Sprintf("Remove Labels from %s", displayName),
			Long:    longDesc,
			Example: exampleText,
			Args:    cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				flags.name = args[0]
				return runConfigRemoveLabels(rtCopy, &flags)
			},
		}
		cmd.AddCommand(subCmd)
	}

	return cmd
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

// validateSubscriptionForResource validates that the current subscription allows
// creating/modifying the specified resource type. Returns an error with exit code 9
// (ExitFeatureNotAvail) if the resource requires a higher subscription tier.
func validateSubscriptionForResource(ctx context.Context, resourceType string) error {
	validator := GetSubscriptionValidator()
	if validator == nil {
		// No validator available, allow the operation
		return nil
	}

	result, err := validator.ValidateResourceAccess(ctx, resourceType)
	if err != nil {
		// Validation error, log warning but don't block
		if IsDebug() {
			fmt.Fprintf(os.Stderr, "Warning: subscription validation failed: %v\n", err)
		}
		return nil
	}

	if !result.IsAllowed {
		// Resource not allowed by subscription
		errMsg := result.ErrorMessage
		if result.Recommendation != "" {
			errMsg += "\n\nTo resolve:\n  - " + result.Recommendation
		}
		errMsg += "\n\nFor more information:\n  f5xcctl subscription show      # View current subscription\n  f5xcctl subscription addons    # View addon services"

		return errors.NewExitError(errors.ExitFeatureNotAvail, errors.ErrFeatureNotAvail, errMsg)
	}

	return nil
}
