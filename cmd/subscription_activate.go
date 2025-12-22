package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	xerrors "github.com/robinmordasiewicz/xcsh/pkg/errors"
	"github.com/robinmordasiewicz/xcsh/pkg/subscription"
)

var (
	activateAddonName string
)

var subscriptionActivateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Activate an addon service subscription.",
	Long: `Activate an addon service for the current tenant.

Activation behavior depends on the addon's activation type:
  - Self-Activation: Activates immediately and the addon becomes available
  - Partially Managed: Creates a pending request with partial backend processing
  - Fully Managed: Creates a pending request requiring SRE approval

Use 'f5xcctl subscription addons --filter available' to see which addons can be activated.
Use 'f5xcctl subscription activation-status' to check pending activation requests.

EXIT CODES:
  0 - Activation successful or pending request created
  1 - Generic error (API failure, invalid arguments)
  9 - Feature not available (access denied, requires upgrade, contact sales)`,
	Example: `  # Activate bot-defense addon
  f5xcctl subscription activate --addon bot-defense

  # Activate addon in specific namespace
  f5xcctl subscription activate --addon api-security -n production

  # Get activation result as JSON
  f5xcctl subscription activate --addon client-side-defense --output-format json`,
	RunE: runSubscriptionActivate,
}

func init() {
	subscriptionCmd.AddCommand(subscriptionActivateCmd)

	subscriptionActivateCmd.Flags().StringVar(&activateAddonName, "addon", "",
		"Name of the addon service to activate (required).")

	_ = subscriptionActivateCmd.MarkFlagRequired("addon")

	// Register completion for --addon flag (shows available addons)
	_ = subscriptionActivateCmd.RegisterFlagCompletionFunc("addon", completeAddonName)
}

func runSubscriptionActivate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	client, err := requireSubscriptionClient()
	if err != nil {
		return err
	}

	namespace := GetSubscriptionNamespace()

	// Attempt activation
	result, err := client.ActivateAddon(ctx, namespace, activateAddonName)

	// Handle activation denied errors with appropriate exit code
	var deniedErr *subscription.ActivationDeniedError
	if errors.As(err, &deniedErr) {
		format := GetOutputFormatWithDefault("table")
		if outputErr := formatOutputWithTableFallback(result, format, func() error {
			return outputActivationTable(result)
		}); outputErr != nil {
			return outputErr
		}
		os.Exit(xerrors.ExitFeatureNotAvail) // Exit code 9
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to activate addon: %w", err)
	}

	// Output result
	format := GetOutputFormatWithDefault("table")
	return formatOutputWithTableFallback(result, format, func() error {
		return outputActivationTable(result)
	})
}

func outputActivationTable(result *subscription.ActivationResponse) error {
	if result == nil {
		return fmt.Errorf("no activation result")
	}

	// Determine status indicator
	var statusLabel string
	if result.IsImmediate {
		statusLabel = "SUCCESS"
	} else if result.IsPending {
		statusLabel = "PENDING"
	} else if result.AccessStatus != subscription.AccessAllowed && result.AccessStatus != "" {
		statusLabel = "DENIED"
	} else {
		statusLabel = "INFO"
	}

	fmt.Printf("ACTIVATION RESULT: %s\n", statusLabel)
	fmt.Println(strings.Repeat("=", 75))
	fmt.Println()

	// Details
	fmt.Printf("  %-20s %s\n", "Addon Service:", result.AddonService)
	if result.Namespace != "" {
		fmt.Printf("  %-20s %s\n", "Namespace:", result.Namespace)
	}
	if result.ActivationType != "" {
		fmt.Printf("  %-20s %s\n", "Activation Type:",
			subscription.ActivationTypeDescription(result.ActivationType))
	}
	if result.SubscriptionState != "" {
		fmt.Printf("  %-20s %s\n", "State:",
			subscription.SubscriptionStateDescription(result.SubscriptionState))
	}
	if result.RequestID != "" {
		fmt.Printf("  %-20s %s\n", "Request ID:", result.RequestID)
	}
	fmt.Println()

	// Message with status icon
	var statusIcon string
	switch statusLabel {
	case "SUCCESS":
		statusIcon = "[OK]"
	case "PENDING":
		statusIcon = "[PENDING]"
	case "DENIED":
		statusIcon = "[DENIED]"
	default:
		statusIcon = "[INFO]"
	}
	fmt.Printf("  %s %s\n", statusIcon, result.Message)
	fmt.Println()

	// Next steps if applicable
	if result.NextSteps != "" {
		fmt.Println("NEXT STEPS")
		fmt.Println(strings.Repeat("-", 75))
		fmt.Printf("  %s\n", result.NextSteps)
		fmt.Println()
	}

	// Guidance for denied access
	if result.AccessStatus != subscription.AccessAllowed && result.AccessStatus != "" {
		fmt.Println("GUIDANCE")
		fmt.Println(strings.Repeat("-", 75))
		switch result.AccessStatus {
		case subscription.AccessUpgradeRequired:
			fmt.Println("  This addon requires a subscription plan upgrade.")
			fmt.Println("  Options:")
			fmt.Println("    1. Upgrade via F5 XC Console: Settings > Subscription > Upgrade Plan")
			fmt.Println("    2. Contact your F5 account manager for plan options")
		case subscription.AccessContactSales:
			fmt.Println("  This addon requires contacting F5 sales for enablement.")
			fmt.Println("  Options:")
			fmt.Println("    1. Contact F5 sales: https://www.f5.com/company/contact/sales")
			fmt.Println("    2. Reach out to your F5 account manager")
		case subscription.AccessDenied:
			fmt.Println("  Access to this addon is denied by policy.")
			fmt.Println("  Check your tenant policies or contact your administrator.")
		case subscription.AccessEOL:
			fmt.Println("  This addon is end-of-life and cannot be activated.")
			fmt.Println("  Consider alternative addons or contact F5 support for migration guidance.")
		case subscription.AccessInternalService:
			fmt.Println("  This addon is an internal service managed by F5.")
			fmt.Println("  It cannot be activated through the API.")
		}
		fmt.Println()
	}

	return nil
}
