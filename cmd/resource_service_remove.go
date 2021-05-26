package cmd

import (
	"fmt"
	"os"

	"github.com/obay/rcmt/rcmtservice"
	"github.com/spf13/cobra"
)

var removeServiceCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a service resource file",
	Long: `This command removes an existing service resource
file. It doesn't do anything more than simply deleting
the file. This will not delete the resource from your 
targets. If the resource file does not exist, you will 
get an error message stating so`,
	Run: func(cmd *cobra.Command, args []string) {
		removeServiceCmdMain(args)
	},
}

func init() {
	serviceCmd.AddCommand(removeServiceCmd)
}

func removeServiceCmdMain(args []string) {
	if len(args) != 1 {
		// To-do: replace the following message with the full help message
		fmt.Fprintf(os.Stderr, "A service name must be provided\n")
		os.Exit(1)
	}
	err := rcmtservice.RemoveService(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
