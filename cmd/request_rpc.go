package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/robinmordasiewicz/vesctl/pkg/client"
	"github.com/robinmordasiewicz/vesctl/pkg/output"
)

var rpcFlags struct {
	inputFile  string
	namespace  string
	httpMethod string
	uri        string
}

var rpcCmd = &cobra.Command{
	Use:     "rpc",
	Short:   "RPC Invocation",
	Long:    `RPC Invocation`,
	Example: `vesctl request rpc registration.CustomAPI.RegistrationApprove -i approval_req.yaml`,
	Args:    cobra.MaximumNArgs(1),
	RunE:    runRPC,
}

func init() {
	requestCmd.AddCommand(rpcCmd)

	rpcCmd.Flags().StringVarP(&rpcFlags.inputFile, "input-file", "i", "", "File containing request data (YAML or JSON)")
	rpcCmd.Flags().StringVarP(&rpcFlags.namespace, "namespace", "n", "", "Namespace for the API call")
	rpcCmd.Flags().StringVar(&rpcFlags.httpMethod, "http-method", "POST", "HTTP method (GET, POST, PUT, DELETE)")
	rpcCmd.Flags().StringVar(&rpcFlags.uri, "uri", "", "URI path for the API call")
}

func runRPC(cmd *cobra.Command, args []string) error {
	apiClient := GetClient()
	if apiClient == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	// Determine the URI
	var uri string
	if rpcFlags.uri != "" {
		uri = rpcFlags.uri
	} else if len(args) > 0 {
		// Convert dotted notation to URI
		uri = convertEndpointToURI(args[0], rpcFlags.namespace)
	} else {
		return fmt.Errorf("either an endpoint argument or --uri flag is required")
	}

	// Load request body if input file is provided
	var requestBody map[string]interface{}
	if rpcFlags.inputFile != "" {
		data, err := os.ReadFile(rpcFlags.inputFile)
		if err != nil {
			return fmt.Errorf("failed to read input file: %w", err)
		}

		// Try YAML first, then JSON
		if err := yaml.Unmarshal(data, &requestBody); err != nil {
			if err := json.Unmarshal(data, &requestBody); err != nil {
				return fmt.Errorf("failed to parse input file (not valid YAML or JSON): %w", err)
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var resp *client.Response
	var err error

	// Execute the request based on HTTP method
	switch strings.ToUpper(rpcFlags.httpMethod) {
	case "GET":
		resp, err = apiClient.Get(ctx, uri, nil)
	case "POST":
		resp, err = apiClient.Post(ctx, uri, requestBody)
	case "PUT":
		resp, err = apiClient.Put(ctx, uri, requestBody)
	case "DELETE":
		resp, err = apiClient.Delete(ctx, uri)
	default:
		return fmt.Errorf("unsupported HTTP method: %s", rpcFlags.httpMethod)
	}

	if err != nil {
		return fmt.Errorf("RPC call failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(resp.Body))
	}

	// Parse and output the response
	var result interface{}
	if len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, &result); err != nil {
			// If not JSON, output as string
			fmt.Println(string(resp.Body))
			return nil
		}
	}

	return output.Print(result, GetOutputFormat())
}

// convertEndpointToURI converts a dotted endpoint notation to a URI path
// e.g., "namespace.CustomAPI.List" -> "/api/web/custom/namespaces/{namespace}/list"
func convertEndpointToURI(endpoint string, namespace string) string {
	parts := strings.Split(endpoint, ".")
	if len(parts) < 2 {
		// If no dots, assume it's already a path
		return endpoint
	}

	// Build URI based on common patterns
	resourceType := strings.ToLower(parts[0])
	apiType := parts[1]
	var operation string
	if len(parts) > 2 {
		operation = strings.ToLower(parts[2])
	}

	// Handle CustomAPI pattern
	if apiType == "CustomAPI" || apiType == "CustomDataAPI" || apiType == "CustomPublicAPI" {
		var path string
		if namespace != "" {
			path = fmt.Sprintf("/api/web/custom/namespaces/%s/%s", namespace, resourceType)
		} else {
			path = fmt.Sprintf("/api/web/custom/%s", resourceType)
		}
		if operation != "" {
			path = fmt.Sprintf("%s/%s", path, operation)
		}
		return path
	}

	// Default: use the endpoint as-is with namespace substitution
	if namespace != "" {
		return fmt.Sprintf("/api/config/namespaces/%s/%ss", namespace, resourceType)
	}
	return fmt.Sprintf("/api/%s", strings.ReplaceAll(endpoint, ".", "/"))
}
