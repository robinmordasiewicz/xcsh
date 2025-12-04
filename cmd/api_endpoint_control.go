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

	"github.com/robinmordasiewicz/f5xc/pkg/output"
)

var controlFlags struct {
	namespace      string
	name           string
	srcService     string
	dstService     string
	srcNamespace   string
	dstNamespace   string
	virtualHost    string
	method         []string
	path           []string
	pathRegex      []string
	action         string
	rateLimitQPS   int
	rateLimitBurst int
	inputFile      string
	dryRun         bool
}

var controlCmd = &cobra.Command{
	Use:   "control",
	Short: "Create layer7 policy for API endpoints",
	Long: `Create layer7 security policies based on discovered API endpoints.

This command allows you to create policies that control access to
specific API endpoints, including:
  - Allow/Deny rules for specific methods and paths
  - Rate limiting per endpoint
  - Cross-namespace policy application

Policies can be defined via command-line flags or input files.`,
	Example: `  # Allow GET requests to /api/users
  f5xc api-endpoint control --name allow-users-read \
    --src-service frontend --dst-service backend \
    --method GET --path /api/users --action allow

  # Deny all requests to admin endpoints
  f5xc api-endpoint control --name deny-admin \
    --dst-service backend --path-regex "/admin/.*" --action deny

  # Rate limit an endpoint
  f5xc api-endpoint control --name rate-limit-api \
    --dst-service api-server --path /api/search \
    --rate-limit-qps 100 --rate-limit-burst 200

  # Create policy from file
  f5xc api-endpoint control -i policy.yaml

  # Dry run to preview the policy
  f5xc api-endpoint control --name test-policy --dry-run --dst-service api`,
	RunE: runControl,
}

func init() {
	apiEndpointCmd.AddCommand(controlCmd)

	controlCmd.Flags().StringVarP(&controlFlags.namespace, "namespace", "n", "shared", "Namespace for the policy")
	controlCmd.Flags().StringVar(&controlFlags.name, "name", "", "Policy name")
	controlCmd.Flags().StringVar(&controlFlags.srcService, "src-service", "", "Source service name")
	controlCmd.Flags().StringVar(&controlFlags.dstService, "dst-service", "", "Destination service name")
	controlCmd.Flags().StringVar(&controlFlags.srcNamespace, "src-namespace", "", "Source service namespace")
	controlCmd.Flags().StringVar(&controlFlags.dstNamespace, "dst-namespace", "", "Destination service namespace")
	controlCmd.Flags().StringVar(&controlFlags.virtualHost, "virtual-host", "", "Virtual host to apply policy")
	controlCmd.Flags().StringSliceVar(&controlFlags.method, "method", nil, "HTTP methods to match (GET, POST, PUT, DELETE, etc.)")
	controlCmd.Flags().StringSliceVar(&controlFlags.path, "path", nil, "Exact paths to match")
	controlCmd.Flags().StringSliceVar(&controlFlags.pathRegex, "path-regex", nil, "Regex patterns for paths to match")
	controlCmd.Flags().StringVar(&controlFlags.action, "action", "allow", "Policy action: allow, deny")
	controlCmd.Flags().IntVar(&controlFlags.rateLimitQPS, "rate-limit-qps", 0, "Rate limit queries per second (0 = no limit)")
	controlCmd.Flags().IntVar(&controlFlags.rateLimitBurst, "rate-limit-burst", 0, "Rate limit burst size")
	controlCmd.Flags().StringVarP(&controlFlags.inputFile, "input-file", "i", "", "Input file (YAML/JSON) containing policy definition")
	controlCmd.Flags().BoolVar(&controlFlags.dryRun, "dry-run", false, "Preview the policy without creating it")
}

// ServicePolicy represents a layer7 service policy
type ServicePolicy struct {
	Metadata PolicyMetadata `json:"metadata" yaml:"metadata"`
	Spec     PolicySpec     `json:"spec" yaml:"spec"`
}

type PolicyMetadata struct {
	Name        string            `json:"name" yaml:"name"`
	Namespace   string            `json:"namespace" yaml:"namespace"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

type PolicySpec struct {
	Algo        string       `json:"algo,omitempty" yaml:"algo,omitempty"`
	AnyServer   bool         `json:"any_server,omitempty" yaml:"any_server,omitempty"`
	ServerName  string       `json:"server_name,omitempty" yaml:"server_name,omitempty"`
	Rules       []PolicyRule `json:"rules,omitempty" yaml:"rules,omitempty"`
	RuleList    *RuleList    `json:"rule_list,omitempty" yaml:"rule_list,omitempty"`
}

type RuleList struct {
	Rules []PolicyRule `json:"rules" yaml:"rules"`
}

type PolicyRule struct {
	Metadata    RuleMetadata        `json:"metadata" yaml:"metadata"`
	Spec        RuleSpec            `json:"spec" yaml:"spec"`
}

type RuleMetadata struct {
	Name string `json:"name" yaml:"name"`
}

type RuleSpec struct {
	Action          string           `json:"action" yaml:"action"`
	AnyClient       bool             `json:"any_client,omitempty" yaml:"any_client,omitempty"`
	ClientSelector  *ClientSelector  `json:"client_selector,omitempty" yaml:"client_selector,omitempty"`
	APIGroupMatcher *APIGroupMatcher `json:"api_group_matcher,omitempty" yaml:"api_group_matcher,omitempty"`
	RequestMatcher  *RequestMatcher  `json:"request_matcher,omitempty" yaml:"request_matcher,omitempty"`
	RateLimiter     *RateLimiter     `json:"rate_limiter,omitempty" yaml:"rate_limiter,omitempty"`
}

type ClientSelector struct {
	Expressions []string `json:"expressions,omitempty" yaml:"expressions,omitempty"`
}

type APIGroupMatcher struct {
	Match []APIGroupMatch `json:"match,omitempty" yaml:"match,omitempty"`
}

type APIGroupMatch struct {
	Method []string `json:"method,omitempty" yaml:"method,omitempty"`
	Path   *PathMatcher `json:"path,omitempty" yaml:"path,omitempty"`
}

type PathMatcher struct {
	Prefix string `json:"prefix,omitempty" yaml:"prefix,omitempty"`
	Exact  string `json:"exact,omitempty" yaml:"exact,omitempty"`
	Regex  string `json:"regex,omitempty" yaml:"regex,omitempty"`
}

type RequestMatcher struct {
	HttpMethod []string     `json:"http_method,omitempty" yaml:"http_method,omitempty"`
	Path       *PathMatcher `json:"path,omitempty" yaml:"path,omitempty"`
}

type RateLimiter struct {
	TotalNumber    int `json:"total_number,omitempty" yaml:"total_number,omitempty"`
	Unit           string `json:"unit,omitempty" yaml:"unit,omitempty"`
	BurstMultiplier int `json:"burst_multiplier,omitempty" yaml:"burst_multiplier,omitempty"`
}

func runControl(cmd *cobra.Command, args []string) error {
	apiClient := GetClient()
	if apiClient == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	var policy *ServicePolicy
	var err error

	if controlFlags.inputFile != "" {
		policy, err = loadPolicyFromFile(controlFlags.inputFile)
		if err != nil {
			return fmt.Errorf("failed to load policy file: %w", err)
		}
	} else {
		if controlFlags.name == "" {
			return fmt.Errorf("--name is required when not using input file")
		}
		policy = buildPolicyFromFlags()
	}

	// Validate the policy
	if err := validatePolicy(policy); err != nil {
		return fmt.Errorf("invalid policy: %w", err)
	}

	// Dry run - just print the policy
	if controlFlags.dryRun {
		output.PrintInfo("Dry run - policy preview:")
		return output.Print(policy, GetOutputFormat())
	}

	// Create the policy
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	requestBody := map[string]interface{}{
		"metadata": policy.Metadata,
		"spec":     policy.Spec,
	}

	path := fmt.Sprintf("/api/config/namespaces/%s/service_policys", policy.Metadata.Namespace)
	resp, err := apiClient.Post(ctx, path, requestBody)
	if err != nil {
		return fmt.Errorf("failed to create service policy: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(resp.Body))
	}

	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	output.PrintInfo(fmt.Sprintf("Service policy '%s' created successfully", policy.Metadata.Name))
	return output.Print(result, GetOutputFormat())
}

func loadPolicyFromFile(path string) (*ServicePolicy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var policy ServicePolicy
	if err := yaml.Unmarshal(data, &policy); err != nil {
		if err := json.Unmarshal(data, &policy); err != nil {
			return nil, fmt.Errorf("failed to parse policy file (not valid YAML or JSON): %w", err)
		}
	}

	// Apply flag overrides
	if controlFlags.namespace != "" && policy.Metadata.Namespace == "" {
		policy.Metadata.Namespace = controlFlags.namespace
	}
	if controlFlags.name != "" && policy.Metadata.Name == "" {
		policy.Metadata.Name = controlFlags.name
	}

	return &policy, nil
}

func buildPolicyFromFlags() *ServicePolicy {
	policy := &ServicePolicy{
		Metadata: PolicyMetadata{
			Name:      controlFlags.name,
			Namespace: controlFlags.namespace,
			Labels: map[string]string{
				"created-by": "f5xc-cli",
			},
		},
		Spec: PolicySpec{
			Algo: "FIRST_MATCH",
		},
	}

	// Set server targeting
	if controlFlags.dstService != "" {
		policy.Spec.ServerName = controlFlags.dstService
	} else {
		policy.Spec.AnyServer = true
	}

	// Build the rule
	rule := PolicyRule{
		Metadata: RuleMetadata{
			Name: fmt.Sprintf("%s-rule", controlFlags.name),
		},
		Spec: RuleSpec{
			Action: strings.ToUpper(controlFlags.action),
		},
	}

	// Client selector
	if controlFlags.srcService != "" {
		srcNs := controlFlags.srcNamespace
		if srcNs == "" {
			srcNs = controlFlags.namespace
		}
		rule.Spec.ClientSelector = &ClientSelector{
			Expressions: []string{
				fmt.Sprintf("service.metadata.name in {%q}", controlFlags.srcService),
			},
		}
	} else {
		rule.Spec.AnyClient = true
	}

	// Build API group matcher for methods and paths
	if len(controlFlags.method) > 0 || len(controlFlags.path) > 0 || len(controlFlags.pathRegex) > 0 {
		apiMatches := []APIGroupMatch{}

		// Handle exact paths
		for _, p := range controlFlags.path {
			match := APIGroupMatch{
				Path: &PathMatcher{Exact: p},
			}
			if len(controlFlags.method) > 0 {
				match.Method = controlFlags.method
			}
			apiMatches = append(apiMatches, match)
		}

		// Handle regex paths
		for _, p := range controlFlags.pathRegex {
			match := APIGroupMatch{
				Path: &PathMatcher{Regex: p},
			}
			if len(controlFlags.method) > 0 {
				match.Method = controlFlags.method
			}
			apiMatches = append(apiMatches, match)
		}

		// If only methods specified, match all paths with those methods
		if len(controlFlags.path) == 0 && len(controlFlags.pathRegex) == 0 && len(controlFlags.method) > 0 {
			apiMatches = append(apiMatches, APIGroupMatch{
				Method: controlFlags.method,
				Path:   &PathMatcher{Prefix: "/"},
			})
		}

		if len(apiMatches) > 0 {
			rule.Spec.APIGroupMatcher = &APIGroupMatcher{
				Match: apiMatches,
			}
		}
	}

	// Rate limiter
	if controlFlags.rateLimitQPS > 0 {
		rule.Spec.RateLimiter = &RateLimiter{
			TotalNumber: controlFlags.rateLimitQPS,
			Unit:        "SECOND",
		}
		if controlFlags.rateLimitBurst > 0 {
			burstMultiplier := controlFlags.rateLimitBurst / controlFlags.rateLimitQPS
			if burstMultiplier < 1 {
				burstMultiplier = 1
			}
			rule.Spec.RateLimiter.BurstMultiplier = burstMultiplier
		}
	}

	policy.Spec.RuleList = &RuleList{
		Rules: []PolicyRule{rule},
	}

	return policy
}

func validatePolicy(policy *ServicePolicy) error {
	if policy.Metadata.Name == "" {
		return fmt.Errorf("policy name is required")
	}

	if policy.Metadata.Namespace == "" {
		return fmt.Errorf("policy namespace is required")
	}

	// Validate action
	if policy.Spec.RuleList != nil {
		for _, rule := range policy.Spec.RuleList.Rules {
			action := strings.ToUpper(rule.Spec.Action)
			switch action {
			case "ALLOW", "DENY", "NEXT_POLICY":
				// valid
			default:
				return fmt.Errorf("invalid action '%s' - must be ALLOW, DENY, or NEXT_POLICY", rule.Spec.Action)
			}
		}
	}

	return nil
}
