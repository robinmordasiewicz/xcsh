package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/robinmordasiewicz/f5xcctl/pkg/cloudstatus"
)

// Incident filtering flags
var (
	incidentsStatus string
	incidentsImpact string
	incidentsSince  string
	incidentsLimit  int
)

var cloudstatusIncidentsCmd = &cobra.Command{
	Use:   "incidents",
	Short: "View and manage F5 Cloud Status incidents.",
	Long:  `View incident history and details for F5 Distributed Cloud services.`,
	Example: `  # List all incidents
  f5xcctl cloudstatus incidents list

  # List active (unresolved) incidents
  f5xcctl cloudstatus incidents active

  # Get incident details
  f5xcctl cloudstatus incidents get <incident-id>

  # Show incident timeline
  f5xcctl cloudstatus incidents updates <incident-id>`,
}

var cloudstatusIncidentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all incidents.",
	Long: `List all incidents with optional filtering.

Supports filtering by:
- Status: investigating, identified, monitoring, resolved, postmortem
- Impact: none, minor, major, critical
- Since: Time range (1h, 1d, 7d, 30d)
- Limit: Maximum number of results`,
	Example: `  # List all incidents
  f5xcctl cloudstatus incidents list

  # Filter by status
  f5xcctl cloudstatus incidents list --status monitoring

  # Filter by impact
  f5xcctl cloudstatus incidents list --impact major

  # Recent incidents (last 7 days)
  f5xcctl cloudstatus incidents list --since 7d

  # Limit results
  f5xcctl cloudstatus incidents list --limit 10`,
	RunE: runIncidentsList,
}

var cloudstatusIncidentsActiveCmd = &cobra.Command{
	Use:   "active",
	Short: "List only unresolved incidents.",
	Long:  `List all currently active (unresolved) incidents.`,
	Example: `  f5xcctl cloudstatus incidents active
  f5xcctl cloudstatus incidents active --output-format json`,
	RunE: runIncidentsActive,
}

var cloudstatusIncidentsGetCmd = &cobra.Command{
	Use:     "get <incident-id>",
	Short:   "Get details for a specific incident.",
	Long:    `Get detailed information about a specific incident by its ID.`,
	Example: `  f5xcctl cloudstatus incidents get kcxnsw71xmwp`,
	Args:    cobra.ExactArgs(1),
	RunE:    runIncidentsGet,
}

var cloudstatusIncidentsUpdatesCmd = &cobra.Command{
	Use:     "updates <incident-id>",
	Short:   "Show incident update timeline.",
	Long:    `Display the chronological timeline of updates for a specific incident.`,
	Example: `  f5xcctl cloudstatus incidents updates kcxnsw71xmwp`,
	Args:    cobra.ExactArgs(1),
	RunE:    runIncidentsUpdates,
}

func init() {
	cloudstatusCmd.AddCommand(cloudstatusIncidentsCmd)
	cloudstatusIncidentsCmd.AddCommand(cloudstatusIncidentsListCmd)
	cloudstatusIncidentsCmd.AddCommand(cloudstatusIncidentsActiveCmd)
	cloudstatusIncidentsCmd.AddCommand(cloudstatusIncidentsGetCmd)
	cloudstatusIncidentsCmd.AddCommand(cloudstatusIncidentsUpdatesCmd)

	// List flags
	cloudstatusIncidentsListCmd.Flags().StringVar(&incidentsStatus, "status", "", "Filter by status (investigating, identified, monitoring, resolved, postmortem).")
	cloudstatusIncidentsListCmd.Flags().StringVar(&incidentsImpact, "impact", "", "Filter by impact level (none, minor, major, critical).")
	cloudstatusIncidentsListCmd.Flags().StringVar(&incidentsSince, "since", "", "Time filter (1h, 1d, 7d, 30d).")
	cloudstatusIncidentsListCmd.Flags().IntVar(&incidentsLimit, "limit", 0, "Limit number of results.")

	// Register completions
	_ = cloudstatusIncidentsListCmd.RegisterFlagCompletionFunc("status", completeIncidentStatus)
	_ = cloudstatusIncidentsListCmd.RegisterFlagCompletionFunc("impact", completeIncidentImpact)
	cloudstatusIncidentsGetCmd.ValidArgsFunction = completeIncidentID
	cloudstatusIncidentsUpdatesCmd.ValidArgsFunction = completeIncidentID
}

func runIncidentsList(cmd *cobra.Command, args []string) error {
	client := GetCloudStatusClient()
	if client == nil {
		return fmt.Errorf("cloudstatus client not initialized")
	}

	resp, err := client.GetIncidents()
	if err != nil {
		return fmt.Errorf("failed to get incidents: %w", err)
	}

	incidents := resp.Incidents

	// Apply status filter
	if incidentsStatus != "" {
		incidents = cloudstatus.FilterIncidentsByStatus(incidents, incidentsStatus)
	}

	// Apply impact filter
	if incidentsImpact != "" {
		incidents = cloudstatus.FilterIncidentsByImpact(incidents, incidentsImpact)
	}

	// Apply time filter
	if incidentsSince != "" {
		since, err := parseDuration(incidentsSince)
		if err != nil {
			return fmt.Errorf("invalid --since value: %w", err)
		}
		sinceTime := time.Now().Add(-since)
		incidents = cloudstatus.FilterIncidentsSince(incidents, sinceTime)
	}

	// Apply limit
	if incidentsLimit > 0 && len(incidents) > incidentsLimit {
		incidents = incidents[:incidentsLimit]
	}

	return outputIncidents(incidents)
}

func runIncidentsActive(cmd *cobra.Command, args []string) error {
	client := GetCloudStatusClient()
	if client == nil {
		return fmt.Errorf("cloudstatus client not initialized")
	}

	resp, err := client.GetUnresolvedIncidents()
	if err != nil {
		return fmt.Errorf("failed to get unresolved incidents: %w", err)
	}

	return outputIncidents(resp.Incidents)
}

func runIncidentsGet(cmd *cobra.Command, args []string) error {
	client := GetCloudStatusClient()
	if client == nil {
		return fmt.Errorf("cloudstatus client not initialized")
	}

	incidentID := args[0]

	// Get all incidents and find the one we're looking for
	resp, err := client.GetIncidents()
	if err != nil {
		return fmt.Errorf("failed to get incidents: %w", err)
	}

	var found *cloudstatus.Incident
	for i, inc := range resp.Incidents {
		if inc.ID == incidentID || strings.Contains(strings.ToLower(inc.Name), strings.ToLower(incidentID)) {
			found = &resp.Incidents[i]
			break
		}
	}

	if found == nil {
		return fmt.Errorf("incident not found: %s", incidentID)
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
		printIncidentDetails(found)
		return nil
	}
}

func runIncidentsUpdates(cmd *cobra.Command, args []string) error {
	client := GetCloudStatusClient()
	if client == nil {
		return fmt.Errorf("cloudstatus client not initialized")
	}

	incidentID := args[0]

	// Get all incidents and find the one we're looking for
	resp, err := client.GetIncidents()
	if err != nil {
		return fmt.Errorf("failed to get incidents: %w", err)
	}

	var found *cloudstatus.Incident
	for i, inc := range resp.Incidents {
		if inc.ID == incidentID {
			found = &resp.Incidents[i]
			break
		}
	}

	if found == nil {
		return fmt.Errorf("incident not found: %s", incidentID)
	}

	format := GetOutputFormat()
	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(found.IncidentUpdates)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(found.IncidentUpdates)
	default:
		fmt.Printf("=== INCIDENT: %s ===\n", found.Name)
		fmt.Printf("Status: %s | Impact: %s\n\n", found.Status, found.Impact)
		fmt.Println("--- Update Timeline ---")
		for _, update := range found.IncidentUpdates {
			fmt.Printf("\n[%s] %s\n", update.DisplayAt.Format("2006-01-02 15:04:05 MST"), update.Status)
			fmt.Printf("%s\n", update.Body)
			if len(update.AffectedComponents) > 0 {
				fmt.Printf("Affected: ")
				names := []string{}
				for _, ac := range update.AffectedComponents {
					names = append(names, ac.Name)
				}
				fmt.Printf("%s\n", strings.Join(names, ", "))
			}
		}
		return nil
	}
}

func outputIncidents(incidents []cloudstatus.Incident) error {
	format := GetOutputFormat()
	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(incidents)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(incidents)
	case "wide":
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		_, _ = fmt.Fprintln(w, "ID\tNAME\tSTATUS\tIMPACT\tSTARTED")
		for _, inc := range incidents {
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				inc.ID,
				inc.Name,
				inc.Status,
				inc.Impact,
				inc.StartedAt.Format("2006-01-02 15:04"))
		}
		return w.Flush()
	default:
		if len(incidents) == 0 {
			fmt.Println("No incidents found")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		_, _ = fmt.Fprintln(w, "NAME\tSTATUS\tIMPACT")
		for _, inc := range incidents {
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", inc.Name, inc.Status, inc.Impact)
		}
		return w.Flush()
	}
}

func printIncidentDetails(inc *cloudstatus.Incident) {
	fmt.Printf("ID:        %s\n", inc.ID)
	fmt.Printf("Name:      %s\n", inc.Name)
	fmt.Printf("Status:    %s\n", inc.Status)
	fmt.Printf("Impact:    %s\n", inc.Impact)
	fmt.Printf("Started:   %s\n", inc.StartedAt.Format("2006-01-02 15:04:05 MST"))
	if inc.ResolvedAt != nil {
		fmt.Printf("Resolved:  %s\n", inc.ResolvedAt.Format("2006-01-02 15:04:05 MST"))
	}
	fmt.Printf("Shortlink: %s\n", inc.Shortlink)

	if len(inc.Components) > 0 {
		fmt.Printf("\nAffected Components:\n")
		for _, comp := range inc.Components {
			fmt.Printf("  - %s (%s)\n", comp.Name, comp.Status)
		}
	}

	if len(inc.IncidentUpdates) > 0 {
		fmt.Printf("\nLatest Update:\n")
		update := inc.IncidentUpdates[0]
		fmt.Printf("  [%s] %s\n", update.DisplayAt.Format("2006-01-02 15:04:05 MST"), update.Status)
		fmt.Printf("  %s\n", update.Body)
	}
}

func parseDuration(s string) (time.Duration, error) {
	// Handle special formats like 1d, 7d, 30d
	if strings.HasSuffix(s, "d") {
		days := strings.TrimSuffix(s, "d")
		var d int
		_, err := fmt.Sscanf(days, "%d", &d)
		if err != nil {
			return 0, fmt.Errorf("invalid duration: %s", s)
		}
		return time.Duration(d) * 24 * time.Hour, nil
	}
	return time.ParseDuration(s)
}
