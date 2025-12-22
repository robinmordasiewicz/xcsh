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

// Component filtering flags
var (
	componentsGroup          string
	componentsStatus         string
	componentsPop            bool
	componentsServices       bool
	componentsDegradedOnly   bool
	componentsWithComponents bool
)

var cloudstatusComponentsCmd = &cobra.Command{
	Use:   "components",
	Short: "Manage and view F5 Cloud Status components.",
	Long:  `View and filter F5 Distributed Cloud service components and their status.`,
	Example: `  # List all components
  f5xcctl cloudstatus components list

  # List degraded components only
  f5xcctl cloudstatus components list --degraded-only

  # List PoP components
  f5xcctl cloudstatus components list --pop

  # Get a specific component
  f5xcctl cloudstatus components get <component-id>

  # List component groups
  f5xcctl cloudstatus components groups`,
}

var cloudstatusComponentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all components with optional filtering.",
	Long: `List all F5 Distributed Cloud service components.

Supports filtering by:
- Group: Filter by component group name or ID
- Status: Filter by operational status
- PoP: Show only Point of Presence components
- Services: Show only service components
- Degraded: Show only non-operational components`,
	Example: `  # List all components
  f5xcctl cloudstatus components list

  # Filter by group
  f5xcctl cloudstatus components list --group "Services"

  # Filter by status
  f5xcctl cloudstatus components list --status degraded_performance

  # Show only degraded components
  f5xcctl cloudstatus components list --degraded-only

  # Show only PoP components
  f5xcctl cloudstatus components list --pop`,
	RunE: runComponentsList,
}

var cloudstatusComponentsGetCmd = &cobra.Command{
	Use:   "get <component-id>",
	Short: "Get details for a specific component.",
	Long:  `Get detailed information about a specific component by its ID or name.`,
	Example: `  # Get by ID
  f5xcctl cloudstatus components get ybcpdlwcdq67

  # Get by name (partial match)
  f5xcctl cloudstatus components get "Portal"`,
	Args: cobra.ExactArgs(1),
	RunE: runComponentsGet,
}

var cloudstatusComponentsGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List component groups.",
	Long:  `List all component groups and their hierarchy.`,
	Example: `  # List groups
  f5xcctl cloudstatus components groups

  # List groups with component counts
  f5xcctl cloudstatus components groups --with-components`,
	RunE: runComponentsGroups,
}

func init() {
	cloudstatusCmd.AddCommand(cloudstatusComponentsCmd)
	cloudstatusComponentsCmd.AddCommand(cloudstatusComponentsListCmd)
	cloudstatusComponentsCmd.AddCommand(cloudstatusComponentsGetCmd)
	cloudstatusComponentsCmd.AddCommand(cloudstatusComponentsGroupsCmd)

	// List flags
	cloudstatusComponentsListCmd.Flags().StringVar(&componentsGroup, "group", "", "Filter by group name or ID.")
	cloudstatusComponentsListCmd.Flags().StringVar(&componentsStatus, "status", "", "Filter by status (operational, degraded_performance, partial_outage, major_outage, under_maintenance).")
	cloudstatusComponentsListCmd.Flags().BoolVar(&componentsPop, "pop", false, "Show only PoP components.")
	cloudstatusComponentsListCmd.Flags().BoolVar(&componentsServices, "services", false, "Show only service components.")
	cloudstatusComponentsListCmd.Flags().BoolVar(&componentsDegradedOnly, "degraded-only", false, "Show only non-operational components.")

	// Groups flags
	cloudstatusComponentsGroupsCmd.Flags().BoolVar(&componentsWithComponents, "with-components", false, "Include component count per group.")

	// Register completions
	_ = cloudstatusComponentsListCmd.RegisterFlagCompletionFunc("status", completeComponentStatus)
	_ = cloudstatusComponentsListCmd.RegisterFlagCompletionFunc("group", completeComponentGroup)
	cloudstatusComponentsGetCmd.ValidArgsFunction = completeComponentID
}

func runComponentsList(cmd *cobra.Command, args []string) error {
	client, err := requireCloudStatusClient()
	if err != nil {
		return err
	}

	resp, err := client.GetComponents()
	if err != nil {
		return fmt.Errorf("failed to get components: %w", err)
	}

	// Apply filters
	components := resp.Components

	// Filter out groups for display (unless specifically requested)
	filtered := []cloudstatus.Component{}
	for _, comp := range components {
		if comp.Group {
			continue // Skip groups in list view
		}

		// Apply group filter
		if componentsGroup != "" {
			if comp.GroupID == nil {
				continue
			}
			// Check if group matches by ID or name
			groupMatches := false
			for _, g := range components {
				if g.Group && g.ID == *comp.GroupID {
					if strings.EqualFold(g.Name, componentsGroup) || g.ID == componentsGroup {
						groupMatches = true
						break
					}
				}
			}
			if !groupMatches {
				continue
			}
		}

		// Apply status filter
		if componentsStatus != "" && comp.Status != componentsStatus {
			continue
		}

		// Apply PoP filter
		if componentsPop {
			if !strings.Contains(strings.ToLower(comp.Description), "pop") &&
				!strings.Contains(strings.ToLower(comp.Description), "edge") {
				continue
			}
		}

		// Apply services filter (components in Services group)
		if componentsServices {
			if comp.GroupID == nil {
				continue
			}
			isService := false
			for _, g := range components {
				if g.Group && g.ID == *comp.GroupID {
					if strings.Contains(strings.ToLower(g.Name), "services") ||
						strings.Contains(strings.ToLower(g.Name), "customer support") {
						isService = true
						break
					}
				}
			}
			if !isService {
				continue
			}
		}

		// Apply degraded filter
		if componentsDegradedOnly && comp.IsOperational() {
			continue
		}

		filtered = append(filtered, comp)
	}

	format := GetOutputFormat()
	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(filtered)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(filtered)
	case "wide":
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		_, _ = fmt.Fprintln(w, "ID\tNAME\tSTATUS\tDESCRIPTION")
		for _, comp := range filtered {
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", comp.ID, comp.Name, comp.Status, comp.Description)
		}
		return w.Flush()
	default:
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		_, _ = fmt.Fprintln(w, "NAME\tSTATUS")
		for _, comp := range filtered {
			_, _ = fmt.Fprintf(w, "%s\t%s\n", comp.Name, comp.Status)
		}
		return w.Flush()
	}
}

func runComponentsGet(cmd *cobra.Command, args []string) error {
	client, err := requireCloudStatusClient()
	if err != nil {
		return err
	}

	componentID := args[0]

	// First try to get by ID
	resp, err := client.GetComponent(componentID)
	if err != nil {
		// If that fails, search by name
		allResp, searchErr := client.GetComponents()
		if searchErr != nil {
			return fmt.Errorf("failed to get component: %w", err)
		}

		var found *cloudstatus.Component
		searchLower := strings.ToLower(componentID)
		for i, comp := range allResp.Components {
			if strings.Contains(strings.ToLower(comp.Name), searchLower) {
				found = &allResp.Components[i]
				break
			}
		}

		if found == nil {
			return fmt.Errorf("component not found: %s", componentID)
		}

		// Get the full component details
		resp, err = client.GetComponent(found.ID)
		if err != nil {
			// Just use what we found
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
				printComponentDetails(found)
				return nil
			}
		}
	}

	format := GetOutputFormat()
	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(resp.Component)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(resp.Component)
	default:
		printComponentDetails(&resp.Component)
		return nil
	}
}

func printComponentDetails(comp *cloudstatus.Component) {
	fmt.Printf("ID:          %s\n", comp.ID)
	fmt.Printf("Name:        %s\n", comp.Name)
	fmt.Printf("Status:      %s\n", comp.Status)
	if comp.Description != "" {
		fmt.Printf("Description: %s\n", comp.Description)
	}
	if comp.GroupID != nil {
		fmt.Printf("Group ID:    %s\n", *comp.GroupID)
	}
	fmt.Printf("Position:    %d\n", comp.Position)
	fmt.Printf("Created:     %s\n", comp.CreatedAt.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("Updated:     %s\n", comp.UpdatedAt.Format("2006-01-02 15:04:05 MST"))
}

func runComponentsGroups(cmd *cobra.Command, args []string) error {
	client, err := requireCloudStatusClient()
	if err != nil {
		return err
	}

	groups, err := client.GetComponentGroups()
	if err != nil {
		return fmt.Errorf("failed to get component groups: %w", err)
	}

	format := GetOutputFormat()
	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(groups)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(groups)
	default:
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		if componentsWithComponents {
			_, _ = fmt.Fprintln(w, "NAME\tID\tCOMPONENTS")
			for _, group := range groups {
				_, _ = fmt.Fprintf(w, "%s\t%s\t%d\n", group.Name, group.ID, group.ComponentCount)
			}
		} else {
			_, _ = fmt.Fprintln(w, "NAME\tID")
			for _, group := range groups {
				_, _ = fmt.Fprintf(w, "%s\t%s\n", group.Name, group.ID)
			}
		}
		return w.Flush()
	}
}
