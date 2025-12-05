package cmd

import (
	"github.com/spf13/cobra"
)

// Site-specific flags matching original vesctl
var (
	siteLogColor    bool
	siteLogFabulous bool
	siteLogLevel    int
)

var siteCmd = &cobra.Command{
	Use:     "site",
	Aliases: []string{"s"},
	Short:   "Manage site creation through view public cloud site apis",
	Long:    `Manage site creation through view public cloud site apis`,
	Example: `vesctl site aws_vpc create`,
}

func init() {
	rootCmd.AddCommand(siteCmd)

	// Site-specific flags matching original vesctl
	siteCmd.PersistentFlags().BoolVar(&siteLogColor, "log-color", true, "enable color for your logs")
	siteCmd.PersistentFlags().BoolVar(&siteLogFabulous, "log-fabulous", true, "enable fabulous writer for your logs")
	siteCmd.PersistentFlags().IntVar(&siteLogLevel, "log-level", 3, "Log Level for Site Deployment")
}
