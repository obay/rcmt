package cmd

import (
	"fmt"
	"os"

	"github.com/obay/rcmt/rcmtpackage"
	"github.com/spf13/cobra"
)

var addPackageCmd = &cobra.Command{
	Use:   "add",
	Short: "Create an rcmt file for a package resource",
	Long: `This command allows you to create an rcmt resource file for a package. 
By defaul, the desired state of the package will be \"installed\". This 
means that if you do rcmt do, the package will be installed. If you want 
to uninstall the package, make sure to update the rcmt file accordingly`,
	Run: func(cmd *cobra.Command, args []string) {
		addPackageCmdMain(args)
	},
}

func init() {
	packageCmd.AddCommand(addPackageCmd)
}

func addPackageCmdMain(args []string) {
	if len(args) != 1 {
		// To-do: replace the following message with the full help message
		fmt.Fprintf(os.Stderr, "A package name must be provided\n")
		os.Exit(1)
	}
	err := rcmtpackage.AddPackage(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
