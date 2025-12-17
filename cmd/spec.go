package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/robinmordasiewicz/f5xcctl/pkg/errors"
	"github.com/robinmordasiewicz/f5xcctl/pkg/subscription"
	"github.com/robinmordasiewicz/f5xcctl/pkg/types"
)

// specFlag controls whether to output machine-readable spec
var specFlag bool

// RegisterSpecFlag registers the --spec flag on the root command
// This should be called from root.go init after rootCmd is defined
func RegisterSpecFlag(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().BoolVar(&specFlag, "spec", false, "Output machine-readable CLI specification in JSON or YAML.")
}

// CLISpec represents the complete CLI specification
type CLISpec struct {
	Name                  string                                `json:"name" yaml:"name"`
	Version               string                                `json:"version" yaml:"version"`
	Description           string                                `json:"description" yaml:"description"`
	Usage                 string                                `json:"usage" yaml:"usage"`
	AIHints               AIHintsSpec                           `json:"ai_hints" yaml:"ai_hints"`
	AuthenticationMethods []AuthMethodSpec                      `json:"authentication_methods" yaml:"authentication_methods"`
	SemanticCategories    SemanticCategoriesSpec                `json:"semantic_categories" yaml:"semantic_categories"`
	FlagRelationships     FlagRelationshipsSpec                 `json:"flag_relationships" yaml:"flag_relationships"`
	Examples              []ExampleSpec                         `json:"examples" yaml:"examples"`
	Workflows             []WorkflowSpec                        `json:"workflows" yaml:"workflows"`
	EnvironmentVariables  []EnvVarSpec                          `json:"environment_variables" yaml:"environment_variables"`
	GlobalFlags           []FlagSpec                            `json:"global_flags" yaml:"global_flags"`
	Commands              []CommandSpec                         `json:"commands" yaml:"commands"`
	ExitCodes             []ExitCodeSpec                        `json:"exit_codes" yaml:"exit_codes"`
	SystemLabels          SystemLabelsSpec                      `json:"system_labels" yaml:"system_labels"`
	ResourceSchemas       map[string]types.ResourceSchemaInfo   `json:"resource_schemas" yaml:"resource_schemas"`
	SubscriptionContext   *subscription.SubscriptionContextSpec `json:"subscription_context,omitempty" yaml:"subscription_context,omitempty"`
	FeatureTierMap        subscription.FeatureTierMapSpec       `json:"feature_tier_map" yaml:"feature_tier_map"`
}

// AIHintsSpec provides guidance for AI agents on how to use the CLI
type AIHintsSpec struct {
	DiscoveryCommand      string   `json:"discovery_command" yaml:"discovery_command"`
	RecommendedAuthMethod string   `json:"recommended_auth_method" yaml:"recommended_auth_method"`
	RequiredForAPICalls   []string `json:"required_for_api_calls" yaml:"required_for_api_calls"`
	OutputFormats         []string `json:"output_formats" yaml:"output_formats"`
	DefaultOutputFormat   string   `json:"default_output_format" yaml:"default_output_format"`
	BestPractices         []string `json:"best_practices" yaml:"best_practices"`
}

// AuthMethodSpec describes an authentication method
type AuthMethodSpec struct {
	Method      string   `json:"method" yaml:"method"`
	Description string   `json:"description" yaml:"description"`
	Flags       []string `json:"flags" yaml:"flags"`
	EnvVars     []string `json:"env_vars" yaml:"env_vars"`
	Priority    int      `json:"priority" yaml:"priority"`
}

// SemanticCategoriesSpec groups flags and env vars by purpose
type SemanticCategoriesSpec struct {
	Authentication []string `json:"authentication" yaml:"authentication"`
	Connection     []string `json:"connection" yaml:"connection"`
	Output         []string `json:"output" yaml:"output"`
	Behavior       []string `json:"behavior" yaml:"behavior"`
}

// FlagRelationshipsSpec describes flag dependencies and conflicts
type FlagRelationshipsSpec struct {
	MutuallyExclusive [][]string           `json:"mutually_exclusive" yaml:"mutually_exclusive"`
	RequiredTogether  [][]string           `json:"required_together" yaml:"required_together"`
	Dependencies      []FlagDependencySpec `json:"dependencies" yaml:"dependencies"`
}

// FlagDependencySpec describes a flag dependency
type FlagDependencySpec struct {
	Flag     string   `json:"flag" yaml:"flag"`
	Requires []string `json:"requires" yaml:"requires"`
}

// ExampleSpec provides a structured usage example
type ExampleSpec struct {
	Task          string            `json:"task" yaml:"task"`
	Command       string            `json:"command" yaml:"command"`
	Description   string            `json:"description" yaml:"description"`
	Category      string            `json:"category" yaml:"category"`
	Prerequisites []string          `json:"prerequisites,omitempty" yaml:"prerequisites,omitempty"`
	EnvVars       map[string]string `json:"env_vars,omitempty" yaml:"env_vars,omitempty"`
}

// WorkflowSpec describes a multi-step workflow
type WorkflowSpec struct {
	Name        string         `json:"name" yaml:"name"`
	Description string         `json:"description" yaml:"description"`
	Steps       []WorkflowStep `json:"steps" yaml:"steps"`
}

// WorkflowStep represents a single step in a workflow
type WorkflowStep struct {
	Step        int    `json:"step" yaml:"step"`
	Description string `json:"description" yaml:"description"`
	Command     string `json:"command" yaml:"command"`
	Optional    bool   `json:"optional,omitempty" yaml:"optional,omitempty"`
}

// CommandSpec represents a command's specification
type CommandSpec struct {
	Path        []string      `json:"path" yaml:"path"`
	Use         string        `json:"use" yaml:"use"`
	Short       string        `json:"short" yaml:"short"`
	Long        string        `json:"long,omitempty" yaml:"long,omitempty"`
	Aliases     []string      `json:"aliases,omitempty" yaml:"aliases,omitempty"`
	Example     string        `json:"example,omitempty" yaml:"example,omitempty"`
	Flags       []FlagSpec    `json:"flags,omitempty" yaml:"flags,omitempty"`
	Subcommands []CommandSpec `json:"subcommands,omitempty" yaml:"subcommands,omitempty"`
}

// FlagSpec represents a flag's specification
type FlagSpec struct {
	Name        string `json:"name" yaml:"name"`
	Shorthand   string `json:"shorthand,omitempty" yaml:"shorthand,omitempty"`
	Type        string `json:"type" yaml:"type"`
	Default     string `json:"default,omitempty" yaml:"default,omitempty"`
	Description string `json:"description" yaml:"description"`
	Required    bool   `json:"required,omitempty" yaml:"required,omitempty"`
	Hidden      bool   `json:"hidden,omitempty" yaml:"hidden,omitempty"`
}

// ExitCodeSpec represents an exit code's specification
type ExitCodeSpec struct {
	Code        int    `json:"code" yaml:"code"`
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
}

// SystemLabelsSpec documents system-managed label patterns for AI assistants.
// F5 XC automatically injects labels with ves.io/* prefix that should not be
// manually set by users or automation tools.
type SystemLabelsSpec struct {
	Description string               `json:"description" yaml:"description"`
	Patterns    []SystemLabelPattern `json:"patterns" yaml:"patterns"`
	Guidance    []string             `json:"guidance" yaml:"guidance"`
}

// SystemLabelPattern describes a system-managed label pattern
type SystemLabelPattern struct {
	Pattern        string `json:"pattern" yaml:"pattern"`
	Description    string `json:"description" yaml:"description"`
	ManagedBy      string `json:"managed_by" yaml:"managed_by"`
	UserModifiable bool   `json:"user_modifiable" yaml:"user_modifiable"`
}

// GenerateSpec generates the CLI specification
func GenerateSpec(cmd *cobra.Command) *CLISpec {
	spec := &CLISpec{
		Name:                  "f5xcctl",
		Version:               Version, // From version.go
		Description:           cmd.Long,
		Usage:                 cmd.Use,
		AIHints:               getAIHints(),
		AuthenticationMethods: getAuthenticationMethods(),
		SemanticCategories:    getSemanticCategories(),
		FlagRelationships:     getFlagRelationships(),
		Examples:              getExamples(),
		Workflows:             getWorkflows(),
		EnvironmentVariables:  EnvVarRegistry,
		GlobalFlags:           extractFlags(cmd.PersistentFlags()),
		Commands:              extractCommands(cmd, []string{}),
		ExitCodes:             getExitCodes(),
		SystemLabels:          getSystemLabels(),
		ResourceSchemas:       types.GetAllResourceSchemas(),
		FeatureTierMap:        subscription.GenerateFeatureTierMapSpec(),
	}

	// Try to populate subscription context if validator is available
	validator := GetSubscriptionValidator()
	if validator != nil {
		ctx := context.Background()
		if subCtx, err := validator.GenerateSubscriptionContext(ctx); err == nil {
			spec.SubscriptionContext = subCtx
		}
	} else {
		// Provide default context based on cached tier or default
		spec.SubscriptionContext = getDefaultSubscriptionContext()
	}

	return spec
}

// getDefaultSubscriptionContext returns a subscription context when validator is unavailable
func getDefaultSubscriptionContext() *subscription.SubscriptionContextSpec {
	cacheInfo := subscription.GetTierCacheInfo()

	tier := cacheInfo.Tier
	if tier == "" {
		tier = "Unknown"
	}

	tierSource := "Unknown"
	if cacheInfo.IsFromEnv {
		tierSource = "F5XC_SUBSCRIPTION_TIER environment variable"
	} else if cacheInfo.IsFromAPI {
		tierSource = "Detected from F5 XC API"
	} else if cacheInfo.IsDefault {
		tierSource = "Default (API unavailable)"
	} else if tier != "Unknown" {
		tierSource = "Cached value"
	}

	// Build feature lists from registry
	registry := subscription.NewFeatureRegistry()

	var available []string
	var restricted []subscription.RestrictedFeatureSpec

	for _, feature := range registry.GetAllFeatures() {
		if subscription.IsTierCached() {
			// Check if current tier meets minimum requirement for this feature
			if isTierSufficientForSpec(tier, feature.MinimumTier) {
				available = append(available, feature.FeatureName)
			} else {
				restricted = append(restricted, subscription.RestrictedFeatureSpec{
					Feature:        feature.FeatureName,
					DisplayName:    feature.DisplayName,
					RequiredTier:   feature.MinimumTier,
					RequiredAddons: feature.RequiredAddons,
					HelpAnnotation: feature.HelpAnnotation,
				})
			}
		} else {
			// If tier unknown, list all features as available with notes
			available = append(available, feature.FeatureName)
		}
	}

	return &subscription.SubscriptionContextSpec{
		CurrentTier:        tier,
		TierSource:         tierSource,
		AvailableFeatures:  available,
		RestrictedFeatures: restricted,
	}
}

// isTierSufficientForSpec checks if currentTier meets requiredTier for spec output
func isTierSufficientForSpec(currentTier, requiredTier string) bool {
	// Tier order: Standard (2) < Advanced (3)
	tierOrder := map[string]int{
		"":                        0,
		subscription.TierNoTier:   0,
		subscription.TierStandard: 2,
		"Standard":                2,
		subscription.TierAdvanced: 3,
		"Advanced":                3,
	}

	currentOrder := tierOrder[currentTier]
	requiredOrder := tierOrder[requiredTier]

	return currentOrder >= requiredOrder
}

// getSystemLabels returns documentation about system-managed labels
func getSystemLabels() SystemLabelsSpec {
	return SystemLabelsSpec{
		Description: "F5 Distributed Cloud automatically injects labels with ves.io/* prefix. " +
			"These labels are system-managed and should NOT be manually set in configuration files. " +
			"The API will overwrite or ignore user-provided ves.io/* labels.",
		Patterns: []SystemLabelPattern{
			{
				Pattern:        "ves.io/app_type",
				Description:    "Application type label auto-generated from resource name during creation",
				ManagedBy:      "api-server",
				UserModifiable: false,
			},
			{
				Pattern:        "ves.io/soft-deleted",
				Description:    "Marks resources pending garbage collection after deletion request",
				ManagedBy:      "garbage-collector",
				UserModifiable: false,
			},
			{
				Pattern:        "ves.io/site-*",
				Description:    "Site identification and topology labels for distributed resources",
				ManagedBy:      "site-controller",
				UserModifiable: false,
			},
			{
				Pattern:        "ves.io/tenant",
				Description:    "Tenant identification derived from authentication context",
				ManagedBy:      "api-gateway",
				UserModifiable: false,
			},
			{
				Pattern:        "ves.io/managed-by-*",
				Description:    "Indicates resources managed by specific controllers or automation",
				ManagedBy:      "various-controllers",
				UserModifiable: false,
			},
			{
				Pattern:        "ves.io/fleet",
				Description:    "Fleet membership label for grouped site management",
				ManagedBy:      "fleet-controller",
				UserModifiable: false,
			},
		},
		Guidance: []string{
			"Do NOT set labels with 'ves.io/' prefix in configuration files",
			"System labels are automatically managed and will be overwritten",
			"Use custom label prefixes (e.g., 'myorg.com/') for user-managed labels",
			"When comparing configurations, ignore ves.io/* labels as they may differ between environments",
			"Terraform and automation tools should use lifecycle.ignore_changes for system labels",
		},
	}
}

// getAIHints returns AI agent guidance
func getAIHints() AIHintsSpec {
	return AIHintsSpec{
		DiscoveryCommand:      "f5xcctl --spec --output-format json",
		RecommendedAuthMethod: "p12-bundle",
		RequiredForAPICalls:   []string{"server-url", "authentication (one of: p12-bundle, cert+key, api-token)"},
		OutputFormats:         []string{"json", "yaml", "table", "text"},
		DefaultOutputFormat:   "table",
		BestPractices: []string{
			"Use --output-format json for programmatic parsing",
			"Set F5XC_API_URL and F5XC_P12_FILE environment variables for persistent configuration",
			"Use f5xcctl configuration list <resource-type> to discover available resources",
			"Use f5xcctl configuration get <resource-type> -n <namespace> <name> to retrieve specific resources",
			"Always specify --namespace or -n for namespace-scoped resources",
			"Use --spec to get complete CLI structure before constructing commands",
			"Check exit codes for programmatic error handling (0=success, 1=generic, 2=validation, 3=auth, 4=connection, 5=not-found, 6=conflict, 7=rate-limit, 8=quota-exceeded, 9=feature-unavailable)",
			"Use resource_schemas in --spec output to understand field constraints, validation rules, and mutual exclusivity",
			"Follow oneof_groups in resource schemas to ensure only one choice per mutually exclusive group is configured",
			"Use decision_tree in resource schemas to determine required fields based on configuration choices",
			"Use 'f5xcctl subscription show' to identify tenant tier before deploying tier-specific features",
			"Use 'f5xcctl subscription validate' before terraform apply to catch quota/feature issues early",
			"Check addon service status with 'f5xcctl subscription addons' before using features like bot-defense",
			"Monitor quota usage with 'f5xcctl subscription quota' to prevent deployment failures",
		},
	}
}

// getAuthenticationMethods returns structured authentication options
func getAuthenticationMethods() []AuthMethodSpec {
	return []AuthMethodSpec{
		{
			Method:      "p12-bundle",
			Description: "PKCS#12 certificate bundle (recommended for automation)",
			Flags:       []string{"--p12-bundle"},
			EnvVars:     []string{"F5XC_P12_FILE", "F5XC_P12_PASSWORD"},
			Priority:    1,
		},
		{
			Method:      "certificate",
			Description: "Separate certificate and key files for mTLS",
			Flags:       []string{"--cert", "--key"},
			EnvVars:     []string{"F5XC_CERT", "F5XC_KEY"},
			Priority:    2,
		},
		{
			Method:      "api-token",
			Description: "API token authentication",
			Flags:       []string{"--api-token"},
			EnvVars:     []string{"F5XC_API_TOKEN"},
			Priority:    3,
		},
		{
			Method:      "hardware-key",
			Description: "YubiKey hardware security module",
			Flags:       []string{"--hardwareKey"},
			EnvVars:     []string{},
			Priority:    4,
		},
	}
}

// getSemanticCategories groups flags by purpose
func getSemanticCategories() SemanticCategoriesSpec {
	return SemanticCategoriesSpec{
		Authentication: []string{
			"--p12-bundle", "--cert", "--key", "--api-token", "--hardwareKey",
			"F5XC_P12_FILE", "F5XC_P12_PASSWORD", "F5XC_CERT", "F5XC_KEY", "F5XC_API_TOKEN",
		},
		Connection: []string{
			"--server-url", "--cacert", "--timeout", "--config",
			"F5XC_API_URL", "F5XC_CACERT", "F5XC_CONFIG",
		},
		Output: []string{
			"--output-format", "--output", "--show-curl", "--spec",
			"F5XC_OUTPUT",
		},
		Behavior: []string{
			"--non-interactive", "--help",
		},
	}
}

// getFlagRelationships returns flag dependencies and conflicts
func getFlagRelationships() FlagRelationshipsSpec {
	return FlagRelationshipsSpec{
		MutuallyExclusive: [][]string{
			{"--p12-bundle", "--cert"},
			{"--p12-bundle", "--key"},
			{"--p12-bundle", "--api-token"},
			{"--cert", "--api-token"},
			{"--key", "--api-token"},
		},
		RequiredTogether: [][]string{
			{"--cert", "--key"},
		},
		Dependencies: []FlagDependencySpec{
			{
				Flag:     "--p12-bundle",
				Requires: []string{"F5XC_P12_PASSWORD (env var)"},
			},
			{
				Flag:     "--api-token",
				Requires: []string{"F5XC_API_TOKEN (env var)"},
			},
		},
	}
}

// getExamples returns structured usage examples
func getExamples() []ExampleSpec {
	return []ExampleSpec{
		{
			Task:        "List all namespaces",
			Command:     "f5xcctl configuration list namespace",
			Description: "Retrieve all namespaces accessible to the authenticated user",
			Category:    "discovery",
		},
		{
			Task:        "List HTTP load balancers in a namespace",
			Command:     "f5xcctl configuration list http_loadbalancer -n <namespace>",
			Description: "List all HTTP load balancer configurations in the specified namespace",
			Category:    "configuration",
		},
		{
			Task:        "Get a specific HTTP load balancer",
			Command:     "f5xcctl configuration get http_loadbalancer -n <namespace> <name>",
			Description: "Retrieve detailed configuration of a specific HTTP load balancer",
			Category:    "configuration",
		},
		{
			Task:        "Create a resource from YAML file",
			Command:     "f5xcctl configuration create http_loadbalancer -n <namespace> -i <file.yaml>",
			Description: "Create a new HTTP load balancer from a YAML specification file",
			Category:    "configuration",
		},
		{
			Task:        "Delete a resource",
			Command:     "f5xcctl configuration delete http_loadbalancer -n <namespace> <name>",
			Description: "Delete an HTTP load balancer by name",
			Category:    "configuration",
		},
		{
			Task:        "List origin pools",
			Command:     "f5xcctl configuration list origin_pool -n <namespace>",
			Description: "List all origin pool configurations in the specified namespace",
			Category:    "configuration",
		},
		{
			Task:        "Output as JSON for parsing",
			Command:     "f5xcctl configuration list namespace --output-format json",
			Description: "Get namespace list in JSON format for programmatic processing",
			Category:    "output",
		},
		{
			Task:        "Discover API endpoints",
			Command:     "f5xcctl api-endpoint list -n <namespace>",
			Description: "List discovered API endpoints within the service mesh",
			Category:    "discovery",
		},
		{
			Task:        "Get CLI specification",
			Command:     "f5xcctl --spec --output-format json",
			Description: "Output complete CLI specification in JSON format for AI/automation tools",
			Category:    "meta",
		},
		{
			Task:        "Show version and build info",
			Command:     "f5xcctl version",
			Description: "Display f5xcctl version, commit, build date, and platform information",
			Category:    "meta",
		},
		{
			Task:        "Show subscription summary",
			Command:     "f5xcctl subscription show --output-format json",
			Description: "Display subscription tier, active addons, and quota summary for AI assistants",
			Category:    "subscription",
		},
		{
			Task:        "List active addon services",
			Command:     "f5xcctl subscription addons --filter active",
			Description: "List all actively subscribed addon services",
			Category:    "subscription",
		},
		{
			Task:        "Check quota usage",
			Command:     "f5xcctl subscription quota -n <namespace>",
			Description: "Display quota limits and current usage in a namespace",
			Category:    "subscription",
		},
		{
			Task:        "Validate resource creation",
			Command:     "f5xcctl subscription validate --resource-type http_loadbalancer --count 5",
			Description: "Check if 5 additional HTTP load balancers can be created within quota",
			Category:    "subscription",
		},
		{
			Task:        "Validate feature availability",
			Command:     "f5xcctl subscription validate --feature bot-defense",
			Description: "Verify if bot-defense addon service is subscribed and available",
			Category:    "subscription",
		},
	}
}

// getWorkflows returns multi-step workflow definitions
func getWorkflows() []WorkflowSpec {
	return []WorkflowSpec{
		{
			Name:        "authenticate-and-discover",
			Description: "Set up authentication and discover available resources",
			Steps: []WorkflowStep{
				{Step: 1, Description: "Set API URL", Command: "export F5XC_API_URL=https://<tenant>.console.ves.volterra.io/api"},
				{Step: 2, Description: "Set P12 credentials", Command: "export F5XC_P12_FILE=/path/to/api-creds.p12 && export F5XC_P12_PASSWORD=<password>"},
				{Step: 3, Description: "Verify authentication", Command: "f5xcctl configuration list namespace"},
				{Step: 4, Description: "Discover resource types", Command: "f5xcctl --help"},
			},
		},
		{
			Name:        "deploy-http-load-balancer",
			Description: "Create an HTTP load balancer with origin pool",
			Steps: []WorkflowStep{
				{Step: 1, Description: "Create origin pool", Command: "f5xcctl configuration create origin_pool -n <namespace> -i origin-pool.yaml"},
				{Step: 2, Description: "Create HTTP load balancer", Command: "f5xcctl configuration create http_loadbalancer -n <namespace> -i http-lb.yaml"},
				{Step: 3, Description: "Verify deployment", Command: "f5xcctl configuration get http_loadbalancer -n <namespace> <name>"},
			},
		},
		{
			Name:        "export-configuration",
			Description: "Export existing configuration for backup or migration",
			Steps: []WorkflowStep{
				{Step: 1, Description: "List resources", Command: "f5xcctl configuration list http_loadbalancer -n <namespace> --output-format json"},
				{Step: 2, Description: "Get specific resource as YAML", Command: "f5xcctl configuration get http_loadbalancer -n <namespace> <name> --output-format yaml > backup.yaml"},
			},
		},
		{
			Name:        "update-configuration",
			Description: "Modify an existing resource configuration",
			Steps: []WorkflowStep{
				{Step: 1, Description: "Export current config", Command: "f5xcctl configuration get http_loadbalancer -n <namespace> <name> --output-format yaml > current.yaml"},
				{Step: 2, Description: "Edit configuration", Command: "# Edit current.yaml with desired changes"},
				{Step: 3, Description: "Apply changes", Command: "f5xcctl configuration replace http_loadbalancer -n <namespace> -i current.yaml"},
				{Step: 4, Description: "Verify update", Command: "f5xcctl configuration get http_loadbalancer -n <namespace> <name>"},
			},
		},
		{
			Name:        "pre-deployment-validation",
			Description: "Validate subscription capabilities before Terraform deployment",
			Steps: []WorkflowStep{
				{Step: 1, Description: "Check subscription tier", Command: "f5xcctl subscription show --output-format json"},
				{Step: 2, Description: "Verify addon services", Command: "f5xcctl subscription addons --filter active --output-format json"},
				{Step: 3, Description: "Check quota availability", Command: "f5xcctl subscription quota -n <namespace> --output-format json"},
				{Step: 4, Description: "Validate specific resources", Command: "f5xcctl subscription validate --resource-type http_loadbalancer --count <planned_count>"},
				{Step: 5, Description: "Validate required features", Command: "f5xcctl subscription validate --feature <required_addon>", Optional: true},
				{Step: 6, Description: "Proceed with deployment", Command: "terraform apply"},
			},
		},
		{
			Name:        "quota-troubleshooting",
			Description: "Diagnose and resolve quota-related deployment failures",
			Steps: []WorkflowStep{
				{Step: 1, Description: "Check overall quota status", Command: "f5xcctl subscription quota -n <namespace>"},
				{Step: 2, Description: "Identify exceeded quotas", Command: "f5xcctl subscription quota -n <namespace> --output-format json | jq '.objects[] | select(.status == \"EXCEEDED\")'"},
				{Step: 3, Description: "List resources of exceeded type", Command: "f5xcctl configuration list <resource_type> -n <namespace>"},
				{Step: 4, Description: "Clean up unused resources", Command: "f5xcctl configuration delete <resource_type> -n <namespace> <unused_name>"},
				{Step: 5, Description: "Revalidate quota", Command: "f5xcctl subscription validate --resource-type <resource_type> --count 1"},
			},
		},
	}
}

// extractCommands recursively extracts command specifications
func extractCommands(cmd *cobra.Command, parentPath []string) []CommandSpec {
	var commands []CommandSpec

	for _, subCmd := range cmd.Commands() {
		// Skip hidden commands
		if subCmd.Hidden {
			continue
		}

		cmdPath := append(parentPath, subCmd.Name())
		cmdSpec := CommandSpec{
			Path:    cmdPath,
			Use:     subCmd.Use,
			Short:   subCmd.Short,
			Long:    subCmd.Long,
			Aliases: subCmd.Aliases,
			Example: subCmd.Example,
			Flags:   extractFlags(subCmd.LocalFlags()),
		}

		// Recursively extract subcommands
		if subCmd.HasSubCommands() {
			cmdSpec.Subcommands = extractCommands(subCmd, cmdPath)
		}

		commands = append(commands, cmdSpec)
	}

	// Sort commands alphabetically
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Use < commands[j].Use
	})

	return commands
}

// extractFlags extracts flag specifications from a FlagSet
func extractFlags(flags *pflag.FlagSet) []FlagSpec {
	var flagSpecs []FlagSpec

	flags.VisitAll(func(f *pflag.Flag) {
		// Skip hidden flags
		if f.Hidden {
			return
		}

		flagSpec := FlagSpec{
			Name:        f.Name,
			Shorthand:   f.Shorthand,
			Type:        f.Value.Type(),
			Default:     f.DefValue,
			Description: f.Usage,
			Required:    isRequiredFlag(f),
		}

		flagSpecs = append(flagSpecs, flagSpec)
	})

	// Sort flags alphabetically
	sort.Slice(flagSpecs, func(i, j int) bool {
		return flagSpecs[i].Name < flagSpecs[j].Name
	})

	return flagSpecs
}

// isRequiredFlag checks if a flag is marked as required via Cobra's annotation system.
// Cobra uses BashCompOneRequiredFlag annotation when MarkFlagRequired() is called.
func isRequiredFlag(f *pflag.Flag) bool {
	if f.Annotations == nil {
		return false
	}
	// Cobra sets this annotation when a flag is marked required
	_, ok := f.Annotations[cobra.BashCompOneRequiredFlag]
	return ok
}

// getExitCodes returns the exit code specifications
func getExitCodes() []ExitCodeSpec {
	return []ExitCodeSpec{
		{Code: errors.ExitSuccess, Name: "ExitSuccess", Description: "Success"},
		{Code: errors.ExitGenericError, Name: "ExitGenericError", Description: "Generic/unknown error"},
		{Code: errors.ExitValidationError, Name: "ExitValidationError", Description: "Invalid arguments or validation failure"},
		{Code: errors.ExitAuthError, Name: "ExitAuthError", Description: "Authentication or authorization failure"},
		{Code: errors.ExitConnectionError, Name: "ExitConnectionError", Description: "Connection or timeout to API"},
		{Code: errors.ExitNotFoundError, Name: "ExitNotFoundError", Description: "Resource not found"},
		{Code: errors.ExitConflictError, Name: "ExitConflictError", Description: "Resource conflict"},
		{Code: errors.ExitRateLimitError, Name: "ExitRateLimitError", Description: "Rate limited"},
		{Code: errors.ExitQuotaExceeded, Name: "ExitQuotaExceeded", Description: "Subscription quota exceeded"},
		{Code: errors.ExitFeatureNotAvail, Name: "ExitFeatureNotAvail", Description: "Feature not available in subscription"},
	}
}

// OutputSpec outputs the CLI specification in the requested format
func OutputSpec(cmd *cobra.Command, format string) error {
	spec := GenerateSpec(cmd)

	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(spec)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(spec)
	default:
		// Default to JSON for machine readability
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(spec)
	}
}

// CheckSpecFlag checks if --spec flag is set and outputs spec if so
// Returns true if spec was output (caller should exit)
func CheckSpecFlag() bool {
	return specFlag
}

// HandleSpecFlag handles the --spec flag and outputs the spec
// This should be called from PersistentPreRunE in root command
// Pass the root command to avoid circular dependency
func HandleSpecFlag(rootCmd *cobra.Command) error {
	if specFlag {
		format := GetOutputFormatWithDefault("json")
		if err := OutputSpec(rootCmd, format); err != nil {
			return fmt.Errorf("failed to output spec: %w", err)
		}
		// Exit after outputting spec
		os.Exit(0)
	}
	return nil
}
