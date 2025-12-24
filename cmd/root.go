package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/robinmordasiewicz/xcsh/pkg/branding"
	"github.com/robinmordasiewicz/xcsh/pkg/client"
	"github.com/robinmordasiewicz/xcsh/pkg/config"
	"github.com/robinmordasiewicz/xcsh/pkg/subscription"
	"github.com/robinmordasiewicz/xcsh/pkg/types"
	"github.com/robinmordasiewicz/xcsh/pkg/validation"
)

var (
	// Config file path
	cfgFile string

	// Connection settings
	serverURL   string
	cert        string
	key         string
	cacert      string
	p12Bundle   string
	hardwareKey bool // Use yubikey for TLS connection
	useAPIToken bool // Use API token from F5XC_API_TOKEN environment variable

	// Output control
	outputFormat string // Output format for command (canonical: --output-format)
	outputDir    string // Output dir for command

	// Behavior flags
	showCurl       bool // Emit requests from program in CURL format
	timeout        int  // Timeout (in seconds) for command to finish
	nonInteractive bool // Fail on missing arguments instead of prompting

	// Internal flags (not exposed to CLI)
	debug bool

	// Global client instance
	apiClient *client.Client

	// Subscription validator for pre-execution validation
	subscriptionValidator *subscription.Validator
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     branding.CLIName,
	Version: Version, // Enables --version and -v flags
	Short:   branding.CLIShortDescription,
	Long:    branding.CLIDescription,
	// Run handles the root command when no subcommand is specified
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle --spec flag for root command
		if CheckSpecFlag() {
			format := GetOutputFormatWithDefault("json")
			return OutputSpec(cmd, format)
		}
		// If no --spec flag, show help
		return cmd.Help()
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip client initialization for non-API commands
		skipCommands := map[string]bool{
			"version":          true,
			"completion":       true,
			"__complete":       true, // Cobra's shell completion handler
			"__completeNoDesc": true, // Cobra's shell completion handler (no descriptions)
			"help":             true,
			branding.CLIName:   true, // Root command itself
		}
		if skipCommands[cmd.Name()] {
			return nil
		}

		// Check if API URL is configured
		if serverURL == "" {
			return fmt.Errorf("F5 Distributed Cloud API URL is not configured.\n\n"+
				"Please set the API URL using one of the following methods:\n"+
				"  1. Environment variable: export F5XC_API_URL=\"https://tenant.console.ves.volterra.io\"\n"+
				"  2. Command-line flag:    --server-url \"https://tenant.console.ves.volterra.io\"\n"+
				"  3. Configuration file:   Add 'server_url' to ~/%s\n\n"+
				"Replace 'tenant' with your actual F5 XC tenant name.\n"+
				"For staging environment, use: https://tenant.staging.volterra.us",
				branding.ConfigFileName)
		}

		// Initialize the API client
		cfg := &client.Config{
			ServerURL: serverURL,
			Cert:      cert,
			Key:       key,
			CACert:    cacert,
			P12Bundle: p12Bundle,
			Debug:     debug,
			Timeout:   timeout,
		}

		// Handle API token authentication
		// Check environment variable first, then fall back to flag
		if token := os.Getenv("F5XC_API_TOKEN"); token != "" {
			cfg.APIToken = token
		} else if useAPIToken {
			return fmt.Errorf("F5XC_API_TOKEN environment variable not set")
		}

		var err error
		apiClient, err = client.New(cfg)
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		// Initialize subscription tier detection and caching
		// This runs in background and doesn't block command execution
		if apiClient != nil && !subscription.IsTierCached() {
			go initSubscriptionContext(cmd.Context())
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	// Configure help system after all commands are registered
	initHelpSystem()
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	pf := rootCmd.PersistentFlags()

	// Connection settings
	pf.StringVarP(&cacert, "cacert", "a", "", "Path to the server CA certificate file for TLS verification.")
	pf.StringVarP(&cert, "cert", "c", "", "Path to the client certificate file for mTLS authentication.")
	pf.StringVar(&cfgFile, "config", "", "Path to configuration file.")
	pf.BoolVar(&hardwareKey, "hardwareKey", false, "Use a YubiKey hardware security module for TLS authentication.")
	pf.StringVarP(&key, "key", "k", "", "Path to the client private key file for mTLS authentication.")

	// Output format
	pf.StringVar(&outputFormat, "output-format", "", "Set the output format to text, json, yaml, or table.")

	pf.StringVarP(&outputDir, "output", "o", "./", "Directory path for command output files.")
	pf.StringVar(&p12Bundle, "p12-bundle", "", "Path to PKCS#12 bundle file containing client certificate and key.")
	pf.StringVarP(&serverURL, "server-url", "u", "", "F5 Distributed Cloud API endpoint URL.")
	pf.BoolVar(&showCurl, "show-curl", false, "Output equivalent curl commands for each API request.")
	pf.IntVar(&timeout, "timeout", 5, "Maximum time in seconds to wait for command completion.")
	pf.BoolVar(&useAPIToken, "api-token", false, "Use API token authentication.")
	pf.BoolVar(&nonInteractive, "non-interactive", false, "Disable interactive prompts and fail if required arguments are missing.")

	// Bind flags to viper (errors are ignored as flags are guaranteed to exist)
	_ = viper.BindPFlag("server-url", pf.Lookup("server-url"))
	_ = viper.BindPFlag("cert", pf.Lookup("cert"))
	_ = viper.BindPFlag("key", pf.Lookup("key"))
	_ = viper.BindPFlag("cacert", pf.Lookup("cacert"))
	_ = viper.BindPFlag("p12-bundle", pf.Lookup("p12-bundle"))
	_ = viper.BindPFlag("config", pf.Lookup("config"))
	_ = viper.BindPFlag("output-format", pf.Lookup("output-format"))

	// Bind environment variables to viper for flags without automatic binding
	_ = viper.BindEnv("config", "F5XC_CONFIG")
	_ = viper.BindEnv("cert", "F5XC_CERT")
	_ = viper.BindEnv("key", "F5XC_KEY")
	_ = viper.BindEnv("cacert", "F5XC_CACERT")
	_ = viper.BindEnv("p12-bundle", "F5XC_P12_FILE")
	_ = viper.BindEnv("output-format", "F5XC_OUTPUT")
	_ = viper.BindEnv("server-url", "F5XC_API_URL")

	// Register --spec flag for machine-readable CLI specification
	RegisterSpecFlag(rootCmd)

	// Set custom version template for clean one-liner output (--version, -v)
	rootCmd.SetVersionTemplate("xcsh version {{.Version}}\n")

	// Set custom help template with Environment Variables section
	rootCmd.SetHelpTemplate(helpTemplateWithEnvVars())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Check F5XC_CONFIG environment variable if cfgFile not set via CLI flag
	if cfgFile == "" {
		if envConfig := os.Getenv("F5XC_CONFIG"); envConfig != "" {
			cfgFile = envConfig
		}
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		viper.SetConfigType("yaml") // .f5xcconfig files are YAML
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			if debug {
				fmt.Fprintln(os.Stderr, "Warning: could not find home directory:", err)
			}
			return
		}

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".f5xcconfig")
	}

	viper.SetEnvPrefix("F5XC")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if debug {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
		applyConfigToFlags()
	} else {
		// No config file found, check environment variables or apply default
		if serverURL == "" {
			// Try F5XC_API_URL environment variable first
			if apiURL := os.Getenv("F5XC_API_URL"); apiURL != "" {
				normalized, err := client.NormalizeAPIURL(apiURL)
				if err != nil {
					// If normalization fails, use the raw URL (backward compatibility)
					if debug {
						fmt.Fprintf(os.Stderr, "Warning: Failed to normalize F5XC_API_URL: %v\n", err)
					}
					serverURL = apiURL
				} else {
					serverURL = normalized
				}
			}
		}
	}
}

// applyConfigToFlags applies viper config values to flags
// Precedence order: CLI flags > Environment variables > Config file > Defaults
func applyConfigToFlags() {
	cfg, err := config.Load(viper.ConfigFileUsed())
	if err != nil {
		// If config file couldn't be loaded, still check environment variables
		applyEnvironmentVariables()
		return
	}

	// Apply config file values first (lowest precedence after defaults)
	if serverURL == "" && cfg.ServerURL != "" {
		normalized, err := client.NormalizeAPIURL(cfg.ServerURL)
		if err != nil {
			// If normalization fails, use the raw URL (backward compatibility)
			if debug {
				fmt.Fprintf(os.Stderr, "Warning: Failed to normalize server URL from config: %v\n", err)
			}
			serverURL = cfg.ServerURL
		} else {
			serverURL = normalized
		}
	}
	if cert == "" && cfg.Cert != "" {
		cert = expandPath(cfg.Cert)
	}
	if key == "" && cfg.Key != "" {
		key = expandPath(cfg.Key)
	}
	if p12Bundle == "" && cfg.P12Bundle != "" {
		p12Bundle = expandPath(cfg.P12Bundle)
	}
	if !useAPIToken && cfg.APIToken {
		useAPIToken = true
	}

	// Apply environment variables (higher precedence than config file)
	applyEnvironmentVariables()
}

// applyEnvironmentVariables applies F5XC_* environment variables to flags
// This is called after config file values are applied, allowing env vars to override
func applyEnvironmentVariables() {
	// F5XC_API_URL overrides server-url (with normalization)
	if envURL := os.Getenv("F5XC_API_URL"); envURL != "" {
		normalized, err := client.NormalizeAPIURL(envURL)
		if err != nil {
			// If normalization fails, use the raw URL (backward compatibility)
			if debug {
				fmt.Fprintf(os.Stderr, "Warning: Failed to normalize F5XC_API_URL: %v\n", err)
			}
			serverURL = envURL
		} else {
			serverURL = normalized
		}
	}

	// F5XC_CERT overrides cert (only if CLI flag not set)
	if cert == "" {
		if envCert := os.Getenv("F5XC_CERT"); envCert != "" {
			cert = expandPath(envCert)
		}
	}

	// F5XC_KEY overrides key (only if CLI flag not set)
	if key == "" {
		if envKey := os.Getenv("F5XC_KEY"); envKey != "" {
			key = expandPath(envKey)
		}
	}

	// F5XC_CACERT overrides cacert (only if CLI flag not set)
	if cacert == "" {
		if envCACert := os.Getenv("F5XC_CACERT"); envCACert != "" {
			cacert = expandPath(envCACert)
		}
	}

	// F5XC_P12_FILE overrides p12-bundle (only if CLI flag not set)
	if p12Bundle == "" {
		if envP12 := os.Getenv("F5XC_P12_FILE"); envP12 != "" {
			p12Bundle = expandPath(envP12)
		}
	}

	// F5XC_OUTPUT overrides output-format (only if CLI flag not set)
	if outputFormat == "" {
		if envOutput := os.Getenv("F5XC_OUTPUT"); envOutput != "" {
			outputFormat = envOutput
		}
	}

	// F5XC_API_TOKEN automatically enables API token authentication
	if !useAPIToken {
		if envToken := os.Getenv("F5XC_API_TOKEN"); envToken != "" {
			useAPIToken = true
		}
	}
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[1:])
	}
	return path
}

// GetClient returns the initialized API client
func GetClient() *client.Client {
	return apiClient
}

// GetOutputFormat returns the current output format (defaults to table for list operations)
func GetOutputFormat() string {
	if outputFormat != "" {
		return outputFormat
	}
	return "table" // Default to table for list operations
}

// GetOutputFormatWithDefault returns the current output format with a custom default
func GetOutputFormatWithDefault(defaultFmt string) string {
	if outputFormat != "" {
		return outputFormat
	}
	return defaultFmt
}

// IsNonInteractive returns whether non-interactive mode is enabled
func IsNonInteractive() bool {
	return nonInteractive
}

// GetOutputDir returns the output directory
func GetOutputDir() string {
	return outputDir
}

// IsDebug returns whether debug mode is enabled
func IsDebug() bool {
	return debug
}

// GetTimeout returns the timeout in seconds
func GetTimeout() int {
	return timeout
}

// ShowCurl returns whether to emit CURL format
func ShowCurl() bool {
	return showCurl
}

// initSubscriptionContext detects and caches the subscription tier.
// This is called once at startup to cache the tier for subsequent commands.
func initSubscriptionContext(ctx context.Context) {
	if apiClient == nil {
		return
	}

	subClient := subscription.NewClient(apiClient)
	subscriptionValidator = subscription.NewValidator(subClient)

	// Detect and cache tier (errors are logged but don't fail the command)
	_, err := subscriptionValidator.DetectAndCacheTier(ctx)
	if err != nil && debug {
		fmt.Fprintf(os.Stderr, "Warning: Failed to detect subscription tier: %v\n", err)
	}
}

// GetSubscriptionValidator returns the subscription validator.
// If not initialized, creates one with the current API client.
func GetSubscriptionValidator() *subscription.Validator {
	if subscriptionValidator == nil && apiClient != nil {
		subClient := subscription.NewClient(apiClient)
		subscriptionValidator = subscription.NewValidator(subClient)
	}
	return subscriptionValidator
}

// EnsureSubscriptionTier ensures the subscription tier is cached.
// This should be called before commands that need tier information synchronously.
func EnsureSubscriptionTier(ctx context.Context) (string, error) {
	validator := GetSubscriptionValidator()
	if validator == nil {
		// No API client, return default
		return "Standard", nil
	}
	return validator.GetCurrentTier(ctx)
}

// helpTemplateWithEnvVars returns a custom help template that includes environment variables
func helpTemplateWithEnvVars() string {
	// Build environment variables section with consistent column alignment
	envVarsSection := "\nEnvironment Variables:\n"

	// Find max name length for alignment
	maxLen := 0
	for _, env := range EnvVarRegistry {
		if len(env.Name) > maxLen {
			maxLen = len(env.Name)
		}
	}

	// Format each env var on single line, matching flag style
	for _, env := range EnvVarRegistry {
		padding := maxLen - len(env.Name) + 3
		if env.RelatedFlag != "" {
			envVarsSection += fmt.Sprintf("  %s%s%s [%s]\n", env.Name, spaces(padding), env.Description, env.RelatedFlag)
		} else {
			envVarsSection += fmt.Sprintf("  %s%s%s\n", env.Name, spaces(padding), env.Description)
		}
	}

	// Add examples section
	examplesSection := `
Examples:
  xcsh identity list namespace                    List all namespaces
  xcsh load_balancer get http_loadbalancer -n shared   Get a specific resource
  xcsh request /api/web/namespaces                     Execute custom API request
  xcsh --spec --output-format json                     Output CLI spec for automation
`

	// Add configuration precedence section
	configSection := `
Configuration:
  Config file:  ~/.f5xcconfig
  Priority:     CLI flags > environment variables > config file > defaults

Learn more:    https://robinmordasiewicz.github.io/xcsh/
`

	// Custom template based on Cobra's default, with additional sections
	return `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}` + envVarsSection + examplesSection + configSection
}

// spaces returns a string of n spaces for alignment
func spaces(n int) string {
	if n <= 0 {
		return " "
	}
	s := ""
	for i := 0; i < n; i++ {
		s += " "
	}
	return s
}

// helpSystemInitialized tracks whether the help system has been set up
var helpSystemInitialized bool

// initHelpSystem configures help command visibility and templates.
// This runs after all commands are registered but before Execute() runs.
func initHelpSystem() {
	if helpSystemInitialized {
		return
	}
	helpSystemInitialized = true

	// Set a hidden help command - users should use --help flag instead
	// The help command still works if typed directly, just hidden from listing
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "help [command]",
		Short:  "Help about any command",
		Hidden: true,
		Run: func(c *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = rootCmd.Help()
				return
			}
			cmd, _, err := rootCmd.Find(args)
			if cmd == nil || err != nil {
				c.Printf("Unknown help topic %#q\n", args)
				_ = rootCmd.Usage()
				return
			}
			_ = cmd.Help()
		},
	})

	// Apply custom usage template that properly filters hidden commands
	applyUsageTemplateRecursively(rootCmd, usageTemplateWithHiddenFilter())

	// Apply custom help template to all commands recursively
	applyHelpTemplateRecursively(rootCmd, helpTemplateWithEnvVars())
}

// usageTemplateWithHiddenFilter returns a usage template that properly filters hidden commands
func usageTemplateWithHiddenFilter() string {
	return `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

Available Commands:{{range $cmds}}{{if (and .IsAvailableCommand (not .Hidden))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}

// applyUsageTemplateRecursively applies custom usage template to all commands
func applyUsageTemplateRecursively(cmd *cobra.Command, template string) {
	cmd.SetUsageTemplate(template)
	for _, subCmd := range cmd.Commands() {
		applyUsageTemplateRecursively(subCmd, template)
	}
}

// applyHelpTemplateRecursively applies custom help template to all commands
// in the command tree. This ensures consistent help formatting with environment
// variables section across all subcommands.
func applyHelpTemplateRecursively(cmd *cobra.Command, template string) {
	cmd.SetHelpTemplate(template)
	for _, subCmd := range cmd.Commands() {
		applyHelpTemplateRecursively(subCmd, template)
	}
}

// ValidateDomainTier checks if the current subscription tier is sufficient for the domain.
// Returns nil if tier is sufficient, or a TierAccessError if access is denied.
// Falls back to allowing access if tier cannot be determined (offline mode).
func ValidateDomainTier(ctx context.Context, domain string) error {
	// Get domain info
	info, found := types.GetDomainInfo(domain)
	if !found {
		return fmt.Errorf("domain not found: %s", domain)
	}

	// Get current subscription tier
	currentTier, err := EnsureSubscriptionTier(ctx)
	if err != nil {
		// If we can't determine tier, allow access (offline/fallback mode)
		// User will see error from API when trying to actually access restricted domain
		return nil
	}

	// Check if tier is sufficient
	if !validation.IsSufficientTier(currentTier, info.RequiresTier) {
		// Create and return a tier access error with user-friendly message
		tierErr := validation.NewTierAccessError(
			domain,
			info.DisplayName,
			currentTier,
			info.RequiresTier,
		)
		return tierErr
	}

	return nil
}

// GetCurrentTierForDomain returns the current subscription tier for the user.
// Falls back to "Standard" if tier cannot be determined.
func GetCurrentTierForDomain(ctx context.Context) string {
	tier, err := EnsureSubscriptionTier(ctx)
	if err != nil {
		return validation.TierStandard // Default to Standard tier
	}
	return tier
}

// CheckAndWarnPreviewDomain checks if a domain is in preview and displays a warning if it is.
// Returns a warning error if the domain is in preview, nil if the domain is stable.
func CheckAndWarnPreviewDomain(domain string) *validation.PreviewWarning {
	// Get domain info
	info, found := types.GetDomainInfo(domain)
	if !found {
		return nil // Domain not found, skip preview check
	}

	// Check if domain is in preview
	if !info.IsPreview {
		return nil // Domain is stable, no warning needed
	}

	// Return preview warning
	return validation.GetPreviewWarning(domain, info.DisplayName)
}
