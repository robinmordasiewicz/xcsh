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

	"github.com/robinmordasiewicz/f5xc/pkg/output"
)

var loginFlags struct {
	tenant    string
	p12Bundle string
	cert      string
	key       string
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to F5 Distributed Cloud",
	Long: `Authenticate with F5 Distributed Cloud.

This command validates your credentials and saves them to the configuration file.
You can authenticate using either:
  - A P12 certificate bundle (set VES_P12_PASSWORD environment variable)
  - Certificate and key files

After successful login, you can use all f5xc commands to manage your resources.`,
	Example: `  # Login with P12 bundle
  f5xc login --tenant my-tenant --p12-bundle ~/.vesctl/my-cert.p12

  # Login with certificate and key
  f5xc login --tenant my-tenant --cert ~/.vesctl/cert.pem --key ~/.vesctl/key.pem

  # Login (using existing configuration)
  f5xc login`,
	RunE: runLogin,
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from F5 Distributed Cloud",
	Long:  `Clear saved credentials from the configuration file.`,
	Example: `  # Log out
  f5xc logout`,
	RunE: runLogout,
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current user information",
	Long:  `Display information about the currently authenticated user.`,
	Example: `  # Show current user
  f5xc whoami`,
	RunE: runWhoami,
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(whoamiCmd)

	loginCmd.Flags().StringVar(&loginFlags.tenant, "tenant", "", "Tenant name (e.g., my-tenant from my-tenant.console.ves.volterra.io)")
	loginCmd.Flags().StringVar(&loginFlags.p12Bundle, "p12-bundle", "", "Path to P12 certificate bundle")
	loginCmd.Flags().StringVar(&loginFlags.cert, "cert", "", "Path to client certificate")
	loginCmd.Flags().StringVar(&loginFlags.key, "key", "", "Path to client key")
}

func runLogin(cmd *cobra.Command, args []string) error {
	configPath := getConfigPath()

	// Load existing config
	config := &ConfigFile{}
	if data, err := os.ReadFile(configPath); err == nil {
		_ = yaml.Unmarshal(data, config)
	}

	// Update from flags
	if loginFlags.tenant != "" {
		serverURL := fmt.Sprintf("https://%s.console.ves.volterra.io/api", loginFlags.tenant)
		config.ServerURLs = []string{serverURL}
		config.Tenant = loginFlags.tenant
	}

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

	// Validate we have credentials
	hasP12 := config.P12Bundle != ""
	hasCertKey := config.Cert != "" && config.Key != ""

	if !hasP12 && !hasCertKey {
		return fmt.Errorf("authentication credentials required: provide --p12-bundle or --cert and --key")
	}

	if hasP12 {
		// Check for P12 password
		if os.Getenv("VES_P12_PASSWORD") == "" {
			fmt.Println("Warning: VES_P12_PASSWORD environment variable not set")
			fmt.Println("Set it with: export VES_P12_PASSWORD='your-password'")
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
	if len(config.ServerURLs) == 0 {
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

	tenant := config.Tenant
	if tenant == "" && len(config.ServerURLs) > 0 {
		// Extract tenant from URL
		parts := strings.Split(config.ServerURLs[0], ".")
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
		_ = yaml.Unmarshal(data, config)
	}

	// Clear credentials but keep server URL
	config.P12Bundle = ""
	config.Cert = ""
	config.Key = ""

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
		return fmt.Errorf("not logged in - run 'f5xc login' first")
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
		_ = yaml.Unmarshal(data, config)
	}

	userInfo := map[string]interface{}{
		"tenant": config.Tenant,
	}

	if len(config.ServerURLs) > 0 {
		userInfo["server"] = config.ServerURLs[0]
	}

	if config.P12Bundle != "" {
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
