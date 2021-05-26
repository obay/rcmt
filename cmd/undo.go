package cmd

import (
	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Undo the tasks in the current folder",
	Long:  "Undo the tasks in the current folder. You can `--auto-approve` if you don't want to be asked for a confirmation",
	Run: func(cmd *cobra.Command, args []string) {
		undoCmdMain(args)
	},
}

func init() {
	rootCmd.AddCommand(undoCmd)
}

func undoCmdMain(args []string) {

}
