package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/robinmordasiewicz/f5xcctl/pkg/cloudstatus"
)

// PoP filtering flags
var popsRegion string

var cloudstatusPopsCmd = &cobra.Command{
	Use:   "pops",
	Short: "View Point of Presence (PoP) status.",
	Long: `View status of F5 Distributed Cloud Points of Presence (PoPs).

PoPs are edge locations distributed globally that provide low-latency
access to F5 XC services. Available regions:
- north-america
- south-america
- europe
- asia
- oceania
- middle-east`,
	Example: `  # List all PoPs
  f5xcctl cloudstatus pops list

  # List PoPs in a specific region
  f5xcctl cloudstatus pops list --region north-america

  # Get regional status summary
  f5xcctl cloudstatus pops status

  # Get status for a specific region
  f5xcctl cloudstatus pops status --region europe`,
}

var cloudstatusPopsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List PoP locations.",
	Long:  `List all Point of Presence (PoP) locations and their current status.`,
	Example: `  # List all PoPs
  f5xcctl cloudstatus pops list

  # List PoPs by region
  f5xcctl cloudstatus pops list --region north-america
  f5xcctl cloudstatus pops list --region europe`,
	RunE: runPopsList,
}

var cloudstatusPopsStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get aggregated status by region.",
	Long:  `Get aggregated status summary for each geographic region.`,
	Example: `  # All regions summary
  f5xcctl cloudstatus pops status

  # Specific region
  f5xcctl cloudstatus pops status --region asia`,
	RunE: runPopsStatus,
}

func init() {
	cloudstatusCmd.AddCommand(cloudstatusPopsCmd)
	cloudstatusPopsCmd.AddCommand(cloudstatusPopsListCmd)
	cloudstatusPopsCmd.AddCommand(cloudstatusPopsStatusCmd)

	// List flags
	cloudstatusPopsListCmd.Flags().StringVar(&popsRegion, "region", "", "Filter by region (north-america, south-america, europe, asia, oceania, middle-east).")

	// Status flags
	cloudstatusPopsStatusCmd.Flags().StringVar(&popsRegion, "region", "", "Filter by region (north-america, south-america, europe, asia, oceania, middle-east).")

	// Register completions
	_ = cloudstatusPopsListCmd.RegisterFlagCompletionFunc("region", completeRegion)
	_ = cloudstatusPopsStatusCmd.RegisterFlagCompletionFunc("region", completeRegion)
}

func runPopsList(cmd *cobra.Command, args []string) error {
	client, err := requireCloudStatusClient()
	if err != nil {
		return err
	}

	pops, err := client.GetPoPs()
	if err != nil {
		return fmt.Errorf("failed to get PoPs: %w", err)
	}

	// Get groups for region detection
	groups, err := client.GetComponentGroups()
	if err != nil {
		return fmt.Errorf("failed to get component groups: %w", err)
	}

	// Apply region filter
	if popsRegion != "" {
		filtered := []cloudstatus.Component{}
		for _, pop := range pops {
			region := cloudstatus.DetectRegion(pop, groups)
			if strings.EqualFold(region, popsRegion) {
				filtered = append(filtered, pop)
			}
		}
		pops = filtered
	}

	format := GetOutputFormat()
	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(pops)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(pops)
	case "wide":
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		_, _ = fmt.Fprintln(w, "NAME\tSTATUS\tSITE CODE\tREGION")
		for _, pop := range pops {
			siteCode := cloudstatus.ExtractSiteCode(pop.Name)
			if siteCode == "" {
				siteCode = cloudstatus.ExtractSiteCode(pop.Description)
			}
			region := cloudstatus.DetectRegion(pop, groups)
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", pop.Name, pop.Status, siteCode, region)
		}
		return w.Flush()
	default:
		if len(pops) == 0 {
			fmt.Println("No PoPs found")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		_, _ = fmt.Fprintln(w, "NAME\tSTATUS")
		for _, pop := range pops {
			_, _ = fmt.Fprintf(w, "%s\t%s\n", pop.Name, pop.Status)
		}
		return w.Flush()
	}
}

func runPopsStatus(cmd *cobra.Command, args []string) error {
	client, err := requireCloudStatusClient()
	if err != nil {
		return err
	}

	statuses, err := client.GetRegionalStatus()
	if err != nil {
		return fmt.Errorf("failed to get regional status: %w", err)
	}

	// Apply region filter
	if popsRegion != "" {
		filtered := []cloudstatus.RegionalStatus{}
		for _, status := range statuses {
			if strings.EqualFold(status.Region.Name, popsRegion) ||
				strings.EqualFold(status.Region.DisplayName, popsRegion) {
				filtered = append(filtered, status)
			}
		}
		statuses = filtered
	}

	format := GetOutputFormat()
	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(statuses)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(statuses)
	default:
		if len(statuses) == 0 {
			fmt.Println("No regional status data found")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		_, _ = fmt.Fprintln(w, "REGION\tSTATUS\tOPERATIONAL\tDEGRADED\tTOTAL")
		for _, status := range statuses {
			if status.TotalCount == 0 {
				continue // Skip regions with no PoPs detected
			}
			_, _ = fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%d\n",
				status.Region.DisplayName,
				status.OverallStatus,
				status.OperationalCount,
				status.DegradedCount,
				status.TotalCount)
		}
		return w.Flush()
	}
}
