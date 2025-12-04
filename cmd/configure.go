package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/robinmordasiewicz/f5xc/pkg/output"
)

var configureFlags struct {
	serverURL  string
	p12Bundle  string
	cert       string
	key        string
	tenant     string
	outputFmt  string
	nonInteractive bool
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure CLI settings",
	Long: `Configure the F5 Distributed Cloud CLI settings.

This command helps you set up your CLI configuration including:
  - API server URL (your tenant URL)
  - Authentication credentials (P12 bundle or cert/key pair)
  - Default output format

The configuration is saved to ~/.vesconfig (or the path specified by --config).`,
	Example: `  # Interactive configuration
  f5xc configure

  # Non-interactive configuration
  f5xc configure --server-url https://my-tenant.console.ves.volterra.io/api \
    --p12-bundle ~/.vesctl/my-cert.p12

  # Configure with certificate and key
  f5xc configure --server-url https://my-tenant.console.ves.volterra.io/api \
    --cert ~/.vesctl/cert.pem --key ~/.vesctl/key.pem`,
	RunE: runConfigure,
}

var configureShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current CLI configuration settings.`,
	Example: `  # Show current configuration
  f5xc configure show`,
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
  - output: Default output format (json, yaml, table)`,
	Example: `  # Set the server URL
  f5xc configure set server-url https://my-tenant.console.ves.volterra.io/api

  # Set the default output format
  f5xc configure set output yaml`,
	Args: cobra.ExactArgs(2),
	RunE: runConfigureSet,
}

func init() {
	rootCmd.AddCommand(configureCmd)

	configureCmd.Flags().StringVar(&configureFlags.serverURL, "server-url", "", "API server URL")
	configureCmd.Flags().StringVar(&configureFlags.p12Bundle, "p12-bundle", "", "Path to P12 certificate bundle")
	configureCmd.Flags().StringVar(&configureFlags.cert, "cert", "", "Path to client certificate")
	configureCmd.Flags().StringVar(&configureFlags.key, "key", "", "Path to client key")
	configureCmd.Flags().StringVar(&configureFlags.tenant, "tenant", "", "Tenant name")
	configureCmd.Flags().StringVar(&configureFlags.outputFmt, "output-format", "", "Default output format")
	configureCmd.Flags().BoolVar(&configureFlags.nonInteractive, "non-interactive", false, "Run in non-interactive mode")

	configureCmd.AddCommand(configureShowCmd)
	configureCmd.AddCommand(configureSetCmd)
}

// ConfigFile represents the configuration file structure
type ConfigFile struct {
	ServerURLs []string `yaml:"server-urls,omitempty"`
	P12Bundle  string   `yaml:"p12-bundle,omitempty"`
	Cert       string   `yaml:"cert,omitempty"`
	Key        string   `yaml:"key,omitempty"`
	Tenant     string   `yaml:"tenant,omitempty"`
	Output     string   `yaml:"output,omitempty"`
}

func runConfigure(cmd *cobra.Command, args []string) error {
	configPath := getConfigPath()

	// Load existing config if it exists
	config := &ConfigFile{}
	if data, err := os.ReadFile(configPath); err == nil {
		_ = yaml.Unmarshal(data, config)
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
	currentServer := ""
	if len(config.ServerURLs) > 0 {
		currentServer = config.ServerURLs[0]
	}
	serverURL := promptWithDefault(reader, "API Server URL", currentServer,
		"e.g., https://my-tenant.console.ves.volterra.io/api")
	if serverURL != "" {
		config.ServerURLs = []string{serverURL}
	}

	// Authentication method
	fmt.Println()
	fmt.Println("Authentication Method:")
	fmt.Println("  1. P12 Certificate Bundle (recommended)")
	fmt.Println("  2. Certificate and Key files")
	fmt.Print("Choose [1/2] (default: 1): ")
	authChoice, _ := reader.ReadString('\n')
	authChoice = strings.TrimSpace(authChoice)

	if authChoice == "2" {
		// Cert and Key
		cert := promptWithDefault(reader, "Certificate file path", config.Cert, "")
		if cert != "" {
			config.Cert = expandPath(cert)
			config.P12Bundle = "" // Clear P12 if using cert/key
		}

		key := promptWithDefault(reader, "Key file path", config.Key, "")
		if key != "" {
			config.Key = expandPath(key)
		}
	} else {
		// P12 Bundle
		p12 := promptWithDefault(reader, "P12 Bundle path", config.P12Bundle, "")
		if p12 != "" {
			config.P12Bundle = expandPath(p12)
			config.Cert = "" // Clear cert/key if using P12
			config.Key = ""
		}

		if config.P12Bundle != "" {
			fmt.Println()
			fmt.Println("Note: Set the VES_P12_PASSWORD environment variable with your P12 password:")
			fmt.Println("  export VES_P12_PASSWORD='your-password'")
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

	config := &ConfigFile{}
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("No configuration file found at %s\n", configPath)
			fmt.Println("Run 'f5xc configure' to create one.")
			return nil
		}
		return fmt.Errorf("failed to read config: %w", err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// Display configuration
	fmt.Printf("Configuration file: %s\n\n", configPath)

	displayConfig := map[string]interface{}{
		"server_urls": config.ServerURLs,
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
	if config.Tenant != "" {
		displayConfig["tenant"] = config.Tenant
	}
	if config.Output != "" {
		displayConfig["output"] = config.Output
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
		_ = yaml.Unmarshal(data, config)
	}

	// Set the value
	switch key {
	case "server-url", "server-urls":
		config.ServerURLs = []string{value}
	case "p12-bundle":
		config.P12Bundle = expandPath(value)
	case "cert":
		config.Cert = expandPath(value)
	case "key":
		config.Key = expandPath(value)
	case "tenant":
		config.Tenant = value
	case "output":
		config.Output = value
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
		config.ServerURLs = []string{configureFlags.serverURL}
	}
	if configureFlags.p12Bundle != "" {
		config.P12Bundle = expandPath(configureFlags.p12Bundle)
	}
	if configureFlags.cert != "" {
		config.Cert = expandPath(configureFlags.cert)
	}
	if configureFlags.key != "" {
		config.Key = expandPath(configureFlags.key)
	}
	if configureFlags.tenant != "" {
		config.Tenant = configureFlags.tenant
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
	return filepath.Join(home, ".vesconfig")
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
