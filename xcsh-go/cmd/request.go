package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var requestCmd = &cobra.Command{
	Use:     "request",
	Aliases: []string{"req", "r"},
	Short:   "Execute custom API requests to F5 Distributed Cloud.",
	Long: `Execute custom API requests to F5 Distributed Cloud services.

This command group provides access to specialized F5 XC API operations
that aren't covered by standard configuration commands. Use these for
advanced operations like secret encryption, RPC calls, and command sequencing.

AVAILABLE SERVICES:
  secrets           Encrypt and manage secrets using F5 XC policy-based encryption
  rpc               Execute raw RPC calls to F5 XC API endpoints
  command-sequence  Run multiple API operations from a sequence file

AI assistants should use 'xcsh request <service> --help' for service-specific
options and available actions.`,
	Example: `  # Encrypt a secret using policy-based encryption
  xcsh request secrets encrypt --policy-doc policy.yaml --public-key key.pem secret

  # Execute a command sequence from file
  xcsh request command-sequence -i commands.yaml

  # Check available request services
  xcsh request --help`,
}

func init() {
	rootCmd.AddCommand(requestCmd)

	// Enable AI-agent-friendly error handling for invalid subcommands
	requestCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("unknown command %q for %q\n\nUsage: xcsh request <service> <action> [flags]\n\nAvailable services:\n  secrets, rpc, command-sequence\n\nRun 'xcsh request --help' for usage", args[0], cmd.CommandPath())
		}
		return cmd.Help()
	}
	requestCmd.SuggestionsMinimumDistance = 2
}
