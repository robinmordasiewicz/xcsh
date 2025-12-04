package cmd

import (
	"github.com/spf13/cobra"
)

var apiEndpointCmd = &cobra.Command{
	Use:     "api-endpoint",
	Aliases: []string{"apie"},
	Short:   "API endpoint discovery and control",
	Long: `Discover and control API endpoints in your service mesh.

This command group provides operations for:
  - Discovering API endpoints between services
  - Creating layer7 policies based on discovered endpoints`,
	Example: `  # Discover API endpoints in a namespace
  f5xc api-endpoint discover --namespace my-ns

  # Create layer7 policy from discovered endpoint
  f5xc api-endpoint control --src-service frontend --dst-service backend --method GET --path /api/users`,
}

func init() {
	rootCmd.AddCommand(apiEndpointCmd)
}
