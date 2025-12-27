package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var loginProfileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage authentication profiles.",
	Long: `Manage F5 Distributed Cloud authentication profiles.

Profiles allow storing multiple sets of credentials and switching between them.
Each profile can use different authentication methods (P12, certificate, API token).

Use this command to:
  - List configured profiles
  - Create a new profile
  - Switch between profiles
  - Delete a profile

Profile Configuration:
  Profiles are stored in ~/.config/xcsh/profiles/<name>.yaml
  Global settings (default profile) are in ~/.config/xcsh/config.yaml

AI assistants should use 'xcsh login profile show' to understand the currently
active profile before executing commands.`,
	Example: `  # List all configured profiles
  xcsh login profile list

  # Show current active profile
  xcsh login profile show

  # Create a new profile
  xcsh login profile create --name staging --api-url https://tenant.console.ves.volterra.io --api-token "your-token"

  # Switch to a profile
  xcsh login profile use staging

  # Delete a profile
  xcsh login profile delete staging`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("unknown command %q for %q\n\nUsage: xcsh login profile <action> [flags]\n\nAvailable Actions:\n  list, show, create, use, delete\n\nRun 'xcsh login profile --help' for usage", args[0], cmd.CommandPath())
		}
		return cmd.Help()
	},
}

func init() {
	loginCmd.AddCommand(loginProfileCmd)
}
