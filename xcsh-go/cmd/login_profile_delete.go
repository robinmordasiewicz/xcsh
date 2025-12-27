package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/robinmordasiewicz/xcsh/pkg/profile"
)

var profileDeleteFlags struct {
	force bool
}

var loginProfileDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete an authentication profile.",
	Long: `Delete an authentication profile.

This permanently removes the profile configuration file.
The default profile cannot be deleted; switch to another profile first.

Use --force to skip confirmation.`,
	Example: `  # Delete a profile
  xcsh login profile delete staging

  # Delete without confirmation
  xcsh login profile delete staging --force`,
	Args: cobra.ExactArgs(1),
	RunE: runLoginProfileDelete,
}

func init() {
	loginProfileCmd.AddCommand(loginProfileDeleteCmd)

	loginProfileDeleteCmd.Flags().BoolVar(&profileDeleteFlags.force, "force", false, "Skip confirmation prompt")
}

func runLoginProfileDelete(cmd *cobra.Command, args []string) error {
	profileName := args[0]

	manager, err := profile.NewManager()
	if err != nil {
		return fmt.Errorf("failed to initialize profile manager: %w", err)
	}

	// Check if profile exists
	if !manager.Exists(profileName) {
		return fmt.Errorf("profile %q not found", profileName)
	}

	// Check if it's the default profile
	if manager.IsDefault(profileName) {
		names, _ := manager.List()
		if len(names) > 1 {
			return fmt.Errorf("cannot delete default profile %q\n\nSwitch to another profile first:\n  xcsh login profile use <other-profile>", profileName)
		}
	}

	// Confirm deletion unless --force
	if !profileDeleteFlags.force {
		fmt.Printf("Delete profile %q? This cannot be undone. [y/N]: ", profileName)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" && response != "yes" && response != "Yes" {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	// Delete the profile
	if err := manager.Delete(profileName); err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	fmt.Printf("Deleted profile: %s\n", profileName)
	return nil
}
