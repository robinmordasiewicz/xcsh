package cmd

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/robinmordasiewicz/f5xcctl/pkg/cloudstatus"
)

// Completion client with longer cache TTL for tab completion
var completionClient *cloudstatus.Client

// getCompletionClient returns a cloudstatus client optimized for tab completion.
// Uses a 5-minute cache TTL and 3-second timeout to avoid blocking the shell.
func getCompletionClient() *cloudstatus.Client {
	if completionClient == nil {
		completionClient = cloudstatus.NewClient(
			cloudstatus.WithTimeout(3*time.Second),
			cloudstatus.WithCache(5*time.Minute),
		)
	}
	return completionClient
}

// Static completion functions - no API calls required

// completeComponentStatus provides completion for component status values.
func completeComponentStatus(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		cloudstatus.ComponentOperational + "\tFully operational",
		cloudstatus.ComponentDegradedPerformance + "\tDegraded performance",
		cloudstatus.ComponentPartialOutage + "\tPartial outage",
		cloudstatus.ComponentMajorOutage + "\tMajor outage",
		cloudstatus.ComponentUnderMaintenance + "\tUnder maintenance",
	}, cobra.ShellCompDirectiveNoFileComp
}

// completeIncidentStatus provides completion for incident status values.
func completeIncidentStatus(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		cloudstatus.IncidentInvestigating + "\tInvestigating the issue",
		cloudstatus.IncidentIdentified + "\tIssue identified",
		cloudstatus.IncidentMonitoring + "\tMonitoring the fix",
		cloudstatus.IncidentResolved + "\tIssue resolved",
		cloudstatus.IncidentPostmortem + "\tPostmortem in progress",
	}, cobra.ShellCompDirectiveNoFileComp
}

// completeIncidentImpact provides completion for incident impact values.
func completeIncidentImpact(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		cloudstatus.ImpactNone + "\tNo impact",
		cloudstatus.ImpactMinor + "\tMinor impact",
		cloudstatus.ImpactMajor + "\tMajor impact",
		cloudstatus.ImpactCritical + "\tCritical impact",
	}, cobra.ShellCompDirectiveNoFileComp
}

// completeMaintenanceStatus provides completion for maintenance status values.
func completeMaintenanceStatus(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		cloudstatus.MaintenanceScheduled + "\tMaintenance scheduled",
		cloudstatus.MaintenanceInProgress + "\tMaintenance in progress",
		cloudstatus.MaintenanceVerifying + "\tVerifying maintenance completion",
		cloudstatus.MaintenanceCompleted + "\tMaintenance completed",
	}, cobra.ShellCompDirectiveNoFileComp
}

// completeRegion provides completion for region values.
func completeRegion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	completions := make([]string, len(cloudstatus.PredefinedRegions))
	for i, region := range cloudstatus.PredefinedRegions {
		completions[i] = region.ID + "\t" + region.DisplayName
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// Dynamic completion functions - require API calls

// completeComponentID provides completion for component IDs.
func completeComponentID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Only complete first argument
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client := getCompletionClient()
	resp, err := client.GetComponents()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	completions := []string{}
	for _, comp := range resp.Components {
		// Skip groups - only show actual components
		if comp.Group {
			continue
		}
		completions = append(completions, comp.ID+"\t"+comp.Name)
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// completeIncidentID provides completion for incident IDs.
func completeIncidentID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Only complete first argument
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client := getCompletionClient()
	resp, err := client.GetIncidents()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	completions := []string{}
	for _, inc := range resp.Incidents {
		// Truncate name if too long
		name := inc.Name
		if len(name) > 50 {
			name = name[:47] + "..."
		}
		completions = append(completions, inc.ID+"\t"+name+" ("+inc.Status+")")
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// completeMaintenanceID provides completion for maintenance IDs.
func completeMaintenanceID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Only complete first argument
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client := getCompletionClient()
	resp, err := client.GetMaintenances()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	completions := []string{}
	for _, maint := range resp.ScheduledMaintenances {
		// Truncate name if too long
		name := maint.Name
		if len(name) > 50 {
			name = name[:47] + "..."
		}
		completions = append(completions, maint.ID+"\t"+name+" ("+maint.Status+")")
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// completeWatchComponents provides completion for the watch --components flag.
// Supports comma-separated values.
func completeWatchComponents(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client := getCompletionClient()
	resp, err := client.GetComponents()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	completions := []string{}
	for _, comp := range resp.Components {
		// Skip groups - only show actual components
		if comp.Group {
			continue
		}
		completions = append(completions, comp.ID+"\t"+comp.Name)
	}

	// Allow comma-separated values (no space after completion)
	return completions, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
}

// completeComponentGroup provides completion for component group filtering.
func completeComponentGroup(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client := getCompletionClient()
	groups, err := client.GetComponentGroups()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	completions := []string{}
	for _, group := range groups {
		completions = append(completions, group.ID+"\t"+group.Name)
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}
