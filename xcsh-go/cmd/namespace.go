package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/robinmordasiewicz/xcsh/pkg/output"
)

var namespaceCmd = &cobra.Command{
	Use:   "namespace [name]",
	Short: "Show or set the default namespace",
	Long: `Show or set the default namespace for CRUD operations.

If no namespace is provided, displays the current default.
If a namespace is provided, sets it as the new default.

Priority order:
  1. F5XC_DEFAULT_NAMESPACE environment variable
  2. Config file default-namespace setting
  3. Fallback to "default"`,
	Args: cobra.MaximumNArgs(1),
	Example: `  # Show current default namespace
  xcsh namespace

  # Set default namespace to 'shared'
  xcsh namespace shared

  # Use environment variable (highest priority)
  export F5XC_DEFAULT_NAMESPACE=production
  xcsh namespace`,
	RunE:              runNamespace,
	ValidArgsFunction: completeNamespaceArg,
}

func init() {
	rootCmd.AddCommand(namespaceCmd)
}

func runNamespace(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		// Show current default namespace
		ns := GetDefaultNamespace()
		source := getNamespaceSource()
		fmt.Printf("Default namespace: %s\n", ns)
		fmt.Printf("Source: %s\n", source)

		// Validate and warn if namespace doesn't exist
		if ns != "default" && !NamespaceExists(ns) {
			output.PrintWarning(fmt.Sprintf("Namespace '%s' does not exist in tenant", ns))
		}
		return nil
	}

	// Set new default namespace
	newNamespace := args[0]

	// Validate namespace exists before setting
	if !NamespaceExists(newNamespace) {
		return fmt.Errorf("namespace '%s' does not exist", newNamespace)
	}

	if err := setDefaultNamespace(newNamespace); err != nil {
		return err
	}

	// Clear cache so next GetValidatedDefaultNamespace() picks up the change
	ClearNamespaceCache()

	output.PrintInfo(fmt.Sprintf("Default namespace set to: %s", newNamespace))
	return nil
}

// getNamespaceSource returns the source of the current default namespace
func getNamespaceSource() string {
	if ns := os.Getenv("F5XC_DEFAULT_NAMESPACE"); ns != "" {
		return "environment variable (F5XC_DEFAULT_NAMESPACE)"
	}

	configPath := getConfigPath()
	if data, err := os.ReadFile(configPath); err == nil {
		if config, err := parseConfigFile(data); err == nil {
			if config.DefaultNamespace != "" {
				return fmt.Sprintf("config file (%s)", configPath)
			}
		}
	}

	return "default fallback"
}

// setDefaultNamespace saves the namespace to the config file
func setDefaultNamespace(namespace string) error {
	configPath := getConfigPath()

	// Load existing config or create new one
	config := &ConfigFile{}
	if data, err := os.ReadFile(configPath); err == nil {
		if parsed, err := parseConfigFile(data); err == nil {
			config = parsed
		}
	}

	// Update namespace
	config.DefaultNamespace = namespace

	// Save configuration
	return saveConfig(config, configPath)
}

// validatedNamespaceCache stores the result of namespace validation
// to avoid repeated API calls and duplicate warnings
var validatedNamespaceCache string

// NamespaceExists checks if a namespace exists in the tenant
func NamespaceExists(namespace string) bool {
	client := GetClient()
	if client == nil {
		// Can't validate without a client, assume it exists
		return true
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	path := fmt.Sprintf("/api/web/namespaces/%s", namespace)
	resp, err := client.Get(ctx, path, nil)
	if err != nil {
		// Network error, assume namespace exists
		return true
	}

	return resp.StatusCode == 200
}

// GetValidatedDefaultNamespace returns the default namespace, validating it exists
// If the configured namespace doesn't exist, falls back to "default" and warns
// Results are cached to avoid repeated API calls and duplicate warnings
func GetValidatedDefaultNamespace() string {
	// Return cached result if available
	if validatedNamespaceCache != "" {
		return validatedNamespaceCache
	}

	ns := GetDefaultNamespace()

	// Skip validation if we don't have a client
	client := GetClient()
	if client == nil {
		return ns
	}

	// "default" namespace always exists, no need to validate
	if ns == "default" {
		validatedNamespaceCache = ns
		return ns
	}

	// Validate the namespace exists
	if !NamespaceExists(ns) {
		source := getNamespaceSource()
		output.PrintWarning(fmt.Sprintf("Namespace '%s' (from %s) does not exist, using 'default'", ns, source))
		validatedNamespaceCache = "default"
		return "default"
	}

	validatedNamespaceCache = ns
	return ns
}

// ClearNamespaceCache clears the validated namespace cache
// Used when namespace is changed via the namespace command
func ClearNamespaceCache() {
	validatedNamespaceCache = ""
}

// completeNamespaceArg provides shell completion for the namespace command argument
func completeNamespaceArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Only complete first argument (namespace name)
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	// Reuse existing namespace completion logic
	return completeNamespace(cmd, args, toComplete)
}
