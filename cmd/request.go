package cmd

import (
	"github.com/spf13/cobra"
)

var requestCmd = &cobra.Command{
	Use:     "request",
	Aliases: []string{"req", "r"},
	Short:   "Execute API requests",
	Long: `Execute direct API requests to F5 Distributed Cloud.

This command group provides low-level access to the F5 XC API for:
  - Generic RPC invocation
  - Secret management operations
  - Command sequence execution`,
	Example: `  # Execute a generic RPC
  f5xc request rpc namespace.CustomAPI.List -i request.yaml

  # Get public key for secret encryption
  f5xc request secrets get-public-key

  # Execute a command sequence
  f5xc request command-sequence -i objects.yaml --operation create`,
}

func init() {
	rootCmd.AddCommand(requestCmd)
}
