package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "Adds/removes resources",
	Long:  "Adds/removes resources (file/package/service) configuration files to the local directory. If no parameters are provided, the resource file will be created using the default settings. You are expected to verify the resource settings and amend it as needed before applying it.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(resourceCmd)
}
