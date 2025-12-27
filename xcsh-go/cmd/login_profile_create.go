package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/robinmordasiewicz/xcsh/pkg/profile"
)

var profileCreateFlags struct {
	name             string
	apiURL           string
	apiToken         string
	p12Bundle        string
	cert             string
	key              string
	defaultNamespace string
	setDefault       bool
}

var loginProfileCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new authentication profile.",
	Long: `Create a new F5 Distributed Cloud authentication profile.

Profiles allow storing multiple sets of credentials for different tenants.
Each profile requires:
  - A unique name
  - An API URL (console URL)
  - One authentication method (API token, P12 bundle, or cert/key pair)

Profile Configuration:
  Profiles are stored in ~/.config/xcsh/profiles/<name>.yaml with 0600 permissions.`,
	Example: `  # Create profile with API token
  xcsh login profile create --name production \
    --api-url https://tenant.console.ves.volterra.io \
    --api-token "your-api-token"

  # Create profile with P12 bundle
  xcsh login profile create --name staging \
    --api-url https://staging.console.ves.volterra.io \
    --p12-bundle ~/.xcsh/staging.p12

  # Create profile with cert/key and set as default
  xcsh login profile create --name dev \
    --api-url https://dev.console.ves.volterra.io \
    --cert ~/.xcsh/cert.pem \
    --key ~/.xcsh/key.pem \
    --set-default`,
	RunE: runLoginProfileCreate,
}

func init() {
	loginProfileCmd.AddCommand(loginProfileCreateCmd)

	loginProfileCreateCmd.Flags().StringVar(&profileCreateFlags.name, "name", "", "Profile name (required)")
	loginProfileCreateCmd.Flags().StringVar(&profileCreateFlags.apiURL, "api-url", "", "F5 XC console URL (required)")
	loginProfileCreateCmd.Flags().StringVar(&profileCreateFlags.apiToken, "api-token", "", "API token for authentication")
	loginProfileCreateCmd.Flags().StringVar(&profileCreateFlags.p12Bundle, "p12-bundle", "", "Path to P12 certificate bundle")
	loginProfileCreateCmd.Flags().StringVar(&profileCreateFlags.cert, "cert", "", "Path to certificate file")
	loginProfileCreateCmd.Flags().StringVar(&profileCreateFlags.key, "key", "", "Path to private key file")
	loginProfileCreateCmd.Flags().StringVar(&profileCreateFlags.defaultNamespace, "default-namespace", "", "Default namespace for operations")
	loginProfileCreateCmd.Flags().BoolVar(&profileCreateFlags.setDefault, "set-default", false, "Set as default profile")

	_ = loginProfileCreateCmd.MarkFlagRequired("name")
	_ = loginProfileCreateCmd.MarkFlagRequired("api-url")
}

func runLoginProfileCreate(cmd *cobra.Command, args []string) error {
	manager, err := profile.NewManager()
	if err != nil {
		return fmt.Errorf("failed to initialize profile manager: %w", err)
	}

	// Create profile from flags
	p := &profile.Profile{
		Name:             profileCreateFlags.name,
		APIURL:           profileCreateFlags.apiURL,
		APIToken:         profileCreateFlags.apiToken,
		P12Bundle:        profileCreateFlags.p12Bundle,
		Cert:             profileCreateFlags.cert,
		Key:              profileCreateFlags.key,
		DefaultNamespace: profileCreateFlags.defaultNamespace,
	}

	// Validate profile
	if err := p.Validate(); err != nil {
		return fmt.Errorf("invalid profile configuration: %w", err)
	}

	// Check if profile already exists
	if manager.Exists(p.Name) {
		return fmt.Errorf("profile %q already exists; use 'xcsh login profile delete %s' first", p.Name, p.Name)
	}

	// Create the profile
	if err := manager.Create(p); err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}

	fmt.Printf("Created profile %q\n", p.Name)
	fmt.Printf("  API URL: %s\n", p.APIURL)
	fmt.Printf("  Auth: %s\n", p.AuthMethod())

	// Set as default if requested or if it's the first profile
	profiles, _ := manager.List()
	if profileCreateFlags.setDefault || len(profiles) == 1 {
		if err := manager.SetDefault(p.Name); err != nil {
			return fmt.Errorf("failed to set default profile: %w", err)
		}
		fmt.Printf("  Status: default\n")
	}

	fmt.Printf("\nProfile saved to: %s/%s.yaml\n", manager.ProfilesDir(), p.Name)
	return nil
}
