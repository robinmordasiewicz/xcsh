package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/robinmordasiewicz/f5xcctl/pkg/client"
	"github.com/robinmordasiewicz/f5xcctl/pkg/config"
)

var (
	// Config file path
	cfgFile string

	// Connection settings (f5xcctl compatible)
	serverURL   string
	cert        string
	key         string
	cacert      string
	p12Bundle   string
	hardwareKey bool // Use yubikey for TLS connection
	useAPIToken bool // Use API token from F5XC_API_TOKEN environment variable

	// Output control (f5xcctl compatible)
	outputFormat string // Output format for command (canonical: --output-format)
	outputDir    string // Output dir for command

	// Behavior flags (f5xcctl compatible)
	showCurl       bool // Emit requests from program in CURL format
	timeout        int  // Timeout (in seconds) for command to finish
	nonInteractive bool // Fail on missing arguments instead of prompting

	// Internal flags (not exposed to CLI)
	debug bool

	// Global client instance
	apiClient *client.Client
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "f5xcctl",
	Short: "Command-line interface for F5 Distributed Cloud services.",
	Long:  `Command-line interface for F5 Distributed Cloud services.`,
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
			"f5xcctl":          true, // Root command itself
		}
		if skipCommands[cmd.Name()] {
			return nil
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
		if useAPIToken {
			token := os.Getenv("F5XC_API_TOKEN")
			if token == "" {
				return fmt.Errorf("F5XC_API_TOKEN environment variable not set")
			}
			cfg.APIToken = token
		}

		var err error
		apiClient, err = client.New(cfg)
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags matching original f5xcctl exactly
	pf := rootCmd.PersistentFlags()

	// Connection settings (f5xcctl compatible)
	pf.StringVarP(&cacert, "cacert", "a", "", "Path to the server CA certificate file for TLS verification.")
	pf.StringVarP(&cert, "cert", "c", "", "Path to the client certificate file for mTLS authentication.")
	pf.StringVar(&cfgFile, "config", "", "Path to configuration file.")
	pf.BoolVar(&hardwareKey, "hardwareKey", false, "Use a YubiKey hardware security module for TLS authentication.")
	pf.StringVarP(&key, "key", "k", "", "Path to the client private key file for mTLS authentication.")

	// Output format: --output-format is canonical, --outfmt is hidden alias for backward compatibility
	pf.StringVar(&outputFormat, "output-format", "", "Set the output format to text, json, yaml, or table.")
	pf.StringVar(&outputFormat, "outfmt", "", "Output format for command (deprecated: use --output-format).")
	_ = pf.MarkHidden("outfmt") // Hide deprecated alias

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

	// Register --spec flag for machine-readable CLI specification
	RegisterSpecFlag(rootCmd)

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

	viper.SetEnvPrefix("VES")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if debug {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
		applyConfigToFlags()
	} else {
		// No config file found, apply default
		if serverURL == "" {
			serverURL = "http://localhost:8001"
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
		// Apply fallback default if still not set
		if serverURL == "" {
			serverURL = "http://localhost:8001"
		}
		return
	}

	// Apply config file values first (lowest precedence after defaults)
	if serverURL == "" && cfg.ServerURL != "" {
		serverURL = cfg.ServerURL
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

	// Apply fallback default if still not set
	if serverURL == "" {
		serverURL = "http://localhost:8001"
	}
}

// applyEnvironmentVariables applies F5XC_* environment variables to flags
// This is called after config file values are applied, allowing env vars to override
func applyEnvironmentVariables() {
	// F5XC_API_URL overrides server-url
	if envURL := os.Getenv("F5XC_API_URL"); envURL != "" {
		serverURL = envURL
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
  f5xcctl configuration list namespace                    List all namespaces
  f5xcctl configuration get http_loadbalancer -n shared   Get a specific resource
  f5xcctl request /api/web/namespaces                     Execute custom API request
  f5xcctl --spec --output-format json                     Output CLI spec for automation
`

	// Add configuration precedence section
	configSection := `
Configuration:
  Config file:  ~/.f5xcconfig
  Priority:     CLI flags > environment variables > config file > defaults

Learn more:    https://robinmordasiewicz.github.io/f5xcctl/
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
