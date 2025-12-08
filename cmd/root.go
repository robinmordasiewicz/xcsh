package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/robinmordasiewicz/vesctl/pkg/client"
	"github.com/robinmordasiewicz/vesctl/pkg/config"
)

var (
	// Config file path
	cfgFile string

	// Connection settings (vesctl compatible)
	serverURLs  []string
	cert        string
	key         string
	cacert      string
	p12Bundle   string
	hardwareKey bool // Use yubikey for TLS connection

	// Output control (vesctl compatible)
	outfmt    string // Output format for command
	outputDir string // Output dir for command

	// Behavior flags (vesctl compatible)
	showCurl bool // Emit requests from program in CURL format
	timeout  int  // Timeout (in seconds) for command to finish

	// Internal flags (not exposed to CLI)
	debug bool

	// Global client instance
	apiClient *client.Client
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vesctl",
	Short: "A command line utility to interact with ves service.",
	Long:  `A command line utility to interact with ves service.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip client initialization for non-API commands
		skipCommands := map[string]bool{
			"version":    true,
			"completion": true,
			"help":       true,
		}
		if skipCommands[cmd.Name()] {
			return nil
		}

		// Initialize the API client
		cfg := &client.Config{
			ServerURLs: serverURLs,
			Cert:       cert,
			Key:        key,
			CACert:     cacert,
			P12Bundle:  p12Bundle,
			Debug:      debug,
			Timeout:    timeout,
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

	// Global flags matching original vesctl exactly
	pf := rootCmd.PersistentFlags()

	// Connection settings (vesctl compatible)
	pf.StringVarP(&cacert, "cacert", "a", "", "Server CA cert file path")
	pf.StringVarP(&cert, "cert", "c", "", "Client cert file path")
	// Get default config path for help text (matches original vesctl behavior)
	defaultConfigPath := "$HOME/.vesconfig"
	if home, err := os.UserHomeDir(); err == nil {
		defaultConfigPath = filepath.Join(home, ".vesconfig")
	}
	pf.StringVar(&cfgFile, "config", "", fmt.Sprintf("A configuration file to use for API gateway URL and credentials (default %q)", defaultConfigPath))
	pf.BoolVar(&hardwareKey, "hardwareKey", false, "Use yubikey for TLS connection")
	pf.StringVarP(&key, "key", "k", "", "Client key file path")
	pf.StringVar(&outfmt, "outfmt", "", "Output format for command")
	pf.StringVarP(&outputDir, "output", "o", "./", "Output dir for command")
	pf.StringVar(&p12Bundle, "p12-bundle", "", "Client P12 bundle (key+cert) file path. Any password for this file should be in environment variable VES_P12_PASSWORD")
	pf.StringSliceVarP(&serverURLs, "server-urls", "u", nil, "API endpoint URL (default [http://localhost:8001])")
	pf.BoolVar(&showCurl, "show-curl", false, "Emit requests from program in CURL format")
	pf.IntVar(&timeout, "timeout", 5, "Timeout (in seconds) for command to finish")

	// Bind flags to viper (errors are ignored as flags are guaranteed to exist)
	_ = viper.BindPFlag("server-urls", pf.Lookup("server-urls"))
	_ = viper.BindPFlag("cert", pf.Lookup("cert"))
	_ = viper.BindPFlag("key", pf.Lookup("key"))
	_ = viper.BindPFlag("cacert", pf.Lookup("cacert"))
	_ = viper.BindPFlag("p12-bundle", pf.Lookup("p12-bundle"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		viper.SetConfigType("yaml") // .vesconfig files are YAML
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
		viper.SetConfigName(".vesconfig")
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
		if len(serverURLs) == 0 {
			serverURLs = []string{"http://localhost:8001"}
		}
	}
}

// applyConfigToFlags applies viper config values to flags
func applyConfigToFlags() {
	cfg, err := config.Load(viper.ConfigFileUsed())
	if err != nil {
		// If config file couldn't be loaded, apply default
		if len(serverURLs) == 0 {
			serverURLs = []string{"http://localhost:8001"}
		}
		return
	}

	// Apply config values if CLI flags not set
	if len(serverURLs) == 0 && len(cfg.ServerURLs) > 0 {
		serverURLs = cfg.ServerURLs
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

	// Apply fallback default if still not set
	if len(serverURLs) == 0 {
		serverURLs = []string{"http://localhost:8001"}
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
	if outfmt != "" {
		return outfmt
	}
	return "table" // Default to table for list operations
}

// GetOutputFormatWithDefault returns the current output format with a custom default
func GetOutputFormatWithDefault(defaultFmt string) string {
	if outfmt != "" {
		return outfmt
	}
	return defaultFmt
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
