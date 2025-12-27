package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/robinmordasiewicz/xcsh/pkg/subscription"
)

var (
	activationStatusAddon string // Optional: filter to specific addon
)

var subscriptionActivationStatusCmd = &cobra.Command{
	Use:     "activation-status",
	Aliases: []string{"status", "pending"},
	Short:   "Check pending addon activation requests.",
	Long: `Display the status of pending addon activation requests.

Shows all addon services with pending activation requests, including:
  - Addon service name
  - Activation type (self, partially managed, fully managed)
  - Current subscription state
  - Expected processing time based on activation type

Use this command after 'xcsh subscription activate' to monitor
the progress of managed activation requests.`,
	Example: `  # Show all pending activation requests
  xcsh subscription activation-status

  # Check status for a specific addon
  xcsh subscription activation-status --addon bot-defense

  # Get status as JSON for automation
  xcsh subscription activation-status --output-format json

  # Check status in a specific namespace
  xcsh subscription activation-status -n production`,
	RunE: runSubscriptionActivationStatus,
}

func init() {
	subscriptionCmd.AddCommand(subscriptionActivationStatusCmd)

	subscriptionActivationStatusCmd.Flags().StringVar(&activationStatusAddon, "addon", "",
		"Filter to a specific addon service.")

	// Register completion for --addon flag (shows pending and active addons)
	_ = subscriptionActivationStatusCmd.RegisterFlagCompletionFunc("addon", completePendingAddonName)
}

func runSubscriptionActivationStatus(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	client, err := requireSubscriptionClient()
	if err != nil {
		return err
	}

	namespace := GetSubscriptionNamespace()
	result, err := client.GetPendingActivations(ctx, namespace)
	if err != nil {
		return fmt.Errorf("failed to get activation status: %w", err)
	}

	// Filter to specific addon if requested
	if activationStatusAddon != "" {
		var filtered []subscription.PendingActivation
		for _, p := range result.PendingActivations {
			if strings.EqualFold(p.AddonService, activationStatusAddon) {
				filtered = append(filtered, p)
			}
		}
		result.PendingActivations = filtered
		result.TotalPending = len(filtered)
	}

	format := GetOutputFormatWithDefault("table")
	return formatOutputWithTableFallback(result, format, func() error {
		return outputActivationStatusTable(result)
	})
}

func outputActivationStatusTable(result *subscription.ActivationStatusResult) error {
	fmt.Println("ADDON ACTIVATION STATUS")
	fmt.Println(strings.Repeat("=", 85))
	fmt.Println()

	// Summary
	fmt.Printf("  Pending Activations: %d\n", result.TotalPending)
	fmt.Printf("  Active Addons: %d\n", len(result.ActiveAddons))
	fmt.Println()

	if result.TotalPending == 0 {
		fmt.Println("  No pending activation requests.")
		fmt.Println()

		if len(result.ActiveAddons) > 0 {
			fmt.Println("ACTIVE ADDONS")
			fmt.Println(strings.Repeat("-", 85))
			for _, addon := range result.ActiveAddons {
				fmt.Printf("  [ACTIVE] %s\n", addon)
			}
			fmt.Println()
		}

		fmt.Println("HINTS")
		fmt.Println(strings.Repeat("-", 85))
		fmt.Println("  Use 'xcsh subscription addons --filter available' to see activatable addons")
		fmt.Println("  Use 'xcsh subscription activate --addon <name>' to activate an addon")
		fmt.Println()
		return nil
	}

	// Pending activations table
	fmt.Println("PENDING ACTIVATIONS")
	fmt.Println(strings.Repeat("-", 85))
	fmt.Printf("  %-25s %-25s %-18s %-15s\n", "ADDON", "TYPE", "STATE", "MESSAGE")
	fmt.Println("  " + strings.Repeat("-", 83))

	for _, p := range result.PendingActivations {
		name := p.AddonService
		if len(name) > 24 {
			name = name[:21] + "..."
		}

		typeDesc := subscription.ActivationTypeDescription(p.ActivationType)
		if len(typeDesc) > 24 {
			typeDesc = typeDesc[:21] + "..."
		}

		stateDesc := subscription.SubscriptionStateDescription(p.SubscriptionState)
		if len(stateDesc) > 17 {
			stateDesc = stateDesc[:14] + "..."
		}

		message := p.Message
		if len(message) > 14 {
			message = message[:11] + "..."
		}

		fmt.Printf("  %-25s %-25s %-18s %-15s\n", name, typeDesc, stateDesc, message)
	}
	fmt.Println()

	// Expectations based on activation type
	fmt.Println("EXPECTED PROCESSING TIMES")
	fmt.Println(strings.Repeat("-", 85))
	fmt.Println("  Self-Activation:           Usually immediate, up to a few minutes")
	fmt.Println("  Partially Managed:         Minutes to hours, depends on backend processing")
	fmt.Println("  Fully Managed:             Up to 24 hours, requires SRE approval")
	fmt.Println()

	fmt.Println("HINTS")
	fmt.Println(strings.Repeat("-", 85))
	fmt.Println("  Run this command again to check for status updates")
	fmt.Println("  Use --output-format json for automation scripts")
	fmt.Println()

	return nil
}
