package cmd

import (
	"github.com/spf13/cobra"
)

var siteCmd = &cobra.Command{
	Use:   "site",
	Short: "Cloud site management operations",
	Long: `Manage cloud sites for F5 Distributed Cloud.

This command group provides operations for managing cloud sites including:
  - AWS VPC sites
  - Azure VNet sites

Each site type supports create, delete, replace, and run (Terraform) operations.`,
	Example: `  # Create an AWS VPC site
  f5xc site aws-vpc create --name my-site --region us-east-1

  # Run Terraform plan for an AWS site
  f5xc site aws-vpc run --name my-site --action plan

  # Delete an Azure VNet site
  f5xc site azure-vnet delete --name my-site`,
}

func init() {
	rootCmd.AddCommand(siteCmd)
}
