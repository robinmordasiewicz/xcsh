//go:build ignore
// +build ignore

// This file generates RPC subcommands from OpenAPI specifications.
// Run with: go run cmd/gen_rpc/main.go

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

// SwaggerFile represents the structure of a Swagger/OpenAPI file
type SwaggerFile struct {
	Paths map[string]json.RawMessage `json:"paths"`
}

// Operation represents an HTTP operation
type Operation struct {
	OperationID  string          `json:"operationId"`
	XVesProtoRPC string          `json:"x-ves-proto-rpc"`
	RequestBody  json.RawMessage `json:"requestBody"`
}

// RequestBody represents the request body structure
type RequestBody struct {
	Content map[string]ContentType `json:"content"`
}

// ContentType represents the content type structure
type ContentType struct {
	Schema SchemaRef `json:"schema"`
}

// SchemaRef represents a schema reference
type SchemaRef struct {
	Ref string `json:"$ref"`
}

// RPCMethod represents an RPC method for code generation
type RPCMethod struct {
	FullName   string // e.g., ves.io.schema.alert.CustomAPI.Alerts
	ShortName  string // e.g., alert.CustomAPI.Alerts
	SchemaType string // e.g., ves.io.schema.alert.Request
}

// apiTypePattern matches any API type like CustomAPI, NamespaceCustomAPI, CustomDataK8SAPI, SignatureCustomApi, etc.
// Uses case-insensitive matching for API/Api suffix
var apiTypePattern = regexp.MustCompile(`\.([A-Za-z0-9]+(?:API|Api))\.`)

func main() {
	specDir := "docs/specifications/api"
	outputFile := "cmd/request_rpc_generated.go"

	methods, err := extractRPCMethods(specDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting RPC methods: %v\n", err)
		os.Exit(1)
	}

	if err := generateCode(methods, outputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating code: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %d RPC subcommands in %s\n", len(methods), outputFile)
}

func extractRPCMethods(specDir string) ([]RPCMethod, error) {
	files, err := filepath.Glob(filepath.Join(specDir, "*.json"))
	if err != nil {
		return nil, err
	}

	methodMap := make(map[string]RPCMethod)

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		var swagger SwaggerFile
		if err := json.Unmarshal(data, &swagger); err != nil {
			continue
		}

		for _, pathItemRaw := range swagger.Paths {
			// Parse the path item as a map
			var pathItem map[string]json.RawMessage
			if err := json.Unmarshal(pathItemRaw, &pathItem); err != nil {
				continue
			}

			// Check each HTTP method
			httpMethods := []string{"get", "post", "put", "delete", "patch"}
			for _, httpMethod := range httpMethods {
				opRaw, exists := pathItem[httpMethod]
				if !exists {
					continue
				}

				var op Operation
				if err := json.Unmarshal(opRaw, &op); err != nil {
					continue
				}

				if op.XVesProtoRPC == "" {
					continue
				}

				// Skip standard CRUD API operations (they're handled by resource commands)
				// Only include custom API types (anything with "Custom" or specific API suffixes)
				if !isCustomAPIMethod(op.XVesProtoRPC) {
					continue
				}

				shortName := convertToShortName(op.XVesProtoRPC)
				if shortName == "" {
					continue
				}

				// Extract schema type from request body or fall back to default
				schemaType := extractSchemaType(op.XVesProtoRPC, op.RequestBody)

				method := RPCMethod{
					FullName:   op.XVesProtoRPC,
					ShortName:  shortName,
					SchemaType: schemaType,
				}

				// Use shortName as key to avoid duplicates
				if _, exists := methodMap[shortName]; !exists {
					methodMap[shortName] = method
				}
			}
		}
	}

	// Convert map to sorted slice
	var methods []RPCMethod
	for _, m := range methodMap {
		methods = append(methods, m)
	}
	sort.Slice(methods, func(i, j int) bool {
		return methods[i].ShortName < methods[j].ShortName
	})

	return methods, nil
}

// isCustomAPIMethod returns true if the RPC method is a custom API (not standard CRUD)
func isCustomAPIMethod(fullName string) bool {
	// Standard CRUD APIs that should be excluded (handled by resource commands)
	// These are like ves.io.schema.namespace.API.List, ves.io.schema.namespace.API.Create
	if strings.Contains(fullName, ".API.") && !strings.Contains(fullName, "CustomAPI") {
		// Check if it's a standard CRUD operation
		crudOps := []string{".API.List", ".API.Get", ".API.Create", ".API.Replace", ".API.Delete"}
		for _, op := range crudOps {
			if strings.HasSuffix(fullName, op) {
				return false
			}
		}
	}

	// Include any API type that contains "Custom" or specific monitoring/data APIs
	customPatterns := []string{
		"CustomAPI", "CustomDataAPI", "CustomPublicAPI",
		"CustomDataK8SAPI", "CustomAlertAPI", "CustomEventAPI",
		"CustomEventDetailsAPI", "CustomGeneralAPI", "CustomMitigationAPI",
		"CustomNetworkAPI", "CustomReportAPI", "CustomFlowConnectionAPI",
		"CustomSiteStatusAPI", "CustomStateAPI",
		"MonitoringAPI", "SuggestionAPI", "ClientRuleAPI",
		"WafExclusionAPI", "ThreatCampaignAPI", "FlowConnectionAPI",
		"ConfigKubeConfigAPI", "UamKubeConfigAPI", "KubeConfigAPI",
		"K8SAPI", "StatusAPI", "StateAPI",
		"UpgradeAPI", "ActionAPI", "UsageAPI",
		"SignatureCustomApi", "WafSignatureChangelogCustomApi",
		"ApiepLBCustomAPI", "ApiepCustomAPI",
	}

	for _, pattern := range customPatterns {
		if strings.Contains(fullName, pattern) {
			return true
		}
	}

	return false
}

// convertToShortName converts ves.io.schema.alert.CustomAPI.Alerts to alert.CustomAPI.Alerts
// or ves.io.schema.namespace.NamespaceCustomAPI.CascadeDelete to namespace.NamespaceCustomAPI.CascadeDelete
func convertToShortName(fullName string) string {
	// Remove ves.io.schema. prefix
	name := strings.TrimPrefix(fullName, "ves.io.schema.")

	// Find the API type using regex (matches any *API pattern)
	matches := apiTypePattern.FindStringIndex(name)
	if matches == nil {
		return ""
	}

	// Extract the API type
	apiMatch := apiTypePattern.FindStringSubmatch(name)
	if len(apiMatch) < 2 {
		return ""
	}
	apiType := apiMatch[1]

	// Find the position of the API type in the string
	apiTypeWithDots := "." + apiType + "."
	apiIdx := strings.Index(name, apiTypeWithDots)
	if apiIdx == -1 {
		return ""
	}

	// Get the resource type (part before the API type)
	resourcePart := name[:apiIdx]

	// For nested resources like "pbac.addon_service", take the last part (addon_service)
	parts := strings.Split(resourcePart, ".")
	resourceType := parts[len(parts)-1]

	// Get the method name (part after API type)
	methodPart := name[apiIdx+len(apiTypeWithDots):]

	return resourceType + "." + apiType + "." + methodPart
}

// extractSchemaType extracts the request schema type from the request body or uses default
func extractSchemaType(fullName string, requestBodyRaw json.RawMessage) string {
	// Try to extract from request body
	if len(requestBodyRaw) > 0 {
		var requestBody RequestBody
		if err := json.Unmarshal(requestBodyRaw, &requestBody); err == nil {
			if content, ok := requestBody.Content["application/json"]; ok {
				if content.Schema.Ref != "" {
					// Extract schema name from $ref like "#/components/schemas/cdn_loadbalancerLilacCDNMetricsRequest"
					ref := content.Schema.Ref
					if idx := strings.LastIndex(ref, "/"); idx != -1 {
						schemaName := ref[idx+1:]
						// Convert schema name to full path
						// e.g., cdn_loadbalancerLilacCDNMetricsRequest -> ves.io.schema.views.cdn_loadbalancer.LilacCDNMetricsRequest
						return convertSchemaNameToFullPath(fullName, schemaName)
					}
				}
			}
		}
	}

	// Fall back to default pattern
	return extractDefaultSchemaType(fullName)
}

// convertSchemaNameToFullPath converts a schema name to a full qualified path
func convertSchemaNameToFullPath(rpcFullName, schemaName string) string {
	// Extract the base package from the RPC name using the API type pattern
	// ves.io.schema.views.cdn_loadbalancer.CustomAPI.CDNMetrics -> ves.io.schema.views.cdn_loadbalancer
	matches := apiTypePattern.FindStringIndex(rpcFullName)
	basePkg := rpcFullName
	if matches != nil {
		basePkg = rpcFullName[:matches[0]]
	}

	// The schema name typically has the package prefix as underscore-separated
	// e.g., cdn_loadbalancerLilacCDNMetricsRequest
	// We need to find where the package ends and the type name starts

	// Get the last part of the package (e.g., "cdn_loadbalancer" from "ves.io.schema.views.cdn_loadbalancer")
	pkgParts := strings.Split(basePkg, ".")
	lastPkgPart := pkgParts[len(pkgParts)-1]

	// Try to find and remove the package prefix from schema name
	if strings.HasPrefix(schemaName, lastPkgPart) {
		typeName := schemaName[len(lastPkgPart):]
		return basePkg + "." + typeName
	}

	// If no match, just use the full schema name with the base package
	return basePkg + "." + schemaName
}

// extractDefaultSchemaType extracts a default schema type from the RPC name
func extractDefaultSchemaType(fullName string) string {
	// Use regex to find any API type pattern
	matches := apiTypePattern.FindStringIndex(fullName)
	if matches != nil {
		return fullName[:matches[0]] + ".Request"
	}

	return fullName + ".Request"
}

var codeTemplate = `// Code generated by cmd/gen_rpc/main.go. DO NOT EDIT.

package cmd

import (
	"github.com/spf13/cobra"
)

// rpcSubcommands is a list of all RPC methods available
var rpcSubcommands = []struct {
	Name       string
	SchemaType string
}{
{{- range .Methods }}
	{"{{ .ShortName }}", "{{ .SchemaType }}"},
{{- end }}
}

func init() {
	// Register all RPC subcommands
	for _, sub := range rpcSubcommands {
		cmd := createRPCSubcommand(sub.Name, sub.SchemaType)
		rpcCmd.AddCommand(cmd)
	}
}

// createRPCSubcommand creates a Cobra command for an RPC method
func createRPCSubcommand(name, schemaType string) *cobra.Command {
	var inputFile string
	var httpMethod string
	var jsonData string
	var uri string

	cmd := &cobra.Command{
		Use:   name,
		Short: "CustomAPI RPC invocation",
		Long:  "CustomAPI RPC invocation",
		Example: "xcsh request rpc registration.CustomAPI.RegistrationApprove -i approval_req.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the RPC flags from local flags
			rpcFlags.inputFile = inputFile
			rpcFlags.httpMethod = httpMethod
			rpcFlags.uri = uri

			// If json-data is provided, we need to handle it differently
			// For now, just call runRPC with the method name as argument
			return runRPC(cmd, []string{name})
		},
	}

	cmd.Flags().StringVar(&httpMethod, "http-method", "", "HTTP Method(POST/GET/DELETE) to use")
	cmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "File containing request message in yaml")
	cmd.Flags().StringVar(&jsonData, "json-data", "", "Inline "+schemaType+" contents in json form")
	cmd.Flags().StringVar(&uri, "uri", "", "URI to use for custom API")

	return cmd
}
`

func generateCode(methods []RPCMethod, outputFile string) error {
	tmpl, err := template.New("code").Parse(codeTemplate)
	if err != nil {
		return err
	}

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	data := struct {
		Methods []RPCMethod
	}{
		Methods: methods,
	}

	return tmpl.Execute(f, data)
}
