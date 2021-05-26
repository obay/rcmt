package cmd

import (
	"fmt"
	"os"

	"github.com/obay/rcmt/rcmtservice"
	"github.com/spf13/cobra"
)

var addServiceCmd = &cobra.Command{
	Use:   "add",
	Short: "Create an rcmt file for a service resource",
	Long: `This command allows you to create an rcmt resource file for a service. 
By defaul, the desired state of the service will be \"runnig\". This 
means that if you do rcmt do, the service will be started. If you want 
to stop the service, make sure to update the rcmt file accordingly`,
	Run: func(cmd *cobra.Command, args []string) {
		addServiceCmdMain(args)
	},
}

func init() {
	serviceCmd.AddCommand(addServiceCmd)

}

func addServiceCmdMain(args []string) {
	if len(args) != 1 {
		// To-do: replace the following message with the full help message
		fmt.Fprintf(os.Stderr, "A service name must be provided\n")
		os.Exit(1)
	}
	err := rcmtservice.AddService(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
