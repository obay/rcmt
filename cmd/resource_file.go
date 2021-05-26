package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "Add/remove file resource files",
	Long: `The file commands allow you to easily create rcmt 
resource files for files. The subcommands 
add/remove only create and delete the files needed. 
There is no interaction with the hosts unless you do 
rcmt do or rcmt undo`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	resourceCmd.AddCommand(fileCmd)
}
