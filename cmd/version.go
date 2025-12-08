package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Build-time variables (set via ldflags)
var (
	Version     = "dev"
	GitCommit   = "unknown"
	BuildDate   = "unknown"
	Branch      = "unknown"
	BuildAuthor = "unknown"
	BuildNumber = "0"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Print build version",
	Long:    `Print build version`,
	Example: `vesctl version`,
	Run: func(cmd *cobra.Command, args []string) {
		// Match original vesctl format exactly:
		// branch: 0-2-35 <br>commit-sha: 997bd8865ab5740ad6a787ac1c4619da3e5761c5 2022-09-27T09:11:13+00:00 mceloud 3089719293 nil <br>
		fmt.Printf("branch: %s <br>commit-sha: %s %s %s %s nil <br>\n",
			Branch, GitCommit, BuildDate, BuildAuthor, BuildNumber)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
