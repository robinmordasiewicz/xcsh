package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/robinmordasiewicz/xcsh/pkg/cloudstatus"
)

// cloudstatus-specific flags
var (
	cloudstatusNoCache  bool          // Bypass cache
	cloudstatusCacheTTL time.Duration // Custom cache TTL
)

// cloudstatusClient is the shared client for all cloudstatus subcommands
var cloudstatusClient *cloudstatus.Client

var cloudstatusCmd = &cobra.Command{
	Use:     "cloudstatus",
	Aliases: []string{"cs", "status"},
	Short:   "Monitor F5 Distributed Cloud service status, incidents, and maintenance.",
	Long: `Monitor F5 Distributed Cloud service status, incidents, and maintenance.

This command group provides access to the F5 Cloud Status API to check service
health, view active incidents, and track scheduled maintenance windows.

No authentication is required - the status API is publicly accessible.

Use 'f5xcctl cloudstatus status' for a quick check, or 'f5xcctl cloudstatus summary'
for a comprehensive overview of all services.`,
	Example: `  # Quick overall status check
  f5xcctl cloudstatus status

  # Quick status check with exit code for scripting
  f5xcctl cloudstatus status --quiet

  # Complete status summary
  f5xcctl cloudstatus summary

  # List all components
  f5xcctl cloudstatus components list

  # Show only degraded components
  f5xcctl cloudstatus components list --degraded-only

  # List active incidents
  f5xcctl cloudstatus incidents active

  # List upcoming maintenance windows
  f5xcctl cloudstatus maintenance upcoming

  # Check PoP status by region
  f5xcctl cloudstatus pops status --region north-america

  # Real-time monitoring
  f5xcctl cloudstatus watch --interval 30

  # Output machine-readable spec for AI agents
  f5xcctl cloudstatus --spec --output-format json`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Handle --spec flag for cloudstatus command
		if CheckSpecFlag() {
			spec := cloudstatus.GenerateSpec()
			format := GetOutputFormatWithDefault("json")
			return OutputCloudStatusSpec(spec, format)
		}

		// Initialize the cloudstatus client with caching options
		opts := []cloudstatus.ClientOption{}
		if cloudstatusNoCache {
			opts = append(opts, cloudstatus.WithoutCache())
		} else if cloudstatusCacheTTL > 0 {
			opts = append(opts, cloudstatus.WithCache(cloudstatusCacheTTL))
		}
		cloudstatusClient = cloudstatus.NewClient(opts...)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cloudstatusCmd)

	// Enable AI-agent-friendly error handling for invalid subcommands
	cloudstatusCmd.RunE = func(cmd *cobra.Command, args []string) error {
		// Handle --spec flag
		if CheckSpecFlag() {
			spec := cloudstatus.GenerateSpec()
			format := GetOutputFormatWithDefault("json")
			return OutputCloudStatusSpec(spec, format)
		}

		if len(args) > 0 {
			return fmt.Errorf("unknown command %q for %q\n\nUsage: f5xcctl cloudstatus <command> [flags]\n\nAvailable Commands:\n  status, summary, components, incidents, maintenance, pops, watch\n\nRun 'f5xcctl cloudstatus --help' for usage", args[0], cmd.CommandPath())
		}
		return cmd.Help()
	}
	cloudstatusCmd.SuggestionsMinimumDistance = 2

	// Cloudstatus-specific flags
	cloudstatusCmd.PersistentFlags().BoolVar(&cloudstatusNoCache, "no-cache", false, "Bypass the response cache for real-time data.")
	cloudstatusCmd.PersistentFlags().DurationVar(&cloudstatusCacheTTL, "cache-ttl", 60*time.Second, "Override the cache time-to-live duration.")

	// Register --spec flag for machine-readable CLI specification
	RegisterSpecFlag(cloudstatusCmd)
}

// GetCloudStatusClient returns the initialized cloudstatus client
func GetCloudStatusClient() *cloudstatus.Client {
	return cloudstatusClient
}

// OutputCloudStatusSpec outputs the cloudstatus specification in the requested format
func OutputCloudStatusSpec(spec *cloudstatus.Spec, format string) error {
	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(spec); err != nil {
			return fmt.Errorf("failed to encode JSON: %w", err)
		}
		os.Exit(0)
		return nil
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		if err := encoder.Encode(spec); err != nil {
			return fmt.Errorf("failed to encode YAML: %w", err)
		}
		os.Exit(0)
		return nil
	default:
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(spec); err != nil {
			return fmt.Errorf("failed to encode JSON: %w", err)
		}
		os.Exit(0)
		return nil
	}
}
