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
	Long:    `Deploy and manage F5 XC sites on public cloud providers.`,
	Example: `f5xcctl site aws_vpc create`,
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
