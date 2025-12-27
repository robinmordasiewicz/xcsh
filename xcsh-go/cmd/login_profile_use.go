package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"

	"github.com/robinmordasiewicz/xcsh/pkg/branding"
	"github.com/robinmordasiewicz/xcsh/pkg/client"
	"github.com/robinmordasiewicz/xcsh/pkg/output"
	"github.com/robinmordasiewicz/xcsh/pkg/profile"
)

var loginProfileUseCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Set the default authentication profile.",
	Long: `Set a profile as the default for all subsequent operations.

The default profile is used when:
  - No --profile flag is specified
  - F5XC_PROFILE environment variable is not set

In REPL mode, use '/login profile use <name>' to switch profiles.`,
	Example: `  # Set default profile
  xcsh login profile use production

  # In REPL mode
  xcsh> /login profile use staging
  Switched to profile: staging`,
	Args: cobra.ExactArgs(1),
	RunE: runLoginProfileUse,
}

func init() {
	loginProfileCmd.AddCommand(loginProfileUseCmd)
}

func runLoginProfileUse(cmd *cobra.Command, args []string) error {
	name := args[0]

	manager, err := profile.NewManager()
	if err != nil {
		return fmt.Errorf("failed to initialize profile manager: %w", err)
	}

	// Check if profile exists
	if !manager.Exists(name) {
		// List available profiles
		names, _ := manager.List()
		if len(names) == 0 {
			return fmt.Errorf("profile %q not found; no profiles configured\n\nCreate a profile with:\n  xcsh login profile create --name <name> --api-url <url> --api-token <token>", name)
		}
		return fmt.Errorf("profile %q not found\n\nAvailable profiles: %v", name, names)
	}

	// Set as default
	if err := manager.SetDefault(name); err != nil {
		return fmt.Errorf("failed to set default profile: %w", err)
	}

	// Also set as current for this session
	if err := manager.SetCurrent(name); err != nil {
		return fmt.Errorf("failed to set current profile: %w", err)
	}

	// Load profile to apply settings
	p, err := manager.Load(name)
	if err != nil {
		return fmt.Errorf("failed to load profile: %w", err)
	}

	// Apply profile settings to global variables and reinitialize client
	if err := applyProfileAndReinitClient(p); err != nil {
		return fmt.Errorf("failed to apply profile settings: %w", err)
	}

	// Display connection banner
	displayProfileSwitchBanner(p)

	return nil
}

// applyProfileAndReinitClient applies profile settings and reinitializes the API client.
func applyProfileAndReinitClient(p *profile.Profile) error {
	// Apply profile settings to global variables
	if p.APIURL != "" {
		normalized, err := client.NormalizeAPIURL(p.APIURL)
		if err != nil {
			serverURL = p.APIURL
		} else {
			serverURL = normalized
		}
	}

	if p.Cert != "" {
		cert = expandPath(p.Cert)
	}
	if p.Key != "" {
		key = expandPath(p.Key)
	}
	if p.P12Bundle != "" {
		p12Bundle = expandPath(p.P12Bundle)
	}

	// Set API token in environment for client initialization
	if p.APIToken != "" {
		if err := os.Setenv("F5XC_API_TOKEN", p.APIToken); err != nil {
			return fmt.Errorf("failed to set API token: %w", err)
		}
	}

	// Reinitialize the API client with new settings
	cfg := &client.Config{
		ServerURL: serverURL,
		Cert:      cert,
		Key:       key,
		CACert:    cacert,
		P12Bundle: p12Bundle,
		Debug:     debug,
		Timeout:   timeout,
	}

	if token := os.Getenv("F5XC_API_TOKEN"); token != "" {
		cfg.APIToken = token
	}

	newClient, err := client.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	apiClient = newClient

	// Apply default namespace from profile if set
	if p.DefaultNamespace != "" {
		// Clear namespace cache so validation uses fresh data
		ClearNamespaceCache()

		// Validate namespace exists
		if NamespaceExists(p.DefaultNamespace) {
			if err := setDefaultNamespace(p.DefaultNamespace); err != nil {
				output.PrintWarning(fmt.Sprintf(
					"Failed to set namespace '%s': %v", p.DefaultNamespace, err))
			}
		} else {
			output.PrintWarning(fmt.Sprintf(
				"Namespace '%s' from profile does not exist, using 'default'",
				p.DefaultNamespace))
			_ = setDefaultNamespace("default")
		}
	}

	return nil
}

// displayProfileSwitchBanner displays connection info after switching profiles.
// Uses buildConnectionInfo() as the single source of truth for connection details.
func displayProfileSwitchBanner(p *profile.Profile) {
	title := "Profile Switched: " + p.Name
	// Reuse buildConnectionInfo() from repl_banner.go for consistent connection details
	lines := buildConnectionInfo()

	// Render the framed box
	renderFramedBox(title, lines)
}

// renderFramedBox renders a title and content lines in a bordered box.
// This is reusable for any framed output display.
func renderFramedBox(title string, lines []string) {
	// Calculate max content width using display width (handles Unicode like •)
	maxWidth := runewidth.StringWidth(title)
	for _, line := range lines {
		if w := runewidth.StringWidth(line); w > maxWidth {
			maxWidth = w
		}
	}

	// Add padding: 2 left + 2 right = 4
	boxWidth := maxWidth + 4

	// Helper to print a padded line with proper display width calculation
	printLine := func(content string) {
		contentWidth := runewidth.StringWidth(content)
		padding := boxWidth - 2 - contentWidth // -2 for left padding
		if padding < 0 {
			padding = 0
		}
		fmt.Printf("%s│%s  %s%s%s%s\n",
			branding.ColorRed, branding.ColorReset,
			branding.ColorBoldWhite, content,
			strings.Repeat(" ", padding),
			branding.ColorRed+"│"+branding.ColorReset)
	}

	// Print banner
	fmt.Println()
	fmt.Printf("%s╭%s╮%s\n", branding.ColorRed, strings.Repeat("─", boxWidth), branding.ColorReset)
	printLine(title)
	fmt.Printf("%s├%s┤%s\n", branding.ColorRed, strings.Repeat("─", boxWidth), branding.ColorReset)
	for _, line := range lines {
		printLine(line)
	}
	fmt.Printf("%s╰%s╯%s\n", branding.ColorRed, strings.Repeat("─", boxWidth), branding.ColorReset)
	fmt.Println()
}
