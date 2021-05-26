package cmd

import (
	"fmt"
	"os"

	"github.com/obay/rcmt/rcmtpackage"
	"github.com/spf13/cobra"
)

var removePackageCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a package resource file",
	Long: `This command removes an existing package resource
file. It doesn't do anything more than simply deleting
the file. This will not delete the resource from your 
targets. If the resource file does not exist, you will 
get an error message stating so`,
	Run: func(cmd *cobra.Command, args []string) {
		removePackageCmdMain(args)
	},
}

func init() {
	packageCmd.AddCommand(removePackageCmd)
}

func removePackageCmdMain(args []string) {
	if len(args) != 1 {
		// To-do: replace the following message with the full help message
		fmt.Fprintf(os.Stderr, "A package name must be provided\n")
		os.Exit(1)
	}
	err := rcmtpackage.RemovePackage(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
