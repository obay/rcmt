package cmd

import (
	"fmt"
	"os"

	"github.com/obay/rcmt/rcmthost"
	"github.com/obay/rcmt/rcmtresource"
	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Applies the tasks in the current folder",
	Long:  "Applies the tasks in the current folder. You can `--auto-approve` if you don't want to be asked for a confirmation.",
	Run: func(cmd *cobra.Command, args []string) {
		doCmdMain(args)
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}

func doCmdMain(args []string) {
	hosts := rcmthost.LoadHosts()
	err := rcmtresource.Do(hosts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
