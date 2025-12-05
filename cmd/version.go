package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/robinmordasiewicz/f5xc/pkg/output"
)

// Build-time variables (set via ldflags)
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// VersionInfo contains version information
type VersionInfo struct {
	Version   string `json:"version" yaml:"version"`
	GitCommit string `json:"git_commit" yaml:"git_commit"`
	BuildDate string `json:"build_date" yaml:"build_date"`
	GoVersion string `json:"go_version" yaml:"go_version"`
	Platform  string `json:"platform" yaml:"platform"`
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Print build version",
	Long:    `Print build version`,
	Example: `vesctl version`,
	Run: func(cmd *cobra.Command, args []string) {
		info := VersionInfo{
			Version:   Version,
			GitCommit: GitCommit,
			BuildDate: BuildDate,
			GoVersion: runtime.Version(),
			Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		}

		format := GetOutputFormat()
		if format == "table" || format == "tsv" {
			// For version, just print simple output
			fmt.Printf("vesctl %s\n", Version)
			return
		}

		if err := output.Print(info, format); err != nil {
			output.PrintError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
