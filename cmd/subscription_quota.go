package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/robinmordasiewicz/f5xcctl/pkg/subscription"
)

var (
	quotaType string
)

var subscriptionQuotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "Display tenant-level quota limits and usage.",
	Long: `Display tenant-level quota limits and current usage.

Quotas are enforced at the TENANT level, not per-namespace. All resources across
all namespaces count toward the same tenant-wide quota limits. This means:
- Quota limits apply globally to your entire F5 XC tenant
- Resource counts accumulate across ALL namespaces
- Querying quotas from any namespace returns identical tenant-wide data

Shows all quota limits with their current usage and percentage used. Quotas at
80%+ usage are flagged as WARNING, and quotas at or above 100% are flagged as
EXCEEDED.

AI assistants should check quota availability before deploying resources to ensure
deployment will not fail due to quota limits.`,
	Example: `  # Show tenant quota usage
  f5xcctl subscription quota

  # Show quota as JSON for automation
  f5xcctl subscription quota --output-format json

  # Show only object quotas
  f5xcctl subscription quota --type objects`,
	RunE: runSubscriptionQuota,
}

func init() {
	subscriptionCmd.AddCommand(subscriptionQuotaCmd)

	subscriptionQuotaCmd.Flags().StringVar(&quotaType, "type", "", "Filter by quota type: objects, resources, apis.")
}

func runSubscriptionQuota(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	client := GetSubscriptionClient()

	if client == nil {
		return fmt.Errorf("subscription client not initialized")
	}

	// Note: Quotas are tenant-level, not namespace-level
	quotaInfo, err := client.GetQuotaInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get quota info: %w", err)
	}

	// Filter by type if specified
	if quotaType != "" {
		switch strings.ToLower(quotaType) {
		case "objects":
			quotaInfo.Resources = nil
			quotaInfo.APIs = nil
		case "resources":
			quotaInfo.Objects = nil
			quotaInfo.APIs = nil
		case "apis":
			quotaInfo.Objects = nil
			quotaInfo.Resources = nil
		}
	}

	// Output based on format
	format := GetOutputFormatWithDefault("table")
	return outputQuota(quotaInfo, format)
}

func outputQuota(quotaInfo *subscription.QuotaUsageInfo, format string) error {
	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(quotaInfo)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(quotaInfo)
	default:
		return outputQuotaTable(quotaInfo)
	}
}

func outputQuotaTable(quotaInfo *subscription.QuotaUsageInfo) error {
	// Count statuses
	var okCount, warningCount, exceededCount int
	for _, q := range quotaInfo.Objects {
		switch q.Status {
		case "OK":
			okCount++
		case "WARNING":
			warningCount++
		case "EXCEEDED":
			exceededCount++
		}
	}

	// Print header
	fmt.Println("TENANT QUOTA USAGE")
	fmt.Println(strings.Repeat("=", 75))
	fmt.Printf("  OK: %d | Warning: %d | Exceeded: %d | Total: %d\n",
		okCount, warningCount, exceededCount, len(quotaInfo.Objects))
	fmt.Println()

	// Print object quotas
	if len(quotaInfo.Objects) > 0 {
		fmt.Println("OBJECT QUOTAS")
		fmt.Println(strings.Repeat("-", 75))
		fmt.Printf("  %-30s %-10s %-10s %-10s %-10s\n",
			"OBJECT TYPE", "USAGE", "LIMIT", "%USED", "STATUS")
		fmt.Println("  " + strings.Repeat("-", 70))

		// Sort by percentage descending (show most used first)
		objects := make([]subscription.QuotaItem, len(quotaInfo.Objects))
		copy(objects, quotaInfo.Objects)
		sort.Slice(objects, func(i, j int) bool {
			return objects[i].Percentage > objects[j].Percentage
		})

		for _, q := range objects {
			name := q.DisplayName
			if name == "" {
				name = q.Name
			}
			if len(name) > 29 {
				name = name[:26] + "..."
			}

			statusStr := getStatusString(q.Status)

			fmt.Printf("  %-30s %-10.0f %-10.0f %-10.0f%% %-10s\n",
				name, q.Usage, q.Limit, q.Percentage, statusStr)
		}
		fmt.Println()
	}

	// Print resource quotas
	if len(quotaInfo.Resources) > 0 {
		fmt.Println("RESOURCE QUOTAS")
		fmt.Println(strings.Repeat("-", 75))
		fmt.Printf("  %-30s %-10s %-10s %-10s %-10s\n",
			"RESOURCE", "USAGE", "LIMIT", "%USED", "STATUS")
		fmt.Println("  " + strings.Repeat("-", 70))

		for _, q := range quotaInfo.Resources {
			name := q.DisplayName
			if name == "" {
				name = q.Name
			}
			if len(name) > 29 {
				name = name[:26] + "..."
			}

			statusStr := getStatusString(q.Status)

			fmt.Printf("  %-30s %-10.0f %-10.0f %-10.0f%% %-10s\n",
				name, q.Usage, q.Limit, q.Percentage, statusStr)
		}
		fmt.Println()
	}

	// Print API rate limits
	if len(quotaInfo.APIs) > 0 {
		fmt.Println("API RATE LIMITS")
		fmt.Println(strings.Repeat("-", 75))
		fmt.Printf("  %-30s %-10s %-10s %-15s\n",
			"API", "RATE", "BURST", "UNIT")
		fmt.Println("  " + strings.Repeat("-", 70))

		for _, r := range quotaInfo.APIs {
			fmt.Printf("  %-30s %-10d %-10d %-15s\n",
				r.Name, r.Rate, r.Burst, r.Unit)
		}
		fmt.Println()
	}

	// Show warnings and hints
	if exceededCount > 0 || warningCount > 0 {
		fmt.Println("ALERTS")
		fmt.Println(strings.Repeat("-", 75))
		if exceededCount > 0 {
			fmt.Printf("  CRITICAL: %d quota(s) exceeded - new resource creation may fail\n", exceededCount)
		}
		if warningCount > 0 {
			fmt.Printf("  WARNING: %d quota(s) at 80%%+ usage - consider cleanup or upgrade\n", warningCount)
		}
		fmt.Println()
	}

	fmt.Println("HINTS")
	fmt.Println(strings.Repeat("-", 75))
	fmt.Println("  Use 'f5xcctl subscription validate --resource-type <type> --count <n>' to check")
	fmt.Println("  if you can create additional resources before deployment.")
	fmt.Println()

	return nil
}

func getStatusString(status string) string {
	switch status {
	case "OK":
		return "OK"
	case "WARNING":
		return "WARNING (!)"
	case "EXCEEDED":
		return "EXCEEDED (!!!)"
	default:
		return status
	}
}
