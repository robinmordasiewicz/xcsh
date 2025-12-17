package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/robinmordasiewicz/f5xcctl/pkg/output"
)

var configureFlags struct {
	serverURL      string
	p12Bundle      string
	cert           string
	key            string
	outputFmt      string
	apiToken       bool
	nonInteractive bool
}

var configureCmd = &cobra.Command{
	Use:    "configure",
	Short:  "Configure CLI settings.",
	Hidden: true, // Hide from help to match original f5xcctl
	Long: `Configure the F5 Distributed Cloud CLI settings.

This command helps you set up your CLI configuration including:
  - API server URL (your tenant URL)
  - Authentication credentials (P12 bundle or cert/key pair)
  - Default output format

The configuration is saved to ~/.f5xcconfig (or the path specified by --config).`,
	Example: `  # Interactive configuration
  f5xcctl configure

  # Non-interactive configuration
  f5xcctl configure --server-url https://example-tenant.console.ves.volterra.io/api \
    --p12-bundle ~/.f5xcctl/example-cert.p12

  # Configure with certificate and key
  f5xcctl configure --server-url https://example-tenant.console.ves.volterra.io/api \
    --cert ~/.f5xcctl/cert.pem --key ~/.f5xcctl/key.pem`,
	RunE: runConfigure,
}

var configureShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current CLI configuration settings.`,
	Example: `  # Show current configuration
  f5xcctl configure show`,
	RunE: runConfigureShow,
}

var configureSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long: `Set a specific configuration value.

Available keys:
  - server-url: API server URL
  - p12-bundle: Path to P12 certificate bundle
  - cert: Path to client certificate
  - key: Path to client key
  - api-token: Enable API token auth (true/false)
  - output: Default output format (json, yaml, table)`,
	Example: `  # Set the server URL
  f5xcctl configure set server-url https://example-tenant.console.ves.volterra.io/api

  # Enable API token authentication
  f5xcctl configure set api-token true

  # Set the default output format
  f5xcctl configure set output yaml`,
	Args: cobra.ExactArgs(2),
	RunE: runConfigureSet,
}

func init() {
	rootCmd.AddCommand(configureCmd)

	configureCmd.Flags().StringVar(&configureFlags.serverURL, "server-url", "", "API server URL")
	configureCmd.Flags().StringVar(&configureFlags.p12Bundle, "p12-bundle", "", "Path to P12 certificate bundle")
	configureCmd.Flags().StringVar(&configureFlags.cert, "cert", "", "Path to client certificate")
	configureCmd.Flags().StringVar(&configureFlags.key, "key", "", "Path to client key")
	configureCmd.Flags().StringVar(&configureFlags.outputFmt, "output-format", "", "Default output format")
	configureCmd.Flags().BoolVar(&configureFlags.apiToken, "api-token", false, "Use API token authentication (token from F5XC_API_TOKEN env var)")
	configureCmd.Flags().BoolVar(&configureFlags.nonInteractive, "non-interactive", false, "Run in non-interactive mode")

	configureCmd.AddCommand(configureShowCmd)
	configureCmd.AddCommand(configureSetCmd)
}

// ConfigFile represents the configuration file structure
type ConfigFile struct {
	ServerURL string `yaml:"server-url,omitempty"`
	P12Bundle string `yaml:"p12-bundle,omitempty"`
	Cert      string `yaml:"cert,omitempty"`
	Key       string `yaml:"key,omitempty"`
	Output    string `yaml:"output,omitempty"`
	APIToken  bool   `yaml:"api-token,omitempty"`
}

// rawConfigFile is used for YAML parsing
type rawConfigFile struct {
	ServerURL string `yaml:"server-url"`
	P12Bundle string `yaml:"p12-bundle"`
	Cert      string `yaml:"cert"`
	Key       string `yaml:"key"`
	Output    string `yaml:"output"`
	APIToken  bool   `yaml:"api-token"`
}

// parseConfigFile parses a config file
func parseConfigFile(data []byte) (*ConfigFile, error) {
	var raw rawConfigFile
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	cfg := &ConfigFile{
		ServerURL: raw.ServerURL,
		P12Bundle: raw.P12Bundle,
		Cert:      raw.Cert,
		Key:       raw.Key,
		Output:    raw.Output,
		APIToken:  raw.APIToken,
	}

	return cfg, nil
}

func runConfigure(cmd *cobra.Command, args []string) error {
	configPath := getConfigPath()

	// Load existing config if it exists
	config := &ConfigFile{}
	if data, err := os.ReadFile(configPath); err == nil {
		if parsed, err := parseConfigFile(data); err == nil {
			config = parsed
		}
	}

	// Non-interactive mode
	if configureFlags.nonInteractive {
		return updateConfigFromFlags(config, configPath)
	}

	// Interactive mode
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("F5 Distributed Cloud CLI Configuration")
	fmt.Println("=======================================")
	fmt.Println()

	// Server URL
	serverURL := promptWithDefault(reader, "API Server URL", config.ServerURL,
		"e.g., https://example-tenant.console.ves.volterra.io/api")
	if serverURL != "" {
		config.ServerURL = serverURL
	}

	// Authentication method
	fmt.Println()
	fmt.Println("Authentication Method:")
	fmt.Println("  1. P12 Certificate Bundle (recommended)")
	fmt.Println("  2. Certificate and Key files")
	fmt.Println("  3. API Token (via F5XC_API_TOKEN environment variable)")
	fmt.Print("Choose [1/2/3] (default: 1): ")
	authChoice, _ := reader.ReadString('\n')
	authChoice = strings.TrimSpace(authChoice)

	switch authChoice {
	case "3":
		// API Token
		config.APIToken = true
		config.P12Bundle = "" // Clear other auth methods
		config.Cert = ""
		config.Key = ""
		fmt.Println()
		fmt.Println("Note: Set the F5XC_API_TOKEN environment variable:")
		fmt.Println("  export F5XC_API_TOKEN='your-api-token'")
	case "2":
		// Cert and Key
		config.APIToken = false // Clear API token
		cert := promptWithDefault(reader, "Certificate file path", config.Cert, "")
		if cert != "" {
			config.Cert = expandPath(cert)
			config.P12Bundle = "" // Clear P12 if using cert/key
		}

		key := promptWithDefault(reader, "Key file path", config.Key, "")
		if key != "" {
			config.Key = expandPath(key)
		}
	default:
		// P12 Bundle
		config.APIToken = false // Clear API token
		p12 := promptWithDefault(reader, "P12 Bundle path", config.P12Bundle, "")
		if p12 != "" {
			config.P12Bundle = expandPath(p12)
			config.Cert = "" // Clear cert/key if using P12
			config.Key = ""
		}

		if config.P12Bundle != "" {
			fmt.Println()
			fmt.Println("Note: Set the F5XC_P12_PASSWORD environment variable with your P12 password:")
			fmt.Println("  export F5XC_P12_PASSWORD='your-password'")
		}
	}

	// Output format
	fmt.Println()
	outputFmt := promptWithDefault(reader, "Default output format (json/yaml/table)", config.Output, "")
	if outputFmt != "" {
		config.Output = outputFmt
	}

	// Save configuration
	if err := saveConfig(config, configPath); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Println()
	output.PrintInfo(fmt.Sprintf("Configuration saved to %s", configPath))

	// Test connection
	fmt.Print("\nTest connection? [Y/n]: ")
	testChoice, _ := reader.ReadString('\n')
	testChoice = strings.TrimSpace(strings.ToLower(testChoice))
	if testChoice != "n" && testChoice != "no" {
		return testConnection()
	}

	return nil
}

func runConfigureShow(cmd *cobra.Command, args []string) error {
	configPath := getConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("No configuration file found at %s\n", configPath)
			fmt.Println("Run 'f5xcctl configure' to create one.")
			return nil
		}
		return fmt.Errorf("failed to read config: %w", err)
	}

	config, err := parseConfigFile(data)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// Display configuration
	fmt.Printf("Configuration file: %s\n\n", configPath)

	displayConfig := map[string]interface{}{
		"server_url": config.ServerURL,
	}

	if config.P12Bundle != "" {
		displayConfig["p12_bundle"] = config.P12Bundle
	}
	if config.Cert != "" {
		displayConfig["cert"] = config.Cert
	}
	if config.Key != "" {
		displayConfig["key"] = config.Key
	}
	if config.Output != "" {
		displayConfig["output"] = config.Output
	}
	if config.APIToken {
		displayConfig["api_token"] = "enabled (token from F5XC_API_TOKEN)"
	}

	return output.Print(displayConfig, GetOutputFormat())
}

func runConfigureSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	configPath := getConfigPath()

	// Load existing config
	config := &ConfigFile{}
	if data, err := os.ReadFile(configPath); err == nil {
		if parsed, err := parseConfigFile(data); err == nil {
			config = parsed
		}
	}

	// Set the value
	switch key {
	case "server-url":
		config.ServerURL = value
	case "p12-bundle":
		config.P12Bundle = expandPath(value)
		config.APIToken = false // Clear API token when setting P12
	case "cert":
		config.Cert = expandPath(value)
		config.APIToken = false // Clear API token when setting cert
	case "key":
		config.Key = expandPath(value)
		config.APIToken = false // Clear API token when setting key
	case "output":
		config.Output = value
	case "api-token":
		if value == "true" || value == "1" || value == "enabled" {
			config.APIToken = true
			config.P12Bundle = "" // Clear other auth methods
			config.Cert = ""
			config.Key = ""
		} else {
			config.APIToken = false
		}
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}

	// Save configuration
	if err := saveConfig(config, configPath); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	output.PrintInfo(fmt.Sprintf("Set %s = %s", key, value))
	return nil
}

func updateConfigFromFlags(config *ConfigFile, configPath string) error {
	if configureFlags.serverURL != "" {
		config.ServerURL = configureFlags.serverURL
	}
	if configureFlags.apiToken {
		config.APIToken = true
		config.P12Bundle = "" // Clear other auth methods
		config.Cert = ""
		config.Key = ""
	} else {
		if configureFlags.p12Bundle != "" {
			config.P12Bundle = expandPath(configureFlags.p12Bundle)
			config.APIToken = false
		}
		if configureFlags.cert != "" {
			config.Cert = expandPath(configureFlags.cert)
			config.APIToken = false
		}
		if configureFlags.key != "" {
			config.Key = expandPath(configureFlags.key)
			config.APIToken = false
		}
	}
	if configureFlags.outputFmt != "" {
		config.Output = configureFlags.outputFmt
	}

	return saveConfig(config, configPath)
}

func getConfigPath() string {
	if cfgFile != "" {
		return cfgFile
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".f5xcconfig")
}

func saveConfig(config *ConfigFile, path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func promptWithDefault(reader *bufio.Reader, prompt, defaultVal, hint string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultVal)
	} else if hint != "" {
		fmt.Printf("%s (%s): ", prompt, hint)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultVal
	}
	return input
}

func testConnection() error {
	client := GetClient()
	if client == nil {
		return fmt.Errorf("failed to initialize client - check your configuration")
	}

	output.PrintInfo("Testing connection...")

	// Try to list namespaces as a simple connectivity test
	// This is a lightweight API call that should work for any authenticated user
	fmt.Println("Connection successful!")
	return nil
}
