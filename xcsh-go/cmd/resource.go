package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/robinmordasiewicz/xcsh/pkg/naming"
	"github.com/robinmordasiewicz/xcsh/pkg/output"
	"github.com/robinmordasiewicz/xcsh/pkg/types"
)

// resourceFlags holds common flags for resource commands
type resourceFlags struct {
	namespace string
	name      string
	file      string
	yes       bool
}

// BuildResourceCommand creates a command group for a resource type
func BuildResourceCommand(rt *types.ResourceType) *cobra.Command {
	displayName := naming.ToHumanReadable(rt.Name)
	cmd := &cobra.Command{
		Use:   rt.CLIName,
		Short: fmt.Sprintf("Manage %s resources", displayName),
		Long:  fmt.Sprintf("Manage %s resources in F5 Distributed Cloud.", displayName),
	}

	// Add subcommands based on supported operations
	if rt.Operations.List {
		cmd.AddCommand(buildListCommand(rt))
	}
	if rt.Operations.Get {
		cmd.AddCommand(buildShowCommand(rt))
	}
	if rt.Operations.Create {
		cmd.AddCommand(buildCreateCommand(rt))
	}
	if rt.Operations.Update {
		cmd.AddCommand(buildUpdateCommand(rt))
	}
	if rt.Operations.Delete {
		cmd.AddCommand(buildDeleteCommand(rt))
	}
	if rt.Operations.Status {
		cmd.AddCommand(buildStatusCommand(rt))
	}

	return cmd
}

// buildListCommand creates the list subcommand
func buildListCommand(rt *types.ResourceType) *cobra.Command {
	var flags resourceFlags
	displayName := naming.ToHumanReadable(rt.Name)

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   fmt.Sprintf("List %s resources", displayName),
		Long:    fmt.Sprintf("List all %s resources in the specified namespace.", displayName),
		Example: buildExample(rt, "list", []string{
			fmt.Sprintf("xcsh %s list --namespace example-namespace", rt.CLIName),
			fmt.Sprintf("xcsh %s list -n example-namespace -o table", rt.CLIName),
		}),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(rt, &flags)
		},
	}

	if rt.SupportsNamespace {
		cmd.Flags().StringVarP(&flags.namespace, "namespace", "n", GetDefaultNamespace(), "Namespace")
	}

	return cmd
}

// buildShowCommand creates the show subcommand
func buildShowCommand(rt *types.ResourceType) *cobra.Command {
	var flags resourceFlags
	displayName := naming.ToHumanReadable(rt.Name)

	cmd := &cobra.Command{
		Use:     "show [name]",
		Aliases: []string{"get", "describe"},
		Short:   fmt.Sprintf("Show details of a %s", displayName),
		Long:    fmt.Sprintf("Display detailed information about a specific %s.", displayName),
		Args:    cobra.ExactArgs(1),
		Example: buildExample(rt, "show", []string{
			fmt.Sprintf("xcsh %s show example-resource --namespace example-namespace", rt.CLIName),
			fmt.Sprintf("xcsh %s show example-resource -n example-namespace -o json", rt.CLIName),
		}),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.name = args[0]
			return runShow(rt, &flags)
		},
	}

	if rt.SupportsNamespace {
		cmd.Flags().StringVarP(&flags.namespace, "namespace", "n", GetDefaultNamespace(), "Namespace")
	}

	return cmd
}

// buildCreateCommand creates the create subcommand
func buildCreateCommand(rt *types.ResourceType) *cobra.Command {
	var flags resourceFlags
	displayName := naming.ToHumanReadable(rt.Name)

	cmd := &cobra.Command{
		Use:   "create",
		Short: fmt.Sprintf("Create a new %s", displayName),
		Long:  fmt.Sprintf("Create a new %s from a YAML or JSON file.", displayName),
		Example: buildExample(rt, "create", []string{
			fmt.Sprintf("xcsh %s create --file resource.yaml", rt.CLIName),
			fmt.Sprintf("xcsh %s create -f resource.json --namespace example-namespace", rt.CLIName),
		}),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(rt, &flags)
		},
	}

	cmd.Flags().StringVarP(&flags.file, "file", "f", "", "Path to resource definition file (YAML or JSON)")
	cmd.Flags().StringVarP(&flags.file, "input-file", "i", "", "Path to resource definition file (xcsh compatible)")
	_ = cmd.MarkFlagRequired("file")

	if rt.SupportsNamespace {
		cmd.Flags().StringVarP(&flags.namespace, "namespace", "n", GetDefaultNamespace(), "Namespace (overrides file)")
	}

	return cmd
}

// buildUpdateCommand creates the update subcommand
func buildUpdateCommand(rt *types.ResourceType) *cobra.Command {
	var flags resourceFlags
	displayName := naming.ToHumanReadable(rt.Name)

	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"replace", "apply"},
		Short:   fmt.Sprintf("Update an existing %s", displayName),
		Long:    fmt.Sprintf("Update an existing %s from a YAML or JSON file.", displayName),
		Example: buildExample(rt, "update", []string{
			fmt.Sprintf("xcsh %s update --file resource.yaml", rt.CLIName),
			fmt.Sprintf("xcsh %s update -f resource.json", rt.CLIName),
		}),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(rt, &flags)
		},
	}

	cmd.Flags().StringVarP(&flags.file, "file", "f", "", "Path to resource definition file (YAML or JSON)")
	cmd.Flags().StringVarP(&flags.file, "input-file", "i", "", "Path to resource definition file (xcsh compatible)")
	_ = cmd.MarkFlagRequired("file")

	if rt.SupportsNamespace {
		cmd.Flags().StringVarP(&flags.namespace, "namespace", "n", GetDefaultNamespace(), "Namespace (overrides file)")
	}

	return cmd
}

// buildDeleteCommand creates the delete subcommand
func buildDeleteCommand(rt *types.ResourceType) *cobra.Command {
	var flags resourceFlags
	displayName := naming.ToHumanReadable(rt.Name)

	cmd := &cobra.Command{
		Use:     "delete [name]",
		Aliases: []string{"rm", "remove"},
		Short:   fmt.Sprintf("Delete a %s", displayName),
		Long:    fmt.Sprintf("Delete a %s resource.", displayName),
		Args:    cobra.ExactArgs(1),
		Example: buildExample(rt, "delete", []string{
			fmt.Sprintf("xcsh %s delete example-resource --namespace example-namespace", rt.CLIName),
			fmt.Sprintf("xcsh %s delete example-resource -n example-namespace --yes", rt.CLIName),
		}),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.name = args[0]
			return runDelete(rt, &flags)
		},
	}

	if rt.SupportsNamespace {
		cmd.Flags().StringVarP(&flags.namespace, "namespace", "n", GetDefaultNamespace(), "Namespace")
	}

	cmd.Flags().BoolVarP(&flags.yes, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}

// buildStatusCommand creates the status subcommand
func buildStatusCommand(rt *types.ResourceType) *cobra.Command {
	var flags resourceFlags
	displayName := naming.ToHumanReadable(rt.Name)

	cmd := &cobra.Command{
		Use:   "status [name]",
		Short: fmt.Sprintf("Show status of a %s", displayName),
		Long:  fmt.Sprintf("Display the operational status of a %s.", displayName),
		Args:  cobra.ExactArgs(1),
		Example: buildExample(rt, "status", []string{
			fmt.Sprintf("xcsh %s status example-resource --namespace example-namespace", rt.CLIName),
		}),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.name = args[0]
			return runStatus(rt, &flags)
		},
	}

	if rt.SupportsNamespace {
		cmd.Flags().StringVarP(&flags.namespace, "namespace", "n", GetDefaultNamespace(), "Namespace")
	}

	return cmd
}

// buildExample formats example strings
func buildExample(rt *types.ResourceType, operation string, examples []string) string {
	result := ""
	for i, ex := range examples {
		if i > 0 {
			result += "\n"
		}
		result += "  " + ex
	}
	return result
}

// runList executes the list operation
func runList(rt *types.ResourceType, flags *resourceFlags) error {
	client := GetClient()
	if client == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	path := rt.BuildAPIPath(flags.namespace, "")
	resp, err := client.Get(ctx, path, nil)
	if err != nil {
		return fmt.Errorf("failed to list resources: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error: %s", string(resp.Body))
	}

	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return output.Print(result, GetOutputFormat())
}

// runShow executes the show operation
func runShow(rt *types.ResourceType, flags *resourceFlags) error {
	client := GetClient()
	if client == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	path := rt.BuildAPIPath(flags.namespace, flags.name)
	resp, err := client.Get(ctx, path, nil)
	if err != nil {
		return fmt.Errorf("failed to get resource: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error: %s", string(resp.Body))
	}

	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return output.Print(result, GetOutputFormat())
}

// runCreate executes the create operation
func runCreate(rt *types.ResourceType, flags *resourceFlags) error {
	client := GetClient()
	if client == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	// Load resource from file
	resource, err := loadResourceFromFile(flags.file)
	if err != nil {
		return fmt.Errorf("failed to load resource file: %w", err)
	}

	// Override namespace if specified
	namespace := flags.namespace
	if namespace == "" {
		if meta, ok := resource["metadata"].(map[string]interface{}); ok {
			if ns, ok := meta["namespace"].(string); ok {
				namespace = ns
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	path := rt.BuildAPIPath(namespace, "")
	resp, err := client.Post(ctx, path, resource)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error: %s", string(resp.Body))
	}

	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	displayName := naming.ToHumanReadable(rt.Name)
	output.PrintInfo(fmt.Sprintf("Created %s successfully", displayName))
	return output.Print(result, GetOutputFormat())
}

// runUpdate executes the update operation
func runUpdate(rt *types.ResourceType, flags *resourceFlags) error {
	client := GetClient()
	if client == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	// Load resource from file
	resource, err := loadResourceFromFile(flags.file)
	if err != nil {
		return fmt.Errorf("failed to load resource file: %w", err)
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

	if flags.namespace != "" {
		namespace = flags.namespace
	}

	if name == "" {
		return fmt.Errorf("resource name is required in metadata")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	path := rt.BuildAPIPath(namespace, name)
	resp, err := client.Put(ctx, path, resource)
	if err != nil {
		return fmt.Errorf("failed to update resource: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error: %s", string(resp.Body))
	}

	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	displayName := naming.ToHumanReadable(rt.Name)
	output.PrintInfo(fmt.Sprintf("Updated %s '%s' successfully", displayName, name))
	return output.Print(result, GetOutputFormat())
}

// runDelete executes the delete operation
func runDelete(rt *types.ResourceType, flags *resourceFlags) error {
	client := GetClient()
	if client == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	displayName := naming.ToHumanReadable(rt.Name)

	// Confirm deletion unless --yes is specified
	if !flags.yes {
		fmt.Printf("Are you sure you want to delete %s '%s'? [y/N]: ", displayName, flags.name)
		var confirm string
		_, _ = fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" && confirm != "yes" {
			output.PrintInfo("Deletion cancelled")
			return nil
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	path := rt.BuildAPIPath(flags.namespace, flags.name)

	// Check if resource type has custom delete configuration
	if rt.DeleteConfig != nil {
		// Apply custom delete path suffix
		if rt.DeleteConfig.PathSuffix != "" {
			path = path + rt.DeleteConfig.PathSuffix
		}

		// Use custom HTTP method (typically POST for cascade_delete)
		if rt.DeleteConfig.Method == "POST" {
			var body interface{}
			if rt.DeleteConfig.IncludeBody {
				body = map[string]interface{}{
					"name": flags.name,
				}
			}
			resp, err := client.Post(ctx, path, body)
			if err != nil {
				return fmt.Errorf("failed to delete resource: %w", err)
			}
			if resp.StatusCode >= 400 {
				return fmt.Errorf("API error: %s", string(resp.Body))
			}
			output.PrintInfo(fmt.Sprintf("Deleted %s '%s' successfully", displayName, flags.name))
			return nil
		}
	}

	// Standard DELETE method
	resp, err := client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error: %s", string(resp.Body))
	}

	output.PrintInfo(fmt.Sprintf("Deleted %s '%s' successfully", displayName, flags.name))
	return nil
}

// runStatus executes the status operation
func runStatus(rt *types.ResourceType, flags *resourceFlags) error {
	client := GetClient()
	if client == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Status endpoint is typically the same as get with /status suffix
	path := rt.BuildAPIPath(flags.namespace, flags.name) + "/status"
	resp, err := client.Get(ctx, path, nil)
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error: %s", string(resp.Body))
	}

	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return output.Print(result, GetOutputFormat())
}

// loadResourceFromFile loads a resource definition from a YAML or JSON file
func loadResourceFromFile(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var resource map[string]interface{}

	// Try YAML first (YAML is a superset of JSON)
	if err := yaml.Unmarshal(data, &resource); err != nil {
		// Try JSON if YAML fails
		if err := json.Unmarshal(data, &resource); err != nil {
			return nil, fmt.Errorf("failed to parse file (not valid YAML or JSON): %w", err)
		}
	}

	return resource, nil
}

// customResourceCommands lists resource types that have custom implementations
// and should not be auto-registered
var customResourceCommands = map[string]bool{
	"site":            true, // Custom implementation in site.go
	"aws-vpc-site":    true, // Custom implementation in site_aws_vpc.go
	"azure-vnet-site": true, // Custom implementation in site_azure_vnet.go
}

// RegisterAllResourceCommands registers all resource type commands
func RegisterAllResourceCommands() {
	for _, rt := range types.All() {
		// Skip resources with custom implementations
		if customResourceCommands[rt.CLIName] {
			continue
		}
		rootCmd.AddCommand(BuildResourceCommand(rt))
	}
}

// init() removed - original xcsh does not have individual resource commands
// All resource operations go through: xcsh configuration <operation> <object_type>
// func init() {
// 	RegisterAllResourceCommands()
// }
