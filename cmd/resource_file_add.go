package cmd

import (
	"fmt"
	"os"

	"github.com/obay/rcmt/rcmtfile"
	"github.com/spf13/cobra"
)

var addFileCmd = &cobra.Command{
	Use:   "add",
	Short: "Create an rcmt file for a file resource",
	Long: `This command allows you to create an rcmt resource file for a file. 
By defaul, the desired state of the file will be \"exists\". This 
means that if you do rcmt do, the file will be created. If you want 
to delete the file, make sure to update the rcmt file accordingly`,
	Run: func(cmd *cobra.Command, args []string) {
		addFileCmdMain(args)
	},
}

func init() {
	fileCmd.AddCommand(addFileCmd)
}

func addFileCmdMain(args []string) {
	if len(args) != 1 {
		// To-do: replace the following message with the full help message
		fmt.Fprintf(os.Stderr, "A file name must be provided\n")
		os.Exit(1)
	}
	err := rcmtfile.AddFile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
