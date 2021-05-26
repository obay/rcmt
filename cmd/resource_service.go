package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Add/remove service resource files",
	Long: `The service commands allow you to easily create rcmt 
	resource files for Debian services. The subcommands 
	add/remove only create and delete the files needed. 
	There is no interaction with the hosts unless you do 
	rcmt do or rcmt undo`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	resourceCmd.AddCommand(serviceCmd)
}
