package cmd

import (
	"fmt"
	"runtime"
	runtimedebug "runtime/debug"

	"github.com/spf13/cobra"
)

// Build-time variables - set via ldflags or auto-detected from VCS info
var (
	Version   = "dev"
	GitCommit = "local"
	BuildDate = "now"
)

func init() {
	// Auto-detect version info from Go's embedded VCS data
	// Only override if still at defaults (not set via ldflags)
	if info, ok := runtimedebug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				if GitCommit == "local" && setting.Value != "" {
					// Use short commit hash (7 chars) like git
					if len(setting.Value) > 7 {
						GitCommit = setting.Value[:7]
					} else {
						GitCommit = setting.Value
					}
				}
			case "vcs.time":
				if BuildDate == "now" && setting.Value != "" {
					BuildDate = setting.Value
				}
			}
		}
	}

	rootCmd.AddCommand(versionCmd)
}

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
