package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/robinmordasiewicz/f5xc/pkg/client"
	"github.com/robinmordasiewicz/f5xc/pkg/config"
)

var (
	// Config file path
	cfgFile string

	// Connection settings
	serverURLs []string
	cert       string
	key        string
	cacert     string
	p12Bundle  string

	// Output control
	outputFormat string
	query        string

	// Behavior flags
	debug      bool
	verbose    bool
	onlyErrors bool
	noWait     bool
	insecure   bool // Skip TLS certificate verification

	// Global client instance
	apiClient *client.Client
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vesctl",
	Short: "F5 Distributed Cloud CLI",
	Long: `F5 Distributed Cloud CLI (vesctl) is a command-line tool for managing
F5 Distributed Cloud (formerly Volterra) resources.

Getting Started:
  vesctl login                    Log in to F5 Distributed Cloud
  vesctl configure                Configure CLI settings
  vesctl --help                   Show help for any command

Common Commands:
  vesctl http-loadbalancer list   List HTTP load balancers
  vesctl origin-pool create       Create an origin pool
  vesctl namespace show           Show namespace details

Configuration:
  Default config file: ~/.vesconfig (YAML format)

  Example ~/.vesconfig:
    server-urls:
      - https://your-tenant.console.ves.volterra.io/api
    p12-bundle: ~/.vesctl/my-cert.p12

Authentication:
  - P12 bundle with VES_P12_PASSWORD environment variable
  - Certificate and key files (--cert and --key)

For more information, visit: https://docs.cloud.f5.com/`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip client initialization for non-API commands
		skipCommands := map[string]bool{
			"version":    true,
			"completion": true,
			"help":       true,
			"configure":  true,
		}
		if skipCommands[cmd.Name()] {
			return nil
		}

		// Initialize the API client
		cfg := &client.Config{
			ServerURLs:         serverURLs,
			Cert:               cert,
			Key:                key,
			CACert:             cacert,
			P12Bundle:          p12Bundle,
			Debug:              debug,
			InsecureSkipVerify: insecure,
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

	// Global flags following Azure CLI patterns
	pf := rootCmd.PersistentFlags()

	// Connection settings
	pf.StringVar(&cfgFile, "config", "", "Path to config file (default: ~/.vesconfig)")
	pf.StringSliceVarP(&serverURLs, "server-urls", "u", nil, "API server URL(s)")
	pf.StringVarP(&cert, "cert", "c", "", "Path to client certificate file")
	pf.StringVarP(&key, "key", "k", "", "Path to client key file")
	pf.StringVar(&cacert, "cacert", "", "Path to CA certificate file")
	pf.StringVar(&p12Bundle, "p12-bundle", "", "Path to P12 certificate bundle")

	// Output control (Azure CLI style)
	pf.StringVarP(&outputFormat, "output", "o", "", "Output format: json, yaml, table, tsv, none")
	pf.StringVar(&query, "query", "", "JMESPath query string for filtering output")

	// Behavior flags (Azure CLI style)
	pf.BoolVar(&debug, "debug", false, "Show all debug logs")
	pf.BoolVar(&verbose, "verbose", false, "Increase logging verbosity")
	pf.BoolVar(&onlyErrors, "only-show-errors", false, "Only show errors, suppressing warnings")
	pf.BoolVar(&noWait, "no-wait", false, "Do not wait for long-running operations to finish")
	pf.BoolVar(&insecure, "insecure", false, "Skip TLS certificate verification (use for staging/testing)")

	// Bind flags to viper (errors are ignored as flags are guaranteed to exist)
	_ = viper.BindPFlag("server-urls", pf.Lookup("server-urls"))
	_ = viper.BindPFlag("cert", pf.Lookup("cert"))
	_ = viper.BindPFlag("key", pf.Lookup("key"))
	_ = viper.BindPFlag("cacert", pf.Lookup("cacert"))
	_ = viper.BindPFlag("p12-bundle", pf.Lookup("p12-bundle"))
	_ = viper.BindPFlag("output", pf.Lookup("output"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
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
	}
}

// applyConfigToFlags applies viper config values to flags
func applyConfigToFlags() {
	cfg, err := config.Load(viper.ConfigFileUsed())
	if err != nil {
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

// GetOutputFormat returns the current output format
func GetOutputFormat() string {
	if outputFormat != "" {
		return outputFormat
	}
	return "yaml"
}

// GetQuery returns the JMESPath query string
func GetQuery() string {
	return query
}

// IsDebug returns whether debug mode is enabled
func IsDebug() bool {
	return debug
}

// IsVerbose returns whether verbose mode is enabled
func IsVerbose() bool {
	return verbose
}
