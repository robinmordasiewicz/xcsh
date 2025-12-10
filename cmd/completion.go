package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish]",
	Short: "Generate shell completion scripts for bash, zsh, or fish.",
	Long: `To load completions:

Bash:
  $ source <(vesctl completion bash)

  # To load completions for each session, execute once:
  Linux:
    $ vesctl completion bash > /etc/bash_completion.d/vesctl
  MacOS:
    $ vesctl completion bash > /usr/local/etc/bash_completion.d/vesctl

Zsh:
  $ source <(vesctl completion zsh)

  # To load completions for each session, execute once:
  $ vesctl completion zsh > "${fpath[1]}/_vesctl"

Fish:
  $ vesctl completion fish | source

  # To load completions for each session, execute once:
  $ vesctl completion fish > ~/.config/fish/completions/vesctl.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			_ = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			_ = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			_ = cmd.Root().GenFishCompletion(os.Stdout, true)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
