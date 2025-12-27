package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var summaryBrief bool

var cloudstatusSummaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Get complete status summary including components, incidents, and maintenance.",
	Long: `Get a comprehensive status summary of F5 Distributed Cloud services.

Includes:
- Overall status indicator
- All components and their current status
- Active and recent incidents
- Scheduled and ongoing maintenance windows

Use --brief for a condensed one-liner per section.`,
	Example: `  # Full summary
  xcsh cloudstatus summary

  # Brief summary
  xcsh cloudstatus summary --brief

  # JSON output for parsing
  xcsh cloudstatus summary --output-format json`,
	RunE: runCloudstatusSummary,
}

func init() {
	cloudstatusCmd.AddCommand(cloudstatusSummaryCmd)

	cloudstatusSummaryCmd.Flags().BoolVar(&summaryBrief, "brief", false, "Condensed one-liner per section.")
}

func runCloudstatusSummary(cmd *cobra.Command, args []string) error {
	client, err := requireCloudStatusClient()
	if err != nil {
		return err
	}

	resp, err := client.GetSummary()
	if err != nil {
		return fmt.Errorf("failed to get summary: %w", err)
	}

	format := GetOutputFormat()
	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(resp)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(resp)
	default:
		if summaryBrief {
			return printBriefSummary(resp)
		}
		return printFullSummary(resp)
	}
}

func printBriefSummary(resp interface{}) error {
	// Type assert to SummaryResponse
	summary, ok := resp.(*interface{})
	if !ok {
		// If we can't assert, use the actual type
		return printBriefSummaryTyped(resp)
	}
	_ = summary
	return printBriefSummaryTyped(resp)
}

func printBriefSummaryTyped(resp interface{}) error {
	client := GetCloudStatusClient()

	// Get fresh data
	summary, err := client.GetSummary()
	if err != nil {
		return err
	}

	// Count statistics
	operationalCount := 0
	degradedCount := 0
	for _, comp := range summary.Components {
		if !comp.Group {
			if comp.IsOperational() {
				operationalCount++
			} else {
				degradedCount++
			}
		}
	}

	activeIncidents := 0
	for _, inc := range summary.Incidents {
		if inc.IsActive() {
			activeIncidents++
		}
	}

	upcomingMaint := 0
	activeMaint := 0
	for _, maint := range summary.ScheduledMaintenances {
		if maint.IsUpcoming() {
			upcomingMaint++
		} else if maint.IsActive() {
			activeMaint++
		}
	}

	fmt.Printf("Status: %s (%s)\n", summary.Status.Indicator, summary.Status.Description)
	fmt.Printf("Components: %d operational, %d degraded\n", operationalCount, degradedCount)
	fmt.Printf("Incidents: %d active\n", activeIncidents)
	fmt.Printf("Maintenance: %d upcoming, %d active\n", upcomingMaint, activeMaint)

	return nil
}

func printFullSummary(resp interface{}) error {
	client := GetCloudStatusClient()

	// Get fresh data
	summary, err := client.GetSummary()
	if err != nil {
		return err
	}

	// Print overall status
	fmt.Println("=== OVERALL STATUS ===")
	fmt.Printf("Indicator: %s\n", summary.Status.Indicator)
	fmt.Printf("Description: %s\n", summary.Status.Description)
	fmt.Printf("Last Updated: %s\n", summary.Page.UpdatedAt.Format("2006-01-02 15:04:05 MST"))
	fmt.Println()

	// Print components summary
	fmt.Println("=== COMPONENTS ===")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "NAME\tSTATUS")

	// Print non-operational components first
	for _, comp := range summary.Components {
		if !comp.Group && comp.IsDegraded() {
			_, _ = fmt.Fprintf(w, "%s\t%s\n", comp.Name, comp.Status)
		}
	}
	_ = w.Flush()

	// Count operational
	operationalCount := 0
	for _, comp := range summary.Components {
		if !comp.Group && comp.IsOperational() {
			operationalCount++
		}
	}
	fmt.Printf("\n(%d components operational, not shown)\n", operationalCount)
	fmt.Println()

	// Print active incidents
	fmt.Println("=== ACTIVE INCIDENTS ===")
	hasActiveIncidents := false
	for _, inc := range summary.Incidents {
		if inc.IsActive() {
			hasActiveIncidents = true
			fmt.Printf("[%s] %s (Impact: %s)\n", inc.Status, inc.Name, inc.Impact)
			fmt.Printf("  Started: %s\n", inc.StartedAt.Format("2006-01-02 15:04:05 MST"))
			if len(inc.IncidentUpdates) > 0 {
				fmt.Printf("  Latest: %s\n", inc.IncidentUpdates[0].Body)
			}
			fmt.Println()
		}
	}
	if !hasActiveIncidents {
		fmt.Println("No active incidents")
		fmt.Println()
	}

	// Print scheduled maintenances
	fmt.Println("=== SCHEDULED MAINTENANCE ===")
	hasMaintenances := false
	for _, maint := range summary.ScheduledMaintenances {
		if !maint.IsCompleted() {
			hasMaintenances = true
			fmt.Printf("[%s] %s\n", maint.Status, maint.Name)
			fmt.Printf("  Scheduled: %s to %s\n",
				maint.ScheduledFor.Format("2006-01-02 15:04 MST"),
				maint.ScheduledUntil.Format("2006-01-02 15:04 MST"))
			fmt.Println()
		}
	}
	if !hasMaintenances {
		fmt.Println("No scheduled maintenance")
	}

	return nil
}
