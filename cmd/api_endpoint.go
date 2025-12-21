package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/spf13/cobra"

	"github.com/robinmordasiewicz/f5xcctl/pkg/output"
)

// API-endpoint flags matching original f5xcctl
var (
	apiEndpointAppType     string
	apiEndpointLogColor    bool
	apiEndpointLogFabulous bool
	apiEndpointLogLevel    int
	apiEndpointRange       string
)

// Discover-specific flags
var (
	discoverNamespace string
)

// Control-specific flags
var (
	controlNs  string
	discoverNs string
	deleteFlag bool
)

var apiEndpointCmd = &cobra.Command{
	Use:   "api-endpoint",
	Short: "Discover and manage API endpoints within F5 XC service mesh.",
	Long: `Discover and manage API endpoints within F5 XC service mesh.

This command group provides tools for API discovery and security policy
generation based on service mesh traffic analysis. F5 XC automatically
discovers API endpoints from observed traffic between services.

WORKFLOW:
  1. Use 'discover' to find API endpoints between services
  2. Review discovered endpoints and their communication patterns
  3. Use 'control' to generate L7 policies based on discoveries

COMMANDS:
  discover  Find API endpoints in service mesh traffic
  control   Generate L7 policies from discovered endpoints

AI assistants should run 'discover' first to understand the service mesh
topology before using 'control' to create security policies.`,
	Example: `  # Discover API endpoints in a namespace
  f5xcctl api-endpoint discover --namespace default --app-type example-app

  # Create L7 policies from discovered endpoints
  f5xcctl api-endpoint control --discover-ns default --app-type example-app

  # Check available commands
  f5xcctl api-endpoint --help`,
}

var apiEndpointDiscoverCmd = &cobra.Command{
	Use:     "discover [<flags>]",
	Aliases: []string{"discover"},
	Short:   "Discover API endpoints between services in a service mesh.",
	Long: `Discover API endpoints between services in a service mesh.

This command performs three steps:
1. Find all nodes and edges in the service mesh graph for an App Type
2. Find all API endpoints discovered between each edge
3. Display the information in tabular format`,
	Example: `f5xcctl api-endpoint discover --namespace default --app-type edge-checkoutcheckout`,
	RunE:    runAPIEndpointDiscover,
}

var apiEndpointControlCmd = &cobra.Command{
	Use:     "control [<flags>]",
	Aliases: []string{"ctrl"},
	Short:   "Create layer 7 policies based on discovered API endpoints.",
	Long: `Create layer 7 policies based on discovered API endpoints.

This command performs three steps:
1. Find all nodes and edges in the service mesh graph for an App Type
2. Find all API endpoints discovered between each edge
3. Create layer 7 policies that allow only known service-to-service communication

You can discover APIs from one namespace and apply policies in another namespace.`,
	Example: `f5xcctl api-endpoint control --discover-ns default --app-type edge-checkoutcheckout`,
	RunE:    runAPIEndpointControl,
}

func init() {
	rootCmd.AddCommand(apiEndpointCmd)

	// Enable AI-agent-friendly error handling for invalid subcommands
	apiEndpointCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("unknown command %q for %q\n\nUsage: f5xcctl api-endpoint <action> [flags]\n\nAvailable actions:\n  discover, control\n\nRun 'f5xcctl api-endpoint --help' for usage", args[0], cmd.CommandPath())
		}
		return cmd.Help()
	}
	apiEndpointCmd.SuggestionsMinimumDistance = 2

	// API-endpoint flags matching original f5xcctl
	apiEndpointCmd.PersistentFlags().StringVar(&apiEndpointAppType, "app-type", "", "App type name labeled on vK8s services or HTTP load balancer objects.")
	apiEndpointCmd.PersistentFlags().BoolVar(&apiEndpointLogColor, "log-color", true, "Enable colored log output.")
	apiEndpointCmd.PersistentFlags().BoolVar(&apiEndpointLogFabulous, "log-fabulous", true, "Enable enhanced log formatting.")
	apiEndpointCmd.PersistentFlags().IntVar(&apiEndpointLogLevel, "log-level", 3, "Set the logging verbosity level (1-5).")
	apiEndpointCmd.PersistentFlags().StringVarP(&apiEndpointRange, "range", "r", "1h", "Time range for querying service mesh data (e.g., '1h', '24h').")

	// Discover command flags
	apiEndpointDiscoverCmd.Flags().StringVar(&discoverNamespace, "namespace", "default", "Namespace containing the service mesh graph.")
	apiEndpointCmd.AddCommand(apiEndpointDiscoverCmd)

	// Control command flags
	apiEndpointControlCmd.Flags().StringVar(&discoverNs, "discover-ns", "default", "Namespace to discover API endpoints from.")
	apiEndpointControlCmd.Flags().StringVar(&controlNs, "control-ns", "", "Namespace where service policies will be applied.")
	apiEndpointControlCmd.Flags().BoolVarP(&deleteFlag, "delete", "d", false, "Delete all managed layer 7 policies.")
	apiEndpointCmd.AddCommand(apiEndpointControlCmd)
}

func runAPIEndpointDiscover(cmd *cobra.Command, args []string) error {
	apiClient := GetClient()
	if apiClient == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Build query parameters
	params := url.Values{}
	if apiEndpointAppType != "" {
		params.Set("app_type", apiEndpointAppType)
	}
	if apiEndpointRange != "" {
		params.Set("range", apiEndpointRange)
	}

	ns := discoverNamespace
	if ns == "" {
		ns = "default"
	}

	path := fmt.Sprintf("/api/config/namespaces/%s/api_endpoint_discovery", ns)
	resp, err := apiClient.Get(ctx, path, params)
	if err != nil {
		return fmt.Errorf("failed to discover API endpoints: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(resp.Body))
	}

	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return output.Print(result, GetOutputFormat())
}

func runAPIEndpointControl(cmd *cobra.Command, args []string) error {
	apiClient := GetClient()
	if apiClient == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Build query parameters
	params := url.Values{}
	if apiEndpointAppType != "" {
		params.Set("app_type", apiEndpointAppType)
	}
	if apiEndpointRange != "" {
		params.Set("range", apiEndpointRange)
	}

	ns := discoverNs
	if ns == "" {
		ns = "default"
	}

	// Note: controlNs flag exists but isn't currently used in the API path
	// It's reserved for future use when applying service policies to a different namespace
	_ = controlNs

	path := fmt.Sprintf("/api/config/namespaces/%s/api_endpoint_control", ns)
	resp, err := apiClient.Get(ctx, path, params)
	if err != nil {
		return fmt.Errorf("failed to get API endpoint control: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(resp.Body))
	}

	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return output.Print(result, GetOutputFormat())
}
