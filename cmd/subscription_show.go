package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/robinmordasiewicz/f5xcctl/pkg/subscription"
)

var subscriptionShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current subscription tier and summary.",
	Long: `Display the current subscription tier and a summary of active services.

Shows the subscription tier (Standard/Advanced), plan information, active addon
services, and a quota usage summary. Use --output-format json for machine-readable
output suitable for AI assistants and automation.

AI assistants should call this command first to understand tenant capabilities
before attempting to deploy resources that may require specific subscription tiers
or addon services.`,
	Example: `  # Show subscription summary in table format
  f5xcctl subscription show

  # Show subscription as JSON for automation
  f5xcctl subscription show --output-format json

  # Show subscription as YAML
  f5xcctl subscription show --output-format yaml`,
	RunE: runSubscriptionShow,
}

func init() {
	subscriptionCmd.AddCommand(subscriptionShowCmd)
}

func runSubscriptionShow(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	client, err := requireSubscriptionClient()
	if err != nil {
		return err
	}

	info, err := client.GetSubscriptionInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get subscription info: %w", err)
	}

	// Output based on format
	format := GetOutputFormatWithDefault("table")
	return formatOutputWithTableFallback(info, format, func() error {
		return outputSubscriptionTable(info)
	})
}

func outputSubscriptionTable(info *subscription.SubscriptionInfo) error {
	// Print header
	fmt.Println("SUBSCRIPTION SUMMARY")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	// Tier and Plan info
	fmt.Printf("  %-20s %s\n", "Tier:", info.Tier)
	if info.Plan.Name != "" {
		fmt.Printf("  %-20s %s\n", "Plan:", info.Plan.DisplayName)
	}
	if info.TenantName != "" {
		fmt.Printf("  %-20s %s\n", "Tenant:", info.TenantName)
	}
	fmt.Println()

	// Active Addons Summary
	totalAddons := len(info.ActiveAddons) + len(info.AvailableAddons)
	fmt.Printf("ACTIVE ADDONS (%d/%d available)\n", len(info.ActiveAddons), totalAddons)
	fmt.Println(strings.Repeat("-", 60))

	if len(info.ActiveAddons) == 0 {
		fmt.Println("  No active addon services")
	} else {
		// Print header
		fmt.Printf("  %-28s %-12s %-15s\n", "NAME", "TIER", "STATUS")
		fmt.Println("  " + strings.Repeat("-", 55))

		// Print active addons (max 5 for summary)
		displayCount := len(info.ActiveAddons)
		if displayCount > 5 {
			displayCount = 5
		}
		for i := 0; i < displayCount; i++ {
			addon := info.ActiveAddons[i]
			name := addon.Name
			if addon.DisplayName != "" {
				name = addon.DisplayName
			}
			if len(name) > 27 {
				name = name[:24] + "..."
			}
			fmt.Printf("  %-28s %-12s %-15s\n",
				name,
				subscription.TierDescription(addon.Tier),
				subscription.StateDescription(addon.State))
		}
		if len(info.ActiveAddons) > 5 {
			fmt.Printf("  ... and %d more active addons\n", len(info.ActiveAddons)-5)
		}
	}
	fmt.Println()

	// Quota Summary
	fmt.Println("QUOTA SUMMARY")
	fmt.Println(strings.Repeat("-", 60))

	summary := info.QuotaSummary
	fmt.Printf("  %-20s %d\n", "Total Limits:", summary.TotalLimits)

	if summary.LimitsExceeded > 0 {
		fmt.Printf("  %-20s %d (CRITICAL)\n", "Limits Exceeded:", summary.LimitsExceeded)
	}
	if summary.LimitsAtRisk > 0 {
		fmt.Printf("  %-20s %d (WARNING)\n", "Limits at Risk:", summary.LimitsAtRisk)
	}
	fmt.Println()

	// Show top quotas approaching limits
	if len(summary.Objects) > 0 {
		fmt.Printf("  %-28s %-10s %-10s %-10s\n", "RESOURCE", "USAGE", "LIMIT", "%USED")
		fmt.Println("  " + strings.Repeat("-", 58))

		// Sort by percentage and show top 5 or those at risk
		displayQuotas := selectTopQuotas(summary.Objects, 5)
		for _, q := range displayQuotas {
			statusIndicator := ""
			switch q.Status {
			case "EXCEEDED":
				statusIndicator = " (!!!)"
			case "WARNING":
				statusIndicator = " (!)"
			}

			name := q.DisplayName
			if name == "" {
				name = q.Name
			}
			if len(name) > 27 {
				name = name[:24] + "..."
			}

			fmt.Printf("  %-28s %-10.0f %-10.0f %5.0f%%%s\n",
				name, q.Usage, q.Limit, q.Percentage, statusIndicator)
		}
	}
	fmt.Println()

	// Hints for AI assistants
	fmt.Println("HINTS")
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println("  Use 'f5xcctl subscription addons' for detailed addon list")
	fmt.Println("  Use 'f5xcctl subscription quota' for full quota details")
	fmt.Println("  Use 'f5xcctl subscription validate' before terraform apply")
	fmt.Println()

	return nil
}

// selectTopQuotas returns the top N quotas by usage percentage, prioritizing those at risk
func selectTopQuotas(quotas []subscription.QuotaItem, n int) []subscription.QuotaItem {
	if len(quotas) <= n {
		return quotas
	}

	// Separate by status
	var exceeded, atRisk, ok []subscription.QuotaItem
	for _, q := range quotas {
		switch q.Status {
		case "EXCEEDED":
			exceeded = append(exceeded, q)
		case "WARNING":
			atRisk = append(atRisk, q)
		default:
			ok = append(ok, q)
		}
	}

	// Combine: exceeded first, then at risk, then OK
	var result []subscription.QuotaItem
	result = append(result, exceeded...)
	result = append(result, atRisk...)
	result = append(result, ok...)

	if len(result) > n {
		return result[:n]
	}
	return result
}
