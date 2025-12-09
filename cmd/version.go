package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// Build-time variables (set via ldflags)
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Display vesctl version and build information",
	Long:    `Display vesctl version and build information.`,
	Example: `vesctl version`,
	Run: func(cmd *cobra.Command, args []string) {
		// Short commit hash (7 chars like GitHub)
		commit := GitCommit
		if len(commit) > 7 {
			commit = commit[:7]
		}

		fmt.Printf("vesctl version %s\n", Version)
		fmt.Printf("  commit:   %s\n", commit)
		fmt.Printf("  built:    %s\n", BuildDate)
		fmt.Printf("  go:       %s\n", runtime.Version())
		fmt.Printf("  platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
