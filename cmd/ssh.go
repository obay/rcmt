package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Verify your SSH access to the hosts",
	Long:  "Verify your SSH access to the hosts",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
}
