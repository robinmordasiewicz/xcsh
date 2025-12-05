package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh]",
	Short: "Generate completion script",
	Long: `To load completions:
Bash:

$ source <(vesctl completion bash)

# To load completions for each session, execute once:
Linux:
  $ vesctl completion bash > /etc/bash_completion.d/yourprogram
MacOS:
  $ vesctl completion bash > /usr/local/etc/bash_completion.d/yourprogram

Zsh:

 $ source <(vesctl completion zsh)

 # To load completions for each session, execute once:
 $ vesctl completion zsh > "${fpath[1]}/_vesctl"
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			_ = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			_ = cmd.Root().GenZshCompletion(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
