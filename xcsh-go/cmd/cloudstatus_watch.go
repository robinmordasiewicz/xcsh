package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/robinmordasiewicz/xcsh/pkg/cloudstatus"
)

// Watch flags
var (
	watchInterval     int
	watchComponents   string
	watchExitOnChange bool
	watchNoClear      bool
)

var cloudstatusWatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Real-time status monitoring.",
	Long: `Monitor F5 Distributed Cloud status in real-time.

Continuously polls the status API and displays updates. Useful for:
- Monitoring during deployments
- Watching for status changes during incidents
- Real-time dashboards

Press Ctrl+C to stop monitoring.`,
	Example: `  # Default monitoring (60s interval)
  xcsh cloudstatus watch

  # Faster polling (30s interval)
  xcsh cloudstatus watch --interval 30

  # Watch specific components
  xcsh cloudstatus watch --components "Portal,DNS"

  # Exit when status changes (for alerting)
  xcsh cloudstatus watch --exit-on-change

  # Keep history visible (no screen clear)
  xcsh cloudstatus watch --no-clear`,
	RunE: runWatch,
}

func init() {
	cloudstatusCmd.AddCommand(cloudstatusWatchCmd)

	cloudstatusWatchCmd.Flags().IntVar(&watchInterval, "interval", 60, "Polling interval in seconds.")
	cloudstatusWatchCmd.Flags().StringVar(&watchComponents, "components", "", "Watch specific components (comma-separated names).")
	cloudstatusWatchCmd.Flags().BoolVar(&watchExitOnChange, "exit-on-change", false, "Exit when status changes.")
	cloudstatusWatchCmd.Flags().BoolVar(&watchNoClear, "no-clear", false, "Don't clear screen between updates.")

	// Register completions
	_ = cloudstatusWatchCmd.RegisterFlagCompletionFunc("components", completeWatchComponents)
}

type watchState struct {
	overallStatus string
	components    map[string]string
	incidents     int
	maintenances  int
}

func runWatch(cmd *cobra.Command, args []string) error {
	client, err := requireCloudStatusClient()
	if err != nil {
		return err
	}

	// Parse component filter
	var watchList []string
	if watchComponents != "" {
		for _, c := range strings.Split(watchComponents, ",") {
			watchList = append(watchList, strings.TrimSpace(c))
		}
	}

	// Set up signal handling for clean exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Initial state
	var lastState *watchState

	ticker := time.NewTicker(time.Duration(watchInterval) * time.Second)
	defer ticker.Stop()

	// Run immediately, then on interval
	printStatus := func() (bool, error) {
		currentState, err := getWatchState(client, watchList)
		if err != nil {
			return false, err
		}

		// Check for changes
		changed := stateChanged(lastState, currentState)

		if !watchNoClear && lastState != nil {
			// Clear screen (ANSI escape codes)
			fmt.Print("\033[H\033[2J")
		}

		printWatchDisplay(currentState, watchList)

		if changed && lastState != nil && watchExitOnChange {
			fmt.Println("\n[STATUS CHANGED - Exiting]")
			return true, nil
		}

		lastState = currentState
		return false, nil
	}

	// Initial display
	shouldExit, err := printStatus()
	if err != nil {
		return err
	}
	if shouldExit {
		return nil
	}

	fmt.Printf("\n[Watching... Interval: %ds, Press Ctrl+C to stop]\n", watchInterval)

	// Watch loop
	for {
		select {
		case <-ticker.C:
			shouldExit, err := printStatus()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}
			if shouldExit {
				return nil
			}
			if !watchNoClear {
				fmt.Printf("\n[Watching... Interval: %ds, Press Ctrl+C to stop]\n", watchInterval)
			}
		case <-sigChan:
			fmt.Println("\n[Stopping watch]")
			return nil
		}
	}
}

func getWatchState(client *cloudstatus.Client, watchList []string) (*watchState, error) {
	// Clear cache to get fresh data
	client.ClearCache()

	summary, err := client.GetSummary()
	if err != nil {
		return nil, err
	}

	state := &watchState{
		overallStatus: summary.Status.Indicator,
		components:    make(map[string]string),
	}

	// Count active incidents
	for _, inc := range summary.Incidents {
		if inc.IsActive() {
			state.incidents++
		}
	}

	// Count active maintenances
	for _, maint := range summary.ScheduledMaintenances {
		if maint.IsActive() {
			state.maintenances++
		}
	}

	// Track component states
	for _, comp := range summary.Components {
		if comp.Group {
			continue
		}

		// If watch list is specified, only track those components
		if len(watchList) > 0 {
			matched := false
			for _, watch := range watchList {
				if strings.Contains(strings.ToLower(comp.Name), strings.ToLower(watch)) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		state.components[comp.ID] = comp.Status
	}

	return state, nil
}

func stateChanged(old, new *watchState) bool {
	if old == nil {
		return false
	}

	if old.overallStatus != new.overallStatus {
		return true
	}

	if old.incidents != new.incidents {
		return true
	}

	if old.maintenances != new.maintenances {
		return true
	}

	// Check component changes
	for id, status := range new.components {
		if oldStatus, ok := old.components[id]; ok {
			if oldStatus != status {
				return true
			}
		}
	}

	return false
}

func printWatchDisplay(state *watchState, watchList []string) {
	now := time.Now().Format("2006-01-02 15:04:05 MST")
	fmt.Printf("=== F5 Cloud Status Monitor ===\n")
	fmt.Printf("Last Updated: %s\n\n", now)

	// Overall status with color indicator
	statusColor := getStatusColor(state.overallStatus)
	fmt.Printf("Overall Status: %s%s%s (%s)\n",
		statusColor,
		strings.ToUpper(state.overallStatus),
		"\033[0m", // Reset color
		cloudstatus.StatusIndicatorDescription(state.overallStatus))

	fmt.Printf("Active Incidents: %d\n", state.incidents)
	fmt.Printf("Active Maintenance: %d\n\n", state.maintenances)

	// Show watched components if any
	if len(state.components) > 0 && len(watchList) > 0 {
		fmt.Println("Watched Components:")
		// Get component details for display
		client := GetCloudStatusClient()
		resp, _ := client.GetSummary()
		if resp != nil {
			for _, comp := range resp.Components {
				if _, ok := state.components[comp.ID]; ok {
					color := getComponentStatusColor(comp.Status)
					fmt.Printf("  %s%s%s: %s\n", color, comp.Name, "\033[0m", comp.Status)
				}
			}
		}
	}

	// Show degraded components
	client := GetCloudStatusClient()
	resp, _ := client.GetSummary()
	if resp != nil {
		degraded := []string{}
		for _, comp := range resp.Components {
			if !comp.Group && comp.IsDegraded() {
				degraded = append(degraded, fmt.Sprintf("%s (%s)", comp.Name, comp.Status))
			}
		}
		if len(degraded) > 0 {
			fmt.Println("\nDegraded Components:")
			for _, d := range degraded {
				fmt.Printf("  ⚠ %s\n", d)
			}
		}
	}
}

func getStatusColor(status string) string {
	switch status {
	case cloudstatus.StatusNone:
		return "\033[32m" // Green
	case cloudstatus.StatusMinor:
		return "\033[33m" // Yellow
	case cloudstatus.StatusMajor:
		return "\033[38;5;208m" // Orange
	case cloudstatus.StatusCritical:
		return "\033[31m" // Red
	case cloudstatus.StatusMaintenance:
		return "\033[34m" // Blue
	default:
		return ""
	}
}

func getComponentStatusColor(status string) string {
	switch status {
	case cloudstatus.ComponentOperational:
		return "\033[32m" // Green
	case cloudstatus.ComponentDegradedPerformance:
		return "\033[33m" // Yellow
	case cloudstatus.ComponentPartialOutage:
		return "\033[38;5;208m" // Orange
	case cloudstatus.ComponentMajorOutage:
		return "\033[31m" // Red
	case cloudstatus.ComponentUnderMaintenance:
		return "\033[34m" // Blue
	default:
		return ""
	}
}
