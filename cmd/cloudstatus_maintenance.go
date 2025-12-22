package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/robinmordasiewicz/xcsh/pkg/cloudstatus"
)

// Maintenance filtering flags
var maintenanceStatus string

var cloudstatusMaintenanceCmd = &cobra.Command{
	Use:     "maintenance",
	Aliases: []string{"maint"},
	Short:   "View scheduled maintenance windows.",
	Long:    `View scheduled maintenance windows for F5 Distributed Cloud services.`,
	Example: `  # List all maintenance windows
  f5xcctl cloudstatus maintenance list

  # List upcoming maintenance
  f5xcctl cloudstatus maintenance upcoming

  # List active (in-progress) maintenance
  f5xcctl cloudstatus maintenance active

  # Get maintenance details
  f5xcctl cloudstatus maintenance get <maintenance-id>`,
}

var cloudstatusMaintenanceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all scheduled maintenance windows.",
	Long: `List all scheduled maintenance windows with optional filtering.

Status values:
- scheduled: Maintenance is planned but not yet started
- in_progress: Maintenance is currently underway
- verifying: Maintenance is being verified
- completed: Maintenance has finished`,
	Example: `  # List all maintenance
  f5xcctl cloudstatus maintenance list

  # Filter by status
  f5xcctl cloudstatus maintenance list --status scheduled`,
	RunE: runMaintenanceList,
}

var cloudstatusMaintenanceUpcomingCmd = &cobra.Command{
	Use:   "upcoming",
	Short: "List upcoming maintenance windows.",
	Long:  `List only upcoming (scheduled but not started) maintenance windows.`,
	Example: `  f5xcctl cloudstatus maintenance upcoming
  f5xcctl cloudstatus maintenance upcoming --output-format json`,
	RunE: runMaintenanceUpcoming,
}

var cloudstatusMaintenanceActiveCmd = &cobra.Command{
	Use:     "active",
	Short:   "List in-progress maintenance windows.",
	Long:    `List maintenance windows that are currently in progress.`,
	Example: `  f5xcctl cloudstatus maintenance active`,
	RunE:    runMaintenanceActive,
}

var cloudstatusMaintenanceGetCmd = &cobra.Command{
	Use:     "get <maintenance-id>",
	Short:   "Get details for a specific maintenance window.",
	Long:    `Get detailed information about a specific maintenance window by its ID.`,
	Example: `  f5xcctl cloudstatus maintenance get xp5l86wjjzyy`,
	Args:    cobra.ExactArgs(1),
	RunE:    runMaintenanceGet,
}

func init() {
	cloudstatusCmd.AddCommand(cloudstatusMaintenanceCmd)
	cloudstatusMaintenanceCmd.AddCommand(cloudstatusMaintenanceListCmd)
	cloudstatusMaintenanceCmd.AddCommand(cloudstatusMaintenanceUpcomingCmd)
	cloudstatusMaintenanceCmd.AddCommand(cloudstatusMaintenanceActiveCmd)
	cloudstatusMaintenanceCmd.AddCommand(cloudstatusMaintenanceGetCmd)

	// List flags
	cloudstatusMaintenanceListCmd.Flags().StringVar(&maintenanceStatus, "status", "", "Filter by status (scheduled, in_progress, verifying, completed).")

	// Register completions
	_ = cloudstatusMaintenanceListCmd.RegisterFlagCompletionFunc("status", completeMaintenanceStatus)
	cloudstatusMaintenanceGetCmd.ValidArgsFunction = completeMaintenanceID
}

func runMaintenanceList(cmd *cobra.Command, args []string) error {
	client, err := requireCloudStatusClient()
	if err != nil {
		return err
	}

	resp, err := client.GetMaintenances()
	if err != nil {
		return fmt.Errorf("failed to get maintenances: %w", err)
	}

	maintenances := resp.ScheduledMaintenances

	// Apply status filter
	if maintenanceStatus != "" {
		maintenances = cloudstatus.FilterMaintenancesByStatus(maintenances, maintenanceStatus)
	}

	return outputMaintenances(maintenances)
}

func runMaintenanceUpcoming(cmd *cobra.Command, args []string) error {
	client, err := requireCloudStatusClient()
	if err != nil {
		return err
	}

	resp, err := client.GetUpcomingMaintenances()
	if err != nil {
		return fmt.Errorf("failed to get upcoming maintenances: %w", err)
	}

	return outputMaintenances(resp.ScheduledMaintenances)
}

func runMaintenanceActive(cmd *cobra.Command, args []string) error {
	client, err := requireCloudStatusClient()
	if err != nil {
		return err
	}

	resp, err := client.GetMaintenances()
	if err != nil {
		return fmt.Errorf("failed to get maintenances: %w", err)
	}

	active := cloudstatus.GetActiveMaintenances(resp.ScheduledMaintenances)

	return outputMaintenances(active)
}

func runMaintenanceGet(cmd *cobra.Command, args []string) error {
	client, err := requireCloudStatusClient()
	if err != nil {
		return err
	}

	maintenanceID := args[0]

	// Get all maintenances and find the one we're looking for
	resp, err := client.GetMaintenances()
	if err != nil {
		return fmt.Errorf("failed to get maintenances: %w", err)
	}

	var found *cloudstatus.ScheduledMaintenance
	for i, maint := range resp.ScheduledMaintenances {
		if maint.ID == maintenanceID {
			found = &resp.ScheduledMaintenances[i]
			break
		}
	}

	if found == nil {
		return fmt.Errorf("maintenance not found: %s", maintenanceID)
	}

	format := GetOutputFormat()
	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(found)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(found)
	default:
		printMaintenanceDetails(found)
		return nil
	}
}

func outputMaintenances(maintenances []cloudstatus.ScheduledMaintenance) error {
	format := GetOutputFormat()
	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(maintenances)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(maintenances)
	case "wide":
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		_, _ = fmt.Fprintln(w, "ID\tNAME\tSTATUS\tSCHEDULED FOR\tSCHEDULED UNTIL")
		for _, maint := range maintenances {
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				maint.ID,
				maint.Name,
				maint.Status,
				maint.ScheduledFor.Format("2006-01-02 15:04 MST"),
				maint.ScheduledUntil.Format("2006-01-02 15:04 MST"))
		}
		return w.Flush()
	default:
		if len(maintenances) == 0 {
			fmt.Println("No maintenance windows found")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		_, _ = fmt.Fprintln(w, "NAME\tSTATUS\tSCHEDULED")
		for _, maint := range maintenances {
			scheduled := fmt.Sprintf("%s - %s",
				maint.ScheduledFor.Format("Jan 2 15:04"),
				maint.ScheduledUntil.Format("Jan 2 15:04 MST"))
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", maint.Name, maint.Status, scheduled)
		}
		return w.Flush()
	}
}

func printMaintenanceDetails(maint *cloudstatus.ScheduledMaintenance) {
	fmt.Printf("ID:              %s\n", maint.ID)
	fmt.Printf("Name:            %s\n", maint.Name)
	fmt.Printf("Status:          %s\n", maint.Status)
	fmt.Printf("Impact:          %s\n", maint.Impact)
	fmt.Printf("Scheduled For:   %s\n", maint.ScheduledFor.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("Scheduled Until: %s\n", maint.ScheduledUntil.Format("2006-01-02 15:04:05 MST"))
	if maint.StartedAt != nil {
		fmt.Printf("Started At:      %s\n", maint.StartedAt.Format("2006-01-02 15:04:05 MST"))
	}
	if maint.ResolvedAt != nil {
		fmt.Printf("Completed At:    %s\n", maint.ResolvedAt.Format("2006-01-02 15:04:05 MST"))
	}
	fmt.Printf("Shortlink:       %s\n", maint.Shortlink)

	if len(maint.Components) > 0 {
		fmt.Printf("\nAffected Components:\n")
		for _, comp := range maint.Components {
			fmt.Printf("  - %s\n", comp.Name)
		}
	}

	if len(maint.IncidentUpdates) > 0 {
		fmt.Printf("\nUpdates:\n")
		for _, update := range maint.IncidentUpdates {
			fmt.Printf("  [%s] %s\n", update.DisplayAt.Format("2006-01-02 15:04:05 MST"), update.Status)
			// Indent the body
			lines := strings.Split(update.Body, "\n")
			for _, line := range lines {
				fmt.Printf("    %s\n", line)
			}
		}
	}
}
