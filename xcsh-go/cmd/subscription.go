package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/robinmordasiewicz/xcsh/pkg/subscription"
)

// subscriptionClient is the shared subscription client for all subscription subcommands
var subscriptionClient *subscription.Client

// subscriptionNamespace is the namespace flag for subscription commands
var subscriptionNamespace string

var subscriptionCmd = &cobra.Command{
	Use:     "subscription",
	Aliases: []string{"sub", "subs"},
	Short:   "Manage and inspect F5 XC subscription, addons, and quotas.",
	Long: `Subscription management commands for F5 Distributed Cloud.

View your current subscription tier (Standard/Advanced), active and available
addon services, tenant-level quota limits and usage, and validate Terraform
plans against subscription capabilities.

AI assistants should use 'xcsh subscription show --output-format json' to
understand tenant capabilities before attempting resource deployments.

SUBSCRIPTION TIERS:
  Standard    Base tier with core F5 XC functionality
  Advanced    Enhanced tier with additional security and management features

ADDON SERVICE STATES:
  AS_SUBSCRIBED     Service is actively subscribed and available
  AS_PENDING        Service subscription is being processed
  AS_NONE           Service is not subscribed
  AS_ERROR          Service subscription has an error

ACCESS STATUS:
  ALLOWED           Can subscribe to or use this service
  UPGRADE_REQUIRED  Requires a plan upgrade to access
  CONTACT_SALES     Requires contacting F5 sales
  DENIED            Access is denied by policy

QUOTA ENFORCEMENT:
  Quotas are enforced at the TENANT level, not per-namespace. Resource counts
  across all namespaces accumulate toward the same tenant-wide quota limits.`,
	Example: `  # Show subscription summary
  xcsh subscription show

  # Show subscription as JSON for automation
  xcsh subscription show --output-format json

  # List active addon services
  xcsh subscription addons --filter active

  # List all addon services including denied ones
  xcsh subscription addons --all

  # Check tenant quota usage
  xcsh subscription quota

  # Check quota usage as JSON
  xcsh subscription quota --output-format json

  # Validate if you can create 5 more HTTP load balancers
  xcsh subscription validate --resource-type http_loadbalancer --count 5

  # Validate if bot-defense feature is available
  xcsh subscription validate --feature bot-defense

  # Get subscription spec for AI assistants
  xcsh subscription --spec --output-format json`,
	// Note: We don't override PersistentPreRunE here to let root's PersistentPreRunE run
	// and initialize the API client. The subscription client is lazily initialized in GetSubscriptionClient()
}

func init() {
	rootCmd.AddCommand(subscriptionCmd)

	// Enable AI-agent-friendly error handling for invalid subcommands
	subscriptionCmd.RunE = func(cmd *cobra.Command, args []string) error {
		// Handle --spec flag
		if CheckSpecFlag() {
			spec := subscription.GenerateSpec()
			format := GetOutputFormatWithDefault("json")
			return OutputSubscriptionSpec(spec, format)
		}

		if len(args) > 0 {
			return fmt.Errorf("unknown command %q for %q\n\nUsage: xcsh subscription <command> [flags]\n\nAvailable Commands:\n  show, addons, quota, validate\n\nRun 'xcsh subscription --help' for usage", args[0], cmd.CommandPath())
		}
		return cmd.Help()
	}
	subscriptionCmd.SuggestionsMinimumDistance = 2

	// Register --spec flag for machine-readable CLI specification
	RegisterSpecFlag(subscriptionCmd)

	// Add persistent flags for subscription commands
	subscriptionCmd.PersistentFlags().StringVarP(&subscriptionNamespace, "namespace", "n", "", "Namespace to check (default: system)")
}

// GetSubscriptionClient returns the subscription client, initializing it if needed
func GetSubscriptionClient() *subscription.Client {
	if subscriptionClient == nil {
		apiClient := GetClient()
		if apiClient != nil {
			subscriptionClient = subscription.NewClient(apiClient)
		}
	}
	return subscriptionClient
}

// GetSubscriptionNamespace returns the namespace for subscription commands
func GetSubscriptionNamespace() string {
	if subscriptionNamespace == "" {
		return "system"
	}
	return subscriptionNamespace
}

// OutputSubscriptionSpec outputs the subscription specification in the requested format
func OutputSubscriptionSpec(spec *subscription.Spec, format string) error {
	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(spec); err != nil {
			return fmt.Errorf("failed to encode JSON: %w", err)
		}
		os.Exit(0)
		return nil
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		if err := encoder.Encode(spec); err != nil {
			return fmt.Errorf("failed to encode YAML: %w", err)
		}
		os.Exit(0)
		return nil
	default:
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(spec); err != nil {
			return fmt.Errorf("failed to encode JSON: %w", err)
		}
		os.Exit(0)
		return nil
	}
}
