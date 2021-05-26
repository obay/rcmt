package cmd

import (
	"fmt"
	"os"

	"github.com/obay/rcmt/rcmthost"
	"github.com/spf13/cobra"
)

var removeHostCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes a host from the hosts.rcmt file.",
	Long:  "Removes a host from the hosts.rcmt file. Only the hostname will be considered in the remove. All other parameters (name, port, username) will be discarded. Think of the hostname as the primary key in the hosts.rcmt file",
	Run: func(cmd *cobra.Command, args []string) {
		removeHostCmdMain(args)
	},
}

func init() {
	hostCmd.AddCommand(removeHostCmd)
}

func removeHostCmdMain(args []string) {
	if len(args) != 1 {
		// To-do: replace the following message with the full help message
		fmt.Fprintf(os.Stderr, "A hostname must be provided\n")
		os.Exit(1)
	}
	sshConnectionParameters := rcmthost.ParseConnectionString(args[0])
	err := rcmthost.RemoveHostUsingSSHConnectionParameters(sshConnectionParameters)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
