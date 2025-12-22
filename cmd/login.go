package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/robinmordasiewicz/xcsh/pkg/output"
)

var loginFlags struct {
	tenant    string
	p12Bundle string
	cert      string
	key       string
	apiToken  bool
}

var loginCmd = &cobra.Command{
	Use:    "login",
	Short:  "Log in to F5 Distributed Cloud.",
	Hidden: true, // Hide from help to match original f5xcctl
	Long: `Authenticate with F5 Distributed Cloud.

This command validates your credentials and saves them to the configuration file.
You can authenticate using:
  - A P12 certificate bundle (set F5XC_P12_PASSWORD environment variable)
  - Certificate and key files
  - API token (set F5XC_API_TOKEN environment variable)

After successful login, you can use all f5xcctl commands to manage your resources.`,
	Example: `  # Login with P12 bundle
  f5xcctl login --tenant example-tenant --p12-bundle ~/.f5xcctl/example-cert.p12

  # Login with certificate and key
  f5xcctl login --tenant example-tenant --cert ~/.f5xcctl/cert.pem --key ~/.f5xcctl/key.pem

  # Login with API token
  export F5XC_API_TOKEN='your-api-token'
  f5xcctl login --tenant example-tenant --api-token

  # Login (using existing configuration)
  f5xcctl login`,
	RunE: runLogin,
}

var logoutCmd = &cobra.Command{
	Use:    "logout",
	Short:  "Log out from F5 Distributed Cloud.",
	Hidden: true, // Hide from help to match original f5xcctl
	Long:   `Clear saved credentials from the configuration file.`,
	Example: `  # Log out
  f5xcctl logout`,
	RunE: runLogout,
}

var whoamiCmd = &cobra.Command{
	Use:    "whoami",
	Short:  "Show current user information",
	Hidden: true, // Hide from help to match original f5xcctl
	Long:   `Display information about the currently authenticated user.`,
	Example: `  # Show current user
  f5xcctl whoami`,
	RunE: runWhoami,
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(whoamiCmd)

	loginCmd.Flags().StringVar(&loginFlags.tenant, "tenant", "", "Tenant name (e.g., example-tenant from example-tenant.console.ves.volterra.io)")
	loginCmd.Flags().StringVar(&loginFlags.p12Bundle, "p12-bundle", "", "Path to P12 certificate bundle")
	loginCmd.Flags().StringVar(&loginFlags.cert, "cert", "", "Path to client certificate")
	loginCmd.Flags().StringVar(&loginFlags.key, "key", "", "Path to client key")
	loginCmd.Flags().BoolVar(&loginFlags.apiToken, "api-token", false, "Use API token from F5XC_API_TOKEN environment variable")
}

func runLogin(cmd *cobra.Command, args []string) error {
	configPath := getConfigPath()

	// Load existing config
	config := &ConfigFile{}
	if data, err := os.ReadFile(configPath); err == nil {
		_ = yaml.Unmarshal(data, config) // Ignore error - we'll use defaults if parsing fails
	}

	// Update from flags
	if loginFlags.tenant != "" {
		config.ServerURL = fmt.Sprintf("https://%s.console.ves.volterra.io/api", loginFlags.tenant)
	}

	// Handle API token authentication
	if loginFlags.apiToken {
		token := os.Getenv("F5XC_API_TOKEN")
		if token == "" {
			return fmt.Errorf("F5XC_API_TOKEN environment variable not set")
		}
		config.APIToken = true
		config.P12Bundle = "" // Clear other auth methods
		config.Cert = ""
		config.Key = ""
	} else {
		config.APIToken = false // Ensure API token is cleared when using other auth methods

		if loginFlags.p12Bundle != "" {
			config.P12Bundle = expandPath(loginFlags.p12Bundle)
			config.Cert = ""
			config.Key = ""
		}

		if loginFlags.cert != "" {
			config.Cert = expandPath(loginFlags.cert)
			config.P12Bundle = ""
		}

		if loginFlags.key != "" {
			config.Key = expandPath(loginFlags.key)
		}
	}

	// Validate we have credentials
	hasP12 := config.P12Bundle != ""
	hasCertKey := config.Cert != "" && config.Key != ""
	hasAPIToken := config.APIToken

	if !hasP12 && !hasCertKey && !hasAPIToken {
		return fmt.Errorf("authentication credentials required: provide --p12-bundle, --cert and --key, or --api-token")
	}

	if hasP12 {
		// Check for P12 password
		if os.Getenv("F5XC_P12_PASSWORD") == "" {
			fmt.Println("Warning: F5XC_P12_PASSWORD environment variable not set")
			fmt.Println("Set it with: export F5XC_P12_PASSWORD='your-password'")
		}

		// Verify P12 file exists
		if _, err := os.Stat(config.P12Bundle); os.IsNotExist(err) {
			return fmt.Errorf("P12 bundle not found: %s", config.P12Bundle)
		}
	}

	if hasCertKey {
		// Verify cert and key files exist
		if _, err := os.Stat(config.Cert); os.IsNotExist(err) {
			return fmt.Errorf("certificate file not found: %s", config.Cert)
		}
		if _, err := os.Stat(config.Key); os.IsNotExist(err) {
			return fmt.Errorf("key file not found: %s", config.Key)
		}
	}

	// Validate we have a server URL
	if config.ServerURL == "" {
		return fmt.Errorf("server URL required: provide --tenant flag")
	}

	// Save configuration
	if err := saveConfig(config, configPath); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	// Test the connection
	output.PrintInfo("Authenticating...")

	// Reinitialize client with new config
	initConfig()

	client := GetClient()
	if client == nil {
		return fmt.Errorf("failed to initialize client - check your credentials")
	}

	// Try to get user info to verify authentication
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, "/api/web/custom/namespace/system/whoami", nil)
	if err != nil {
		// Try alternative endpoint
		resp, err = client.Get(ctx, "/api/web/namespaces", nil)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("authentication failed (status %d): check your credentials", resp.StatusCode)
	}

	// Extract tenant from URL
	tenant := ""
	if config.ServerURL != "" {
		parts := strings.Split(config.ServerURL, ".")
		if len(parts) > 0 {
			tenant = strings.TrimPrefix(parts[0], "https://")
		}
	}

	fmt.Println()
	output.PrintInfo(fmt.Sprintf("Successfully logged in to %s", tenant))
	output.PrintInfo(fmt.Sprintf("Configuration saved to %s", configPath))

	return nil
}

func runLogout(cmd *cobra.Command, args []string) error {
	configPath := getConfigPath()

	// Load existing config
	config := &ConfigFile{}
	if data, err := os.ReadFile(configPath); err == nil {
		_ = yaml.Unmarshal(data, config) // Ignore error - we'll use defaults if parsing fails
	}

	// Clear credentials but keep server URL
	config.P12Bundle = ""
	config.Cert = ""
	config.Key = ""
	config.APIToken = false

	// Save configuration
	if err := saveConfig(config, configPath); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	output.PrintInfo("Logged out successfully")
	output.PrintInfo("Credentials cleared from configuration")

	return nil
}

func runWhoami(cmd *cobra.Command, args []string) error {
	client := GetClient()
	if client == nil {
		return fmt.Errorf("not logged in - run 'f5xcctl login' first")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to get user info
	resp, err := client.Get(ctx, "/api/web/custom/namespace/system/whoami", nil)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	if resp.StatusCode >= 400 {
		// Try alternative approach - get tenant info
		resp, err = client.Get(ctx, "/api/web/namespaces/system", nil)
		if err != nil {
			return fmt.Errorf("failed to get user info: %w", err)
		}
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to get user info (status %d)", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Build user info display
	configPath := getConfigPath()
	config := &ConfigFile{}
	if data, err := os.ReadFile(configPath); err == nil {
		_ = yaml.Unmarshal(data, config) // Ignore error - we'll use defaults if parsing fails
	}

	userInfo := map[string]interface{}{}

	if config.ServerURL != "" {
		userInfo["server"] = config.ServerURL
		// Extract tenant from URL
		parts := strings.Split(config.ServerURL, ".")
		if len(parts) > 0 {
			userInfo["tenant"] = strings.TrimPrefix(parts[0], "https://")
		}
	}

	if config.APIToken {
		userInfo["auth_method"] = "API Token"
	} else if config.P12Bundle != "" {
		userInfo["auth_method"] = "P12 Bundle"
		userInfo["p12_bundle"] = filepath.Base(config.P12Bundle)
	} else if config.Cert != "" {
		userInfo["auth_method"] = "Certificate"
		userInfo["cert"] = filepath.Base(config.Cert)
	}

	// Add info from API response if available
	if name, ok := result["name"].(string); ok {
		userInfo["namespace"] = name
	}

	return output.Print(userInfo, GetOutputFormat())
}
