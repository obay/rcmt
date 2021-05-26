package cmd

import (
	"fmt"

	"github.com/obay/rcmt/rcmtresource"
	"github.com/spf13/cobra"
)

var listResourceCmd = &cobra.Command{
	Use:   "list",
	Short: "List all resources in the current folder",
	Long: `This command lists all rcmt resources in the current folder. It is not possible to point rcmt 
to pick up configuration from any other folder. This command will go through all *.rcmt files
and parse them to check their validaity and print an output of all resources`,
	Run: func(cmd *cobra.Command, args []string) {
		listResourceCmdMain(args)
	},
}

func init() {
	resourceCmd.AddCommand(listResourceCmd)
}

func listResourceCmdMain(args []string) {
	resources := rcmtresource.LoadResources()
	fmt.Printf("%v", resources)
}
