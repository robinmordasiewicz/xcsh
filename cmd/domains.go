package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/robinmordasiewicz/f5xcctl/pkg/naming"
	"github.com/robinmordasiewicz/f5xcctl/pkg/types"
)

// init registers all domain commands dynamically
func init() {
	// Register domain commands for all domains in DomainRegistry
	for domain := range types.DomainRegistry {
		rootCmd.AddCommand(buildDomainCmd(domain))
	}
}

// buildDomainCmd creates a top-level domain command with all operation subcommands
func buildDomainCmd(domain string) *cobra.Command {
	info, _ := types.GetDomainInfo(domain)

	cmd := &cobra.Command{
		Use:     domain,
		Aliases: info.Aliases,
		Short:   fmt.Sprintf("Manage %s resources", info.DisplayName),
		Long: fmt.Sprintf(`Manage F5 Distributed Cloud %s resources.

%s

OPERATIONS:
  list           List resources of a type (optionally filtered by namespace)
  get            Retrieve a specific resource by name
  create         Create a new resource from YAML/JSON file
  replace        Replace an existing resource completely
  apply          Create or update (upsert) a resource
  delete         Remove a resource by name
  status         Check the operational status of a resource
  patch          Partially update a resource
  add-labels     Add labels to a resource
  remove-labels  Remove labels from a resource

Run 'f5xcctl %s --help' for more information.`, info.DisplayName, info.Description, domain),
	}

	// Enable error handling for invalid subcommands
	cmd.RunE = func(c *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("unknown command %q for %q\n\nUsage: f5xcctl %s <action> [resource-type] [name] [flags]\n\nAvailable actions:\n  list, get, create, replace, apply, delete, status, patch, add-labels, remove-labels\n\nRun 'f5xcctl %s --help' for usage", args[0], c.CommandPath(), domain, domain)
		}
		return c.Help()
	}
	cmd.SuggestionsMinimumDistance = 2

	// Add operation subcommands for this domain
	cmd.AddCommand(buildDomainListCmd(domain))
	cmd.AddCommand(buildDomainGetCmd(domain))
	cmd.AddCommand(buildDomainCreateCmd(domain))
	cmd.AddCommand(buildDomainDeleteCmd(domain))
	cmd.AddCommand(buildDomainReplaceCmd(domain))
	cmd.AddCommand(buildDomainStatusCmd(domain))
	cmd.AddCommand(buildDomainApplyCmd(domain))
	cmd.AddCommand(buildDomainPatchCmd(domain))
	cmd.AddCommand(buildDomainAddLabelsCmd(domain))
	cmd.AddCommand(buildDomainRemoveLabelsCmd(domain))

	return cmd
}

// buildDomainListCmd creates the list operation for a domain
func buildDomainListCmd(domain string) *cobra.Command {
	var flags configurationFlags

	domainInfo, _ := types.GetDomainInfo(domain)

	cmd := &cobra.Command{
		Use:   "list",
		Short: fmt.Sprintf("List %s resources", domainInfo.DisplayName),
		Long: fmt.Sprintf(`List all %s resources in the specified namespace.

%s

Returns a list of configurations with names, namespaces, and metadata.
Use --namespace to filter by namespace, or --output-format to control output format.`, domainInfo.DisplayName, domainInfo.Description),
	}

	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Target namespace for the operation.")

	// Register flag completions
	_ = cmd.RegisterFlagCompletionFunc("namespace", completeNamespace)

	// Get resources for this domain (cross-domain enabled)
	resources := types.GetByDomain(domain)
	sortResourcesByName(resources)

	// Add resource type subcommands for this domain
	for _, rt := range resources {
		rtCopy := rt

		displayName := naming.ToHumanReadable(rt.Name)
		longDesc := fmt.Sprintf("List all %s resources in the specified namespace.\n\n", displayName)
		if rt.Description != "" {
			longDesc += rt.Description + "\n\n"
		}
		longDesc += "Returns a list of configurations with names, namespaces, and metadata."

		exampleText := fmt.Sprintf(`  # List all %s in default namespace
  f5xcctl %s list %s

  # List %s in a specific namespace
  f5xcctl %s list %s -n production

  # List with JSON output
  f5xcctl %s list %s --output-format json`, rt.Name, domain, rt.Name, rt.Name, domain, rt.Name, domain, rt.Name)

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

// buildDomainGetCmd creates the get operation for a domain
func buildDomainGetCmd(domain string) *cobra.Command {
	var flags configurationFlags

	domainInfo, _ := types.GetDomainInfo(domain)

	cmd := &cobra.Command{
		Use:   "get",
		Short: fmt.Sprintf("Retrieve a %s resource by name", domainInfo.DisplayName),
		Long: fmt.Sprintf(`Retrieve a specific %s resource by name from F5 Distributed Cloud.

%s

Returns the full configuration including metadata and spec.
Use --response-format replace-request to get output suitable for editing and replacing.`, domainInfo.DisplayName, domainInfo.Description),
	}

	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Target namespace for the operation.")
	cmd.PersistentFlags().StringVar(&flags.responseFormat, "response-format", "read", "Response format: 'read' for display or 'replace-request' for editing.")

	// Register flag completions
	_ = cmd.RegisterFlagCompletionFunc("namespace", completeNamespace)

	// Get resources for this domain
	resources := types.GetByDomain(domain)
	sortResourcesByName(resources)

	// Add resource type subcommands for this domain
	for _, rt := range resources {
		rtCopy := rt

		displayName := naming.ToHumanReadable(rt.Name)
		longDesc := fmt.Sprintf("Retrieve a specific %s resource by name.\n\n", displayName)
		if rt.Description != "" {
			longDesc += rt.Description + "\n\n"
		}
		longDesc += "Returns the full configuration including metadata, spec, and system metadata."

		exampleText := fmt.Sprintf(`  # Get a specific %s
  f5xcctl %s get %s example-resource

  # Get with replace-request format for editing
  f5xcctl %s get %s example-resource --response-format replace-request

  # Get from a specific namespace
  f5xcctl %s get %s example-resource -n production`, rt.Name, domain, rt.Name, domain, rt.Name, domain, rt.Name)

		subCmd := &cobra.Command{
			Use:               fmt.Sprintf("%s <name>", rt.Name),
			Short:             formatShortWithTier("Get", displayName, rt.Name),
			Long:              longDesc,
			Example:           exampleText,
			Args:              cobra.ExactArgs(1),
			ValidArgsFunction: completeResourceName(domain, rtCopy.Name),
			RunE: func(cmd *cobra.Command, args []string) error {
				flags.name = args[0]
				return runConfigGet(rtCopy, &flags)
			},
		}
		_ = subCmd.RegisterFlagCompletionFunc("namespace", completeNamespace)
		cmd.AddCommand(subCmd)
	}

	return cmd
}

// buildDomainCreateCmd creates the create operation for a domain
func buildDomainCreateCmd(domain string) *cobra.Command {
	var flags configurationFlags

	domainInfo, _ := types.GetDomainInfo(domain)

	cmd := &cobra.Command{
		Use:   "create",
		Short: fmt.Sprintf("Create a new %s resource", domainInfo.DisplayName),
		Long: fmt.Sprintf(`Create a new %s resource in F5 Distributed Cloud.

%s

Provide a YAML or JSON file with the resource configuration using --input-file.`, domainInfo.DisplayName, domainInfo.Description),
	}

	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Target namespace for the operation.")
	cmd.PersistentFlags().StringVarP(&flags.inputFile, "input-file", "i", "", "Path to YAML/JSON file containing the resource configuration.")
	cmd.PersistentFlags().StringVar(&flags.mode, "mode", "raw", "Input mode: 'raw' for direct config or 'form' for interactive mode.")

	// Register flag completions
	_ = cmd.RegisterFlagCompletionFunc("namespace", completeNamespace)

	// Get resources for this domain
	resources := types.GetByDomain(domain)
	sortResourcesByName(resources)

	// Add resource type subcommands for this domain
	for _, rt := range resources {
		if !rt.Operations.Create {
			continue // Skip resources that don't support create
		}

		rtCopy := rt

		displayName := naming.ToHumanReadable(rt.Name)
		longDesc := fmt.Sprintf("Create a new %s resource.\n\n", displayName)
		if rt.Description != "" {
			longDesc += rt.Description + "\n\n"
		}
		longDesc += "Provide a YAML or JSON file with the resource configuration."

		exampleText := fmt.Sprintf(`  # Create from file
  f5xcctl %s create %s -n my-namespace -i config.yaml

  # Create with JSON input
  f5xcctl %s create %s -n my-namespace -i config.json`, domain, rt.Name, domain, rt.Name)

		subCmd := &cobra.Command{
			Use:     rt.Name,
			Short:   formatShortWithTier("Create", displayName, rt.Name),
			Long:    longDesc,
			Example: exampleText,
			RunE: func(cmd *cobra.Command, args []string) error {
				return runConfigCreate(rtCopy, &flags)
			},
		}
		_ = subCmd.RegisterFlagCompletionFunc("namespace", completeNamespace)
		cmd.AddCommand(subCmd)
	}

	return cmd
}

// buildDomainDeleteCmd creates the delete operation for a domain
func buildDomainDeleteCmd(domain string) *cobra.Command {
	var flags configurationFlags

	domainInfo, _ := types.GetDomainInfo(domain)

	cmd := &cobra.Command{
		Use:   "delete",
		Short: fmt.Sprintf("Delete a %s resource", domainInfo.DisplayName),
		Long: fmt.Sprintf(`Delete a %s resource from F5 Distributed Cloud.

%s

Requires confirmation unless --yes is specified.`, domainInfo.DisplayName, domainInfo.Description),
	}

	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Target namespace for the operation.")
	cmd.PersistentFlags().BoolVar(&flags.yes, "yes", false, "Skip confirmation prompt.")

	// Register flag completions
	_ = cmd.RegisterFlagCompletionFunc("namespace", completeNamespace)

	// Get resources for this domain
	resources := types.GetByDomain(domain)
	sortResourcesByName(resources)

	// Add resource type subcommands for this domain
	for _, rt := range resources {
		if !rt.Operations.Delete {
			continue // Skip resources that don't support delete
		}

		rtCopy := rt

		displayName := naming.ToHumanReadable(rt.Name)
		longDesc := fmt.Sprintf("Delete a %s resource by name.\n\n", displayName)
		if rt.Description != "" {
			longDesc += rt.Description + "\n\n"
		}
		longDesc += "Requires confirmation unless --yes is specified."

		exampleText := fmt.Sprintf(`  # Delete a resource (with confirmation)
  f5xcctl %s delete %s -n my-namespace my-resource

  # Delete without confirmation
  f5xcctl %s delete %s -n my-namespace my-resource --yes`, domain, rt.Name, domain, rt.Name)

		subCmd := &cobra.Command{
			Use:               fmt.Sprintf("%s <name>", rt.Name),
			Short:             formatShortWithTier("Delete", displayName, rt.Name),
			Long:              longDesc,
			Example:           exampleText,
			Args:              cobra.ExactArgs(1),
			ValidArgsFunction: completeResourceName(domain, rtCopy.Name),
			RunE: func(cmd *cobra.Command, args []string) error {
				flags.name = args[0]
				return runConfigDelete(rtCopy, &flags)
			},
		}
		_ = subCmd.RegisterFlagCompletionFunc("namespace", completeNamespace)
		cmd.AddCommand(subCmd)
	}

	return cmd
}

// buildDomainReplaceCmd creates the replace operation for a domain
func buildDomainReplaceCmd(domain string) *cobra.Command {
	var flags configurationFlags

	domainInfo, _ := types.GetDomainInfo(domain)

	cmd := &cobra.Command{
		Use:   "replace",
		Short: fmt.Sprintf("Replace a %s resource completely", domainInfo.DisplayName),
		Long: fmt.Sprintf(`Replace a %s resource completely in F5 Distributed Cloud.

%s

This performs a complete replacement of the resource with the provided configuration.
Use apply for create-or-update semantics.`, domainInfo.DisplayName, domainInfo.Description),
	}

	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Target namespace for the operation.")
	cmd.PersistentFlags().StringVarP(&flags.inputFile, "input-file", "i", "", "Path to YAML/JSON file containing the updated configuration.")

	// Register flag completions
	_ = cmd.RegisterFlagCompletionFunc("namespace", completeNamespace)

	// Get resources for this domain
	resources := types.GetByDomain(domain)
	sortResourcesByName(resources)

	// Add resource type subcommands for this domain
	for _, rt := range resources {
		if !rt.Operations.Update {
			continue // Skip resources that don't support update
		}

		rtCopy := rt

		displayName := naming.ToHumanReadable(rt.Name)
		longDesc := fmt.Sprintf("Replace a %s resource.\n\n", displayName)
		if rt.Description != "" {
			longDesc += rt.Description + "\n\n"
		}
		longDesc += "This performs a complete replacement of the resource."

		exampleText := fmt.Sprintf(`  # Replace from file
  f5xcctl %s replace %s -n my-namespace my-resource -i updated-config.yaml`, domain, rt.Name)

		subCmd := &cobra.Command{
			Use:     fmt.Sprintf("%s <name>", rt.Name),
			Short:   formatShortWithTier("Replace", displayName, rt.Name),
			Long:    longDesc,
			Example: exampleText,
			Args:    cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				flags.name = args[0]
				return runConfigReplace(rtCopy, &flags)
			},
		}
		_ = subCmd.RegisterFlagCompletionFunc("namespace", completeNamespace)
		cmd.AddCommand(subCmd)
	}

	return cmd
}

// buildDomainStatusCmd creates the status operation for a domain
func buildDomainStatusCmd(domain string) *cobra.Command {
	var flags configurationFlags

	domainInfo, _ := types.GetDomainInfo(domain)

	cmd := &cobra.Command{
		Use:   "status",
		Short: fmt.Sprintf("Check the operational status of a %s resource", domainInfo.DisplayName),
		Long: fmt.Sprintf(`Check the operational status of a %s resource.

%s

Returns the current operational state and any relevant status information.`, domainInfo.DisplayName, domainInfo.Description),
	}

	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Target namespace for the operation.")

	// Register flag completions
	_ = cmd.RegisterFlagCompletionFunc("namespace", completeNamespace)

	// Get resources for this domain
	resources := types.GetByDomain(domain)
	sortResourcesByName(resources)

	// Add resource type subcommands for this domain
	for _, rt := range resources {
		if !rt.Operations.Status {
			continue // Skip resources that don't support status
		}

		rtCopy := rt

		displayName := naming.ToHumanReadable(rt.Name)
		longDesc := fmt.Sprintf("Check the operational status of a %s resource.\n\n", displayName)
		if rt.Description != "" {
			longDesc += rt.Description + "\n\n"
		}

		exampleText := fmt.Sprintf(`  # Check status
  f5xcctl %s status %s -n my-namespace my-resource`, domain, rt.Name)

		subCmd := &cobra.Command{
			Use:               fmt.Sprintf("%s <name>", rt.Name),
			Short:             formatShortWithTier("Status", displayName, rt.Name),
			Long:              longDesc,
			Example:           exampleText,
			Args:              cobra.ExactArgs(1),
			ValidArgsFunction: completeResourceName(domain, rtCopy.Name),
			RunE: func(cmd *cobra.Command, args []string) error {
				flags.name = args[0]
				return runConfigStatus(rtCopy, &flags)
			},
		}
		_ = subCmd.RegisterFlagCompletionFunc("namespace", completeNamespace)
		cmd.AddCommand(subCmd)
	}

	return cmd
}

// buildDomainApplyCmd creates the apply operation for a domain
func buildDomainApplyCmd(domain string) *cobra.Command {
	var flags configurationFlags

	domainInfo, _ := types.GetDomainInfo(domain)

	cmd := &cobra.Command{
		Use:   "apply",
		Short: fmt.Sprintf("Create or update (upsert) a %s resource", domainInfo.DisplayName),
		Long: fmt.Sprintf(`Create or update (upsert) a %s resource in F5 Distributed Cloud.

%s

If the resource exists, it will be updated. If it doesn't exist, it will be created.`, domainInfo.DisplayName, domainInfo.Description),
	}

	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Target namespace for the operation.")
	cmd.PersistentFlags().StringVarP(&flags.inputFile, "input-file", "i", "", "Path to YAML/JSON file containing the resource configuration.")

	// Register flag completions
	_ = cmd.RegisterFlagCompletionFunc("namespace", completeNamespace)

	// Get resources for this domain
	resources := types.GetByDomain(domain)
	sortResourcesByName(resources)

	// Add resource type subcommands for this domain
	for _, rt := range resources {
		if !rt.Operations.Create && !rt.Operations.Update {
			continue // Skip resources that don't support create or update
		}

		rtCopy := rt

		displayName := naming.ToHumanReadable(rt.Name)
		longDesc := fmt.Sprintf("Create or update (upsert) a %s resource.\n\n", displayName)
		if rt.Description != "" {
			longDesc += rt.Description + "\n\n"
		}
		longDesc += "If the resource exists, it will be updated. If it doesn't exist, it will be created."

		exampleText := fmt.Sprintf(`  # Apply from file
  f5xcctl %s apply %s -n my-namespace -i config.yaml`, domain, rt.Name)

		subCmd := &cobra.Command{
			Use:     rt.Name,
			Short:   formatShortWithTier("Apply", displayName, rt.Name),
			Long:    longDesc,
			Example: exampleText,
			RunE: func(cmd *cobra.Command, args []string) error {
				return runConfigApply(rtCopy, &flags)
			},
		}
		_ = subCmd.RegisterFlagCompletionFunc("namespace", completeNamespace)
		cmd.AddCommand(subCmd)
	}

	return cmd
}

// buildDomainPatchCmd creates the patch operation for a domain
func buildDomainPatchCmd(domain string) *cobra.Command {
	var flags configurationFlags

	domainInfo, _ := types.GetDomainInfo(domain)

	cmd := &cobra.Command{
		Use:   "patch",
		Short: fmt.Sprintf("Partially update a %s resource", domainInfo.DisplayName),
		Long: fmt.Sprintf(`Partially update a %s resource in F5 Distributed Cloud.

%s

Only specified fields will be updated. Other fields remain unchanged.`, domainInfo.DisplayName, domainInfo.Description),
	}

	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Target namespace for the operation.")
	cmd.PersistentFlags().StringVarP(&flags.inputFile, "input-file", "i", "", "Path to YAML/JSON file containing the fields to patch.")

	// Register flag completions
	_ = cmd.RegisterFlagCompletionFunc("namespace", completeNamespace)

	// Get resources for this domain
	resources := types.GetByDomain(domain)
	sortResourcesByName(resources)

	// Add resource type subcommands for this domain
	for _, rt := range resources {
		if !rt.Operations.Update {
			continue // Skip resources that don't support update
		}

		rtCopy := rt

		displayName := naming.ToHumanReadable(rt.Name)
		longDesc := fmt.Sprintf("Partially update a %s resource.\n\n", displayName)
		if rt.Description != "" {
			longDesc += rt.Description + "\n\n"
		}
		longDesc += "Only specified fields will be updated."

		exampleText := fmt.Sprintf(`  # Patch from file
  f5xcctl %s patch %s -n my-namespace my-resource -i patch.yaml`, domain, rt.Name)

		subCmd := &cobra.Command{
			Use:               fmt.Sprintf("%s <name>", rt.Name),
			Short:             formatShortWithTier("Patch", displayName, rt.Name),
			Long:              longDesc,
			Example:           exampleText,
			Args:              cobra.ExactArgs(1),
			ValidArgsFunction: completeResourceName(domain, rtCopy.Name),
			RunE: func(cmd *cobra.Command, args []string) error {
				flags.name = args[0]
				return runConfigPatch(rtCopy, &flags)
			},
		}
		_ = subCmd.RegisterFlagCompletionFunc("namespace", completeNamespace)
		cmd.AddCommand(subCmd)
	}

	return cmd
}

// buildDomainAddLabelsCmd creates the add-labels operation for a domain
func buildDomainAddLabelsCmd(domain string) *cobra.Command {
	var flags configurationFlags

	domainInfo, _ := types.GetDomainInfo(domain)

	cmd := &cobra.Command{
		Use:   "add-labels",
		Short: fmt.Sprintf("Add labels to a %s resource", domainInfo.DisplayName),
		Long: fmt.Sprintf(`Add labels to a %s resource in F5 Distributed Cloud.

%s

Specify label key-value pairs using --label-key and --label-value flags.`, domainInfo.DisplayName, domainInfo.Description),
	}

	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Target namespace for the operation.")
	cmd.PersistentFlags().StringArrayVar(&flags.labelKeys, "label-key", []string{}, "Label keys to add.")
	cmd.PersistentFlags().StringArrayVar(&flags.labelValues, "label-value", []string{}, "Corresponding label values.")

	// Register flag completions
	_ = cmd.RegisterFlagCompletionFunc("namespace", completeNamespace)
	_ = cmd.RegisterFlagCompletionFunc("label-key", completeLabelKey)

	// Get resources for this domain
	resources := types.GetByDomain(domain)
	sortResourcesByName(resources)

	// Add resource type subcommands for this domain
	for _, rt := range resources {
		rtCopy := rt

		displayName := naming.ToHumanReadable(rt.Name)
		longDesc := fmt.Sprintf("Add labels to a %s resource.\n\n", displayName)
		if rt.Description != "" {
			longDesc += rt.Description + "\n\n"
		}

		exampleText := fmt.Sprintf(`  # Add labels
  f5xcctl %s add-labels %s -n my-namespace my-resource --label-key env --label-value prod`, domain, rt.Name)

		subCmd := &cobra.Command{
			Use:               fmt.Sprintf("%s <name>", rt.Name),
			Short:             formatShortWithTier("Add-labels", displayName, rt.Name),
			Long:              longDesc,
			Example:           exampleText,
			Args:              cobra.ExactArgs(1),
			ValidArgsFunction: completeResourceName(domain, rtCopy.Name),
			RunE: func(cmd *cobra.Command, args []string) error {
				flags.name = args[0]
				return runConfigAddLabels(rtCopy, &flags)
			},
		}
		_ = subCmd.RegisterFlagCompletionFunc("namespace", completeNamespace)
		_ = subCmd.RegisterFlagCompletionFunc("label-key", completeLabelKey)
		cmd.AddCommand(subCmd)
	}

	return cmd
}

// buildDomainRemoveLabelsCmd creates the remove-labels operation for a domain
func buildDomainRemoveLabelsCmd(domain string) *cobra.Command {
	var flags configurationFlags

	domainInfo, _ := types.GetDomainInfo(domain)

	cmd := &cobra.Command{
		Use:   "remove-labels",
		Short: fmt.Sprintf("Remove labels from a %s resource", domainInfo.DisplayName),
		Long: fmt.Sprintf(`Remove labels from a %s resource in F5 Distributed Cloud.

%s

Specify the label keys to remove using --label-key flags.`, domainInfo.DisplayName, domainInfo.Description),
	}

	cmd.PersistentFlags().StringVarP(&flags.namespace, "namespace", "n", "default", "Target namespace for the operation.")
	cmd.PersistentFlags().StringArrayVar(&flags.labelKeys, "label-key", []string{}, "Label keys to remove.")

	// Register flag completions
	_ = cmd.RegisterFlagCompletionFunc("namespace", completeNamespace)
	_ = cmd.RegisterFlagCompletionFunc("label-key", completeLabelKey)

	// Get resources for this domain
	resources := types.GetByDomain(domain)
	sortResourcesByName(resources)

	// Add resource type subcommands for this domain
	for _, rt := range resources {
		rtCopy := rt

		displayName := naming.ToHumanReadable(rt.Name)
		longDesc := fmt.Sprintf("Remove labels from a %s resource.\n\n", displayName)
		if rt.Description != "" {
			longDesc += rt.Description + "\n\n"
		}

		exampleText := fmt.Sprintf(`  # Remove labels
  f5xcctl %s remove-labels %s -n my-namespace my-resource --label-key env`, domain, rt.Name)

		subCmd := &cobra.Command{
			Use:               fmt.Sprintf("%s <name>", rt.Name),
			Short:             formatShortWithTier("Remove-labels", displayName, rt.Name),
			Long:              longDesc,
			Example:           exampleText,
			Args:              cobra.ExactArgs(1),
			ValidArgsFunction: completeResourceName(domain, rtCopy.Name),
			RunE: func(cmd *cobra.Command, args []string) error {
				flags.name = args[0]
				return runConfigRemoveLabels(rtCopy, &flags)
			},
		}
		_ = subCmd.RegisterFlagCompletionFunc("namespace", completeNamespace)
		_ = subCmd.RegisterFlagCompletionFunc("label-key", completeLabelKey)
		cmd.AddCommand(subCmd)
	}

	return cmd
}

// sortResourcesByName sorts resources by name for consistent output
func sortResourcesByName(resources []*types.ResourceType) {
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].Name < resources[j].Name
	})
}
