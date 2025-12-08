package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/spf13/cobra"

	"github.com/robinmordasiewicz/vesctl/pkg/output"
)

// API-endpoint flags matching original vesctl
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
	controlNs   string
	discoverNs  string
	deleteFlag  bool
)

var apiEndpointCmd = &cobra.Command{
	Use:     "api-endpoint",
	Short:   "Execute commands for API Endpoint Discovery and Control",
	Long:    `Execute commands for API Endpoint Discovery and Control`,
	Example: `vesctl api-endpoint discover --namespace default`,
}

var apiEndpointDiscoverCmd = &cobra.Command{
	Use:     "discover [<flags>]",
	Aliases: []string{"discover"},
	Short:   "Discover all API endpoints for each of the edge (client --> server), belonging to an App Type",
	Long: `It will run three steps process
1. Find all the nodes and edges in a service mesh graph of an App Type
2. Find all the api endpoints discovered between an edge
3. Display the information in tabular format`,
	Example: `vesctl api-endpoint discover --namespace default --app-type edge-checkoutcheckout`,
	RunE:    runAPIEndpointDiscover,
}

var apiEndpointControlCmd = &cobra.Command{
	Use:     "control [<flags>]",
	Aliases: []string{"ctrl"},
	Short:   "Find all the edges (client --> server) and create respective layer7 policies",
	Long: `It will run three steps process
1. Find all the nodes and edges in a service mesh graph of an App Type
2. Find all the api endpoints discovered between an edge
3. Create layer7 policies which will only allow communication between known edges(services) and
which are communicating with known api endpoint (method, url combination)

This command allows to discover the api's from one namespace and apply the policy in another namespace`,
	Example: `vesctl api-endpoint control --discover-ns default --app-type edge-checkoutcheckout`,
	RunE:    runAPIEndpointControl,
}

func init() {
	rootCmd.AddCommand(apiEndpointCmd)

	// API-endpoint flags matching original vesctl
	apiEndpointCmd.PersistentFlags().StringVar(&apiEndpointAppType, "app-type", "", "App Type name labelled on vk8s service or http lb objects, defaults to the value of namespace/discover-ns")
	apiEndpointCmd.PersistentFlags().BoolVar(&apiEndpointLogColor, "log-color", true, "enable color for your logs")
	apiEndpointCmd.PersistentFlags().BoolVar(&apiEndpointLogFabulous, "log-fabulous", true, "enable fabulous writer for your logs")
	apiEndpointCmd.PersistentFlags().IntVar(&apiEndpointLogLevel, "log-level", 3, "Log Level for Site Deployment")
	apiEndpointCmd.PersistentFlags().StringVarP(&apiEndpointRange, "range", "r", "1h", "range for which the graph and edges will be queried")

	// Discover command flags
	apiEndpointDiscoverCmd.Flags().StringVar(&discoverNamespace, "namespace", "default", "namespace where the service-mesh graph exists")
	apiEndpointCmd.AddCommand(apiEndpointDiscoverCmd)

	// Control command flags
	apiEndpointControlCmd.Flags().StringVar(&discoverNs, "discover-ns", "default", "Namespace where the service-mesh graph exists")
	apiEndpointControlCmd.Flags().StringVar(&controlNs, "control-ns", "", "Namespace on which the service policy will be applied, defaults to the value of discover-ns")
	apiEndpointControlCmd.Flags().BoolVarP(&deleteFlag, "delete", "d", false, "delete all managed policies")
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
