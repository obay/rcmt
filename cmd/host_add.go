package cmd

import (
	"fmt"
	"os"

	"github.com/obay/rcmt/rcmthost"
	"github.com/spf13/cobra"
)

var addHostCmd = &cobra.Command{
	Use:   "add [username@]<hostname>[:port]",
	Short: "Adds a host to the hosts.rcmt file.",
	Long:  "Adds a host to the hosts.rcmt file.",
	Run: func(cmd *cobra.Command, args []string) {
		addHostCmdMain(args)
	},
}

var Name string

func init() {
	hostCmd.AddCommand(addHostCmd)
	addHostCmd.Flags().StringVarP(&Name, "name", "n", "", "Name for the host (required)")
	addHostCmd.MarkFlagRequired("name")
}

func addHostCmdMain(args []string) {
	if len(args) != 1 {
		// To-do: replace the following message with the full help message
		fmt.Fprintf(os.Stderr, "A hostname must be provided\n")
		os.Exit(1)
	}
	sshConnectionParameters := rcmthost.ParseConnectionString(args[0])
	sshConnectionParameters.Name = Name
	err := rcmthost.AddHostUsingSSHConnectionParameters(sshConnectionParameters)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
