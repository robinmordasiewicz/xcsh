package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/robinmordasiewicz/xcsh/pkg/subscription"
)

var (
	addonsShowAll bool
	addonsFilter  string
)

var subscriptionAddonsCmd = &cobra.Command{
	Use:   "addons",
	Short: "List all addon services with activation status.",
	Long: `List all addon services available for the subscription.

Shows each addon service with its tier (BASIC, STANDARD, ADVANCED, PREMIUM),
subscription state (SUBSCRIBED, PENDING, NONE, ERROR), and access status
(ALLOWED, UPGRADE_REQUIRED, CONTACT_SALES, DENIED).

Use --filter to show only specific addon types:
  - active:    Show only actively subscribed addons
  - available: Show addons that can be subscribed to
  - denied:    Show addons that require upgrade or sales contact

AI assistants should check addon status before deploying resources that depend
on specific services like bot-defense, api-security, or web-app-firewall.`,
	Example: `  # List all addon services (excluding denied by default)
  f5xcctl subscription addons

  # List all addon services including denied ones
  f5xcctl subscription addons --all

  # List only active addons
  f5xcctl subscription addons --filter active

  # List available addons that can be activated
  f5xcctl subscription addons --filter available

  # List addons that need upgrade or sales contact
  f5xcctl subscription addons --filter denied

  # Output as JSON for automation
  f5xcctl subscription addons --output-format json`,
	RunE: runSubscriptionAddons,
}

func init() {
	subscriptionCmd.AddCommand(subscriptionAddonsCmd)

	subscriptionAddonsCmd.Flags().BoolVar(&addonsShowAll, "all", false, "Show all addon services including denied ones.")
	subscriptionAddonsCmd.Flags().StringVar(&addonsFilter, "filter", "", "Filter by status: active, available, denied.")

	// Register completion for --filter flag
	_ = subscriptionAddonsCmd.RegisterFlagCompletionFunc("filter", completeAddonFilter)
}

func runSubscriptionAddons(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	client, err := requireSubscriptionClient()
	if err != nil {
		return err
	}

	namespace := GetSubscriptionNamespace()
	addons, err := client.GetAddonServices(ctx, namespace)
	if err != nil {
		return fmt.Errorf("failed to get addon services: %w", err)
	}

	// Apply filter
	if addonsFilter != "" {
		addons = subscription.FilterAddons(addons, addonsFilter)
	} else if !addonsShowAll {
		// By default, exclude denied addons unless --all is specified
		var filtered []subscription.AddonServiceInfo
		for _, addon := range addons {
			if !addon.IsDenied() {
				filtered = append(filtered, addon)
			}
		}
		addons = filtered
	}

	// Output based on format
	format := GetOutputFormatWithDefault("table")
	return formatOutputWithTableFallback(addons, format, func() error {
		return outputAddonsTable(addons)
	})
}

func outputAddonsTable(addons []subscription.AddonServiceInfo) error {
	if len(addons) == 0 {
		fmt.Println("No addon services found matching the criteria.")
		fmt.Println("\nUse --all to show all addon services, or --filter to change criteria.")
		return nil
	}

	// Count by status
	var activeCount, availableCount, deniedCount int
	for _, addon := range addons {
		if addon.IsActive() {
			activeCount++
		} else if addon.IsAvailable() {
			availableCount++
		} else if addon.IsDenied() {
			deniedCount++
		}
	}

	// Print header
	fmt.Println("ADDON SERVICES")
	fmt.Println(strings.Repeat("=", 90))
	fmt.Printf("  Active: %d | Available: %d | Restricted: %d | Total: %d\n",
		activeCount, availableCount, deniedCount, len(addons))
	fmt.Println()

	// Print table header
	fmt.Printf("  %-25s %-25s %-12s %-12s %-12s\n",
		"NAME", "DISPLAY NAME", "TIER", "STATUS", "ACCESS")
	fmt.Println("  " + strings.Repeat("-", 86))

	// Print addons
	for _, addon := range addons {
		name := addon.Name
		if len(name) > 24 {
			name = name[:21] + "..."
		}

		displayName := addon.DisplayName
		if displayName == "" {
			displayName = "-"
		}
		if len(displayName) > 24 {
			displayName = displayName[:21] + "..."
		}

		tierDesc := subscription.TierDescription(addon.Tier)
		stateDesc := subscription.StateDescription(addon.State)
		accessDesc := getShortAccessStatus(addon.AccessStatus)

		fmt.Printf("  %-25s %-25s %-12s %-12s %-12s\n",
			name, displayName, tierDesc, stateDesc, accessDesc)
	}
	fmt.Println()

	// Show hints based on what's displayed
	if deniedCount > 0 {
		fmt.Println("NOTES")
		fmt.Println(strings.Repeat("-", 90))
		fmt.Println("  Addons with 'Upgrade Req' need a plan upgrade to access")
		fmt.Println("  Addons with 'Sales' require contacting F5 sales")
		fmt.Println("  Use 'f5xcctl subscription validate --feature <name>' to check specific addon availability")
		fmt.Println()
	}

	return nil
}

func getShortAccessStatus(status string) string {
	switch status {
	case subscription.AccessAllowed:
		return "Allowed"
	case subscription.AccessDenied:
		return "Denied"
	case subscription.AccessUpgradeRequired:
		return "Upgrade Req"
	case subscription.AccessContactSales:
		return "Sales"
	case subscription.AccessInternalService:
		return "Internal"
	default:
		return "Unknown"
	}
}
