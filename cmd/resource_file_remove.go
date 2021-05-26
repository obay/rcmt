package cmd

import (
	"fmt"
	"os"

	"github.com/obay/rcmt/rcmtfile"
	"github.com/spf13/cobra"
)

var removeFileCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a file resource file",
	Long: `This command removes an existing file resource
file. It doesn't do anything more than simply deleting
the file. This will not delete the resource from your 
targets. If the resource file does not exist, you will 
get an error message stating so`,
	Run: func(cmd *cobra.Command, args []string) {
		removeFileCmdMain(args)
	},
}

func init() {
	fileCmd.AddCommand(removeFileCmd)
}

func removeFileCmdMain(args []string) {
	if len(args) != 1 {
		// To-do: replace the following message with the full help message
		fmt.Fprintf(os.Stderr, "A file name must be provided\n")
		os.Exit(1)
	}
	err := rcmtfile.RemoveFile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
