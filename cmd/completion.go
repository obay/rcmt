package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

  $ source <(rcmt completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ rcmt completion bash > /etc/bash_completion.d/rcmt
  # macOS:
  $ rcmt completion bash > /usr/local/etc/bash_completion.d/rcmt

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ rcmt completion zsh > "${fpath[1]}/_rcmt"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ rcmt completion fish | source

  # To load completions for each session, execute once:
  $ rcmt completion fish > ~/.config/fish/completions/rcmt.fish

PowerShell:

  PS> rcmt completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> rcmt completion powershell > rcmt.ps1
  # and source this file from your PowerShell profile.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintf(os.Stderr, "One of the options [bash|zsh|fish|powershell] must be provided\n")
			os.Exit(1)
		}
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
