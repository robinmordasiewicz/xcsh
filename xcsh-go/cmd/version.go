package cmd

import (
	"fmt"
	"runtime"
	runtimedebug "runtime/debug"
	"time"

	"github.com/spf13/cobra"

	"github.com/robinmordasiewicz/xcsh/pkg/branding"
)

// Build-time version variables: set via ldflags during release, or auto-detected from VCS info
var (
	Version   = "dev"
	GitCommit = "local"
	BuildDate = "now"
)

func init() {
	// Generate development version if not set via ldflags
	// Format: YYYYMMDDHHMM-ALPHA (e.g., 202512261430-ALPHA)
	if Version == "dev" {
		Version = time.Now().UTC().Format("200601021504") + "-ALPHA"
	}

	// Update rootCmd.Version since it was initialized before init() ran
	rootCmd.Version = Version

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
	Use:    "version",
	Hidden: true, // Hide from help - users should use --version or -v flag
	Short:  fmt.Sprintf("Display %s version and build information.", branding.CLIName),
	Long: fmt.Sprintf(`Display %s version and build information.

Shows the current version, git commit hash, build date, Go version,
and platform information. Useful for debugging, support requests,
and verifying installation.

OUTPUT FIELDS:
  version   Release version (semver format)
  commit    Git commit hash (7-character short form)
  built     Build timestamp (ISO 8601 format)
  go        Go runtime version used to compile
  platform  Operating system and architecture (e.g., darwin/arm64)`, branding.CLIName),
	Example: fmt.Sprintf(`  # Show version information
  %s version`, branding.CLIName),
	Run: func(cmd *cobra.Command, args []string) {
		// Short commit hash (7 chars like GitHub)
		commit := GitCommit
		if len(commit) > 7 {
			commit = commit[:7]
		}

		fmt.Printf("%s version %s\n", branding.CLIName, Version)
		fmt.Printf("  commit:   %s\n", commit)
		fmt.Printf("  built:    %s\n", BuildDate)
		fmt.Printf("  go:       %s\n", runtime.Version())
		fmt.Printf("  platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}
