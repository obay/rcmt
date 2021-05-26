package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Add/remove Debian package resource files",
	Long: `The package commands allow you to easily create rcmt 
resource files for Debian packages. The subcommands 
add/remove only create and delete the files needed. 
There is no interaction with the hosts unless you do 
rcmt do or rcmt undo`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	resourceCmd.AddCommand(packageCmd)
}
