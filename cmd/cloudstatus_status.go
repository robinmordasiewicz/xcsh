package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/robinmordasiewicz/f5xcctl/pkg/cloudstatus"
)

var statusQuiet bool

var cloudstatusStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get overall F5 Cloud status indicator.",
	Long: `Get the overall F5 Distributed Cloud status indicator.

Returns a quick snapshot of the overall system health. Use --quiet mode
in CI/CD pipelines to get an exit code without output:

Exit Codes:
  0 - All systems operational (none)
  1 - Minor system issue (minor)
  2 - Major system issue (major)
  3 - Critical system outage (critical)
  4 - System under maintenance (maintenance)`,
	Example: `  # Quick status check
  f5xcctl cloudstatus status

  # CI/CD health gate (exits with appropriate code)
  f5xcctl cloudstatus status --quiet
  if [ $? -ne 0 ]; then echo "F5 XC has issues"; fi

  # JSON output for parsing
  f5xcctl cloudstatus status --output-format json`,
	RunE: runCloudstatusStatus,
}

func init() {
	cloudstatusCmd.AddCommand(cloudstatusStatusCmd)

	cloudstatusStatusCmd.Flags().BoolVarP(&statusQuiet, "quiet", "q", false, "Suppress output, return exit code only.")
}

func runCloudstatusStatus(cmd *cobra.Command, args []string) error {
	client := GetCloudStatusClient()
	if client == nil {
		return fmt.Errorf("cloudstatus client not initialized")
	}

	resp, err := client.GetStatus()
	if err != nil {
		if statusQuiet {
			os.Exit(cloudstatus.ExitCodeAPIError)
		}
		return fmt.Errorf("failed to get status: %w", err)
	}

	// In quiet mode, just exit with appropriate code
	if statusQuiet {
		exitCode := cloudstatus.StatusIndicatorToExitCode(resp.Status.Indicator)
		os.Exit(exitCode)
	}

	// Output based on format
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
		// Table format
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		_, _ = fmt.Fprintln(w, "STATUS\tDESCRIPTION")
		_, _ = fmt.Fprintf(w, "%s\t%s\n", resp.Status.Indicator, resp.Status.Description)
		return w.Flush()
	}
}
