package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "Add/remove/list hosts",
	Long:  "Add/remove/list hosts",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(hostCmd)
}
