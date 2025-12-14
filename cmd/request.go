package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var requestCmd = &cobra.Command{
	Use:     "request",
	Aliases: []string{"req", "r"},
	Short:   "Execute custom API requests to F5 Distributed Cloud.",
	Long:    `Execute custom API requests to F5 Distributed Cloud.`,
	Example: `f5xcctl request secrets encrypt --policy-doc temp_policy --public-key pub_key secret`,
}

func init() {
	rootCmd.AddCommand(requestCmd)

	// Enable AI-agent-friendly error handling for invalid subcommands
	requestCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("unknown command %q for %q\n\nUsage: f5xcctl request <service> <action> [flags]\n\nAvailable services:\n  secrets, rpc, command-sequence\n\nRun 'f5xcctl request --help' for usage", args[0], cmd.CommandPath())
		}
		return cmd.Help()
	}
	requestCmd.SuggestionsMinimumDistance = 2
}
