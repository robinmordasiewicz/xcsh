package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Site-specific flags matching original f5xcctl
var (
	siteLogColor    bool
	siteLogFabulous bool
	siteLogLevel    int
)

var siteCmd = &cobra.Command{
	Use:     "site",
	Aliases: []string{"s"},
	Short:   "Deploy and manage F5 XC sites on public cloud providers.",
	Long: `Deploy and manage F5 Distributed Cloud Customer Edge (CE) sites.

Sites are deployment points for F5 XC services in your infrastructure.
This command group supports creating, managing, and monitoring sites
across multiple cloud providers using Terraform automation.

SITE TYPES:
  aws_vpc      Deploy CE nodes in AWS Virtual Private Cloud
  azure_vnet   Deploy CE nodes in Azure Virtual Network

LIFECYCLE:
  1. Create site configuration in F5 XC console
  2. Run Terraform to provision cloud infrastructure
  3. Monitor site status until nodes report ONLINE
  4. Optionally destroy infrastructure and delete site

AI assistants should use 'f5xcctl site <provider> --help' for provider-specific
options and 'f5xcctl site <provider> run --help' for Terraform actions.`,
	Example: `  # Create an AWS VPC site from configuration
  f5xcctl site aws_vpc create -i aws-site.yaml

  # Delete an AWS VPC site
  f5xcctl site aws_vpc delete --name example-site

  # Run Terraform to provision AWS infrastructure
  f5xcctl site aws_vpc run --name example-site --action apply --auto-approve

  # Create an Azure VNet site
  f5xcctl site azure_vnet create -i azure-site.yaml

  # Check available site commands
  f5xcctl site aws_vpc --help`,
}

func init() {
	rootCmd.AddCommand(siteCmd)

	// Enable AI-agent-friendly error handling for invalid subcommands
	siteCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("unknown command %q for %q\n\nUsage: f5xcctl site <provider> <action> [flags]\n\nAvailable providers:\n  aws_vpc, azure_vnet\n\nRun 'f5xcctl site --help' for usage", args[0], cmd.CommandPath())
		}
		return cmd.Help()
	}
	siteCmd.SuggestionsMinimumDistance = 2

	// Site-specific flags matching original f5xcctl
	siteCmd.PersistentFlags().BoolVar(&siteLogColor, "log-color", true, "Enable colored log output.")
	siteCmd.PersistentFlags().BoolVar(&siteLogFabulous, "log-fabulous", true, "Enable enhanced log formatting.")
	siteCmd.PersistentFlags().IntVar(&siteLogLevel, "log-level", 3, "Set the logging verbosity level for site operations.")
}
