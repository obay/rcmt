package cmd

import (
	"fmt"
	"os"

	"github.com/obay/rcmt/helpers"
	"github.com/obay/rcmt/rcmthost"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listHostCmd = &cobra.Command{
	Use:   "list",
	Short: "List the hosts",
	Long:  "List the hosts in the hosts.rcmt file. This is the list of hosts where the tasks will be applied to. You can't have multiple hosts in the hosts.rcmt with the same hostname.",
	Run: func(cmd *cobra.Command, args []string) {
		listHostCmdMain(args)
	},
}

var (
	printJSON bool
)

func init() {
	hostCmd.AddCommand(listHostCmd)
	listHostCmd.PersistentFlags().BoolVarP(&printJSON, "json", "j", false, "JSON output")
}

func listHostCmdMain(args []string) {
	hosts := rcmthost.LoadHosts()
	if len(hosts) == 0 {
		fmt.Fprintf(os.Stderr, "No hosts found in the hosts.rcmt file or the file does not exist.\n")
		os.Exit(1)
	}
	if printJSON {
		jsonOutput, err := rcmthost.SerializeHostsToJSON(hosts)
		helpers.Check(err)
		fmt.Println(jsonOutput)
	} else {
		tabulateOutput(hosts)
	}
}

func tabulateOutput(hosts []rcmthost.HostDetails) {
	var data [][]string
	for _, host := range hosts {
		data = append(data, []string{host.Name, host.Hostname, host.Username, host.Port})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Hostname", "Username", "Port"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data) // Add Bulk Data
	table.Render()

}
