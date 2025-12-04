package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/robinmordasiewicz/f5xc/pkg/output"
)

var discoverFlags struct {
	namespace     string
	srcService    string
	dstService    string
	virtualHost   string
	timeRange     string
	includeSystem bool
	format        string
}

var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover API endpoints in service mesh",
	Long: `Discover API endpoints between services in your service mesh.

This command queries the F5 Distributed Cloud API to discover:
  - API endpoints that have been observed between services
  - HTTP methods and paths being used
  - Traffic patterns and request counts

The results can be used to understand your API surface and
create appropriate security policies.`,
	Example: `  # Discover all endpoints in a namespace
  f5xc api-endpoint discover --namespace my-ns

  # Discover endpoints between specific services
  f5xc api-endpoint discover --namespace my-ns --src-service frontend --dst-service backend

  # Discover endpoints for a virtual host
  f5xc api-endpoint discover --virtual-host my-vh --namespace my-ns

  # Get results in table format
  f5xc api-endpoint discover --namespace my-ns --format table`,
	RunE: runDiscover,
}

func init() {
	apiEndpointCmd.AddCommand(discoverCmd)

	discoverCmd.Flags().StringVarP(&discoverFlags.namespace, "namespace", "n", "", "Namespace to discover endpoints in (required)")
	discoverCmd.Flags().StringVar(&discoverFlags.srcService, "src-service", "", "Filter by source service name")
	discoverCmd.Flags().StringVar(&discoverFlags.dstService, "dst-service", "", "Filter by destination service name")
	discoverCmd.Flags().StringVar(&discoverFlags.virtualHost, "virtual-host", "", "Filter by virtual host")
	discoverCmd.Flags().StringVar(&discoverFlags.timeRange, "time-range", "24h", "Time range for discovery (e.g., 1h, 24h, 7d)")
	discoverCmd.Flags().BoolVar(&discoverFlags.includeSystem, "include-system", false, "Include system/internal endpoints")
	discoverCmd.Flags().StringVar(&discoverFlags.format, "format", "", "Output format: table, json, yaml")
	discoverCmd.MarkFlagRequired("namespace")
}

// DiscoveredEndpoint represents a discovered API endpoint
type DiscoveredEndpoint struct {
	Method       string `json:"method" yaml:"method"`
	Path         string `json:"path" yaml:"path"`
	SrcService   string `json:"src_service" yaml:"src_service"`
	DstService   string `json:"dst_service" yaml:"dst_service"`
	VirtualHost  string `json:"virtual_host,omitempty" yaml:"virtual_host,omitempty"`
	RequestCount int64  `json:"request_count" yaml:"request_count"`
	AvgLatencyMs float64 `json:"avg_latency_ms" yaml:"avg_latency_ms"`
	ErrorRate    float64 `json:"error_rate" yaml:"error_rate"`
	LastSeen     string `json:"last_seen" yaml:"last_seen"`
}

func runDiscover(cmd *cobra.Command, args []string) error {
	apiClient := GetClient()
	if apiClient == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Build the discovery request
	requestBody := map[string]interface{}{
		"namespace": discoverFlags.namespace,
	}

	if discoverFlags.srcService != "" {
		requestBody["src_service"] = discoverFlags.srcService
	}
	if discoverFlags.dstService != "" {
		requestBody["dst_service"] = discoverFlags.dstService
	}
	if discoverFlags.virtualHost != "" {
		requestBody["virtual_host"] = discoverFlags.virtualHost
	}

	// Parse time range
	timeRange := parseTimeRange(discoverFlags.timeRange)
	requestBody["start_time"] = time.Now().Add(-timeRange).Format(time.RFC3339)
	requestBody["end_time"] = time.Now().Format(time.RFC3339)

	// Query the API discovery endpoint
	path := fmt.Sprintf("/api/web/namespaces/%s/api_endpoints/discover", discoverFlags.namespace)
	resp, err := apiClient.Post(ctx, path, requestBody)
	if err != nil {
		return fmt.Errorf("failed to discover endpoints: %w", err)
	}

	if resp.StatusCode >= 400 {
		// Try alternative endpoint paths
		path = fmt.Sprintf("/api/data/namespaces/%s/api_endpoints", discoverFlags.namespace)
		resp, err = apiClient.Get(ctx, path, nil)
		if err != nil {
			return fmt.Errorf("failed to discover endpoints: %w", err)
		}
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(resp.Body))
	}

	// Parse the response
	var rawResponse map[string]interface{}
	if err := json.Unmarshal(resp.Body, &rawResponse); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract and format discovered endpoints
	endpoints := extractDiscoveredEndpoints(rawResponse)

	// Filter endpoints if needed
	if !discoverFlags.includeSystem {
		endpoints = filterSystemEndpoints(endpoints)
	}

	// Output results
	outputFormat := discoverFlags.format
	if outputFormat == "" {
		outputFormat = GetOutputFormat()
	}

	if outputFormat == "table" {
		return printEndpointsTable(endpoints)
	}

	return output.Print(endpoints, outputFormat)
}

func parseTimeRange(rangeStr string) time.Duration {
	rangeStr = strings.ToLower(rangeStr)

	if strings.HasSuffix(rangeStr, "d") {
		days := 1
		fmt.Sscanf(rangeStr, "%dd", &days)
		return time.Duration(days) * 24 * time.Hour
	}

	if strings.HasSuffix(rangeStr, "h") {
		hours := 24
		fmt.Sscanf(rangeStr, "%dh", &hours)
		return time.Duration(hours) * time.Hour
	}

	if strings.HasSuffix(rangeStr, "m") {
		minutes := 60
		fmt.Sscanf(rangeStr, "%dm", &minutes)
		return time.Duration(minutes) * time.Minute
	}

	// Default to 24 hours
	return 24 * time.Hour
}

func extractDiscoveredEndpoints(response map[string]interface{}) []DiscoveredEndpoint {
	endpoints := []DiscoveredEndpoint{}

	// Try to extract from items array
	if items, ok := response["items"].([]interface{}); ok {
		for _, item := range items {
			if endpoint := parseEndpointItem(item); endpoint != nil {
				endpoints = append(endpoints, *endpoint)
			}
		}
	}

	// Try to extract from api_endpoints array
	if apiEndpoints, ok := response["api_endpoints"].([]interface{}); ok {
		for _, item := range apiEndpoints {
			if endpoint := parseEndpointItem(item); endpoint != nil {
				endpoints = append(endpoints, *endpoint)
			}
		}
	}

	// Try to extract from edges in a graph response
	if edges, ok := response["edges"].([]interface{}); ok {
		for _, edge := range edges {
			if edgeMap, ok := edge.(map[string]interface{}); ok {
				if endpoints, ok := edgeMap["endpoints"].([]interface{}); ok {
					for _, ep := range endpoints {
						if endpoint := parseEndpointItem(ep); endpoint != nil {
							endpoints = append(endpoints, *endpoint)
						}
					}
				}
			}
		}
	}

	return endpoints
}

func parseEndpointItem(item interface{}) *DiscoveredEndpoint {
	itemMap, ok := item.(map[string]interface{})
	if !ok {
		return nil
	}

	endpoint := &DiscoveredEndpoint{}

	if method, ok := itemMap["method"].(string); ok {
		endpoint.Method = method
	} else if method, ok := itemMap["http_method"].(string); ok {
		endpoint.Method = method
	}

	if path, ok := itemMap["path"].(string); ok {
		endpoint.Path = path
	} else if path, ok := itemMap["api_path"].(string); ok {
		endpoint.Path = path
	} else if path, ok := itemMap["uri"].(string); ok {
		endpoint.Path = path
	}

	if src, ok := itemMap["src_service"].(string); ok {
		endpoint.SrcService = src
	} else if src, ok := itemMap["source_service"].(string); ok {
		endpoint.SrcService = src
	}

	if dst, ok := itemMap["dst_service"].(string); ok {
		endpoint.DstService = dst
	} else if dst, ok := itemMap["destination_service"].(string); ok {
		endpoint.DstService = dst
	} else if dst, ok := itemMap["service"].(string); ok {
		endpoint.DstService = dst
	}

	if vh, ok := itemMap["virtual_host"].(string); ok {
		endpoint.VirtualHost = vh
	}

	if count, ok := itemMap["request_count"].(float64); ok {
		endpoint.RequestCount = int64(count)
	} else if count, ok := itemMap["count"].(float64); ok {
		endpoint.RequestCount = int64(count)
	}

	if latency, ok := itemMap["avg_latency_ms"].(float64); ok {
		endpoint.AvgLatencyMs = latency
	} else if latency, ok := itemMap["latency"].(float64); ok {
		endpoint.AvgLatencyMs = latency
	}

	if errRate, ok := itemMap["error_rate"].(float64); ok {
		endpoint.ErrorRate = errRate
	}

	if lastSeen, ok := itemMap["last_seen"].(string); ok {
		endpoint.LastSeen = lastSeen
	} else if lastSeen, ok := itemMap["timestamp"].(string); ok {
		endpoint.LastSeen = lastSeen
	}

	// Only return if we have at least method and path
	if endpoint.Method == "" && endpoint.Path == "" {
		return nil
	}

	return endpoint
}

func filterSystemEndpoints(endpoints []DiscoveredEndpoint) []DiscoveredEndpoint {
	filtered := []DiscoveredEndpoint{}
	systemPaths := []string{"/healthz", "/readyz", "/livez", "/metrics", "/_internal", "/favicon.ico"}

	for _, ep := range endpoints {
		isSystem := false
		for _, syspath := range systemPaths {
			if strings.HasPrefix(ep.Path, syspath) {
				isSystem = true
				break
			}
		}
		if !isSystem {
			filtered = append(filtered, ep)
		}
	}

	return filtered
}

func printEndpointsTable(endpoints []DiscoveredEndpoint) error {
	if len(endpoints) == 0 {
		fmt.Println("No endpoints discovered")
		return nil
	}

	// Print header
	fmt.Printf("%-8s %-40s %-20s %-20s %-12s\n", "METHOD", "PATH", "SRC SERVICE", "DST SERVICE", "REQUESTS")
	fmt.Println(strings.Repeat("-", 104))

	// Print rows
	for _, ep := range endpoints {
		path := ep.Path
		if len(path) > 40 {
			path = path[:37] + "..."
		}
		srcSvc := ep.SrcService
		if len(srcSvc) > 20 {
			srcSvc = srcSvc[:17] + "..."
		}
		dstSvc := ep.DstService
		if len(dstSvc) > 20 {
			dstSvc = dstSvc[:17] + "..."
		}

		fmt.Printf("%-8s %-40s %-20s %-20s %-12d\n",
			ep.Method,
			path,
			srcSvc,
			dstSvc,
			ep.RequestCount,
		)
	}

	fmt.Printf("\nTotal: %d endpoints discovered\n", len(endpoints))
	return nil
}
