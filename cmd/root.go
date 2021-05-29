package cmd

import (
	"fmt"

	"github.com/obay/rcmt/helpers"
	"github.com/spf13/cobra"
	"github.com/tcnksm/go-latest"

	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "rcmt",
	Short: "rcmt configuration managment tool",
	Long: `rcmt stands for Rudimentary Configuration Management Tool.
It is a configuration tool created as part of Slack interview
challenge to demonstrate development and tooling skills. Think
of rcmt as Ansible using Golang`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
	checkLatestVersion()
}

func checkLatestVersion() {
	githubTag := &latest.GithubTag{
		Owner:             "obay",
		Repository:        "rcmt",
		FixVersionStrFunc: latest.DeleteFrontV(),
	}
	res, _ := latest.Check(githubTag, VersionString)
	if res.Outdated {
		helpers.PrintWarning(fmt.Sprintf("You are running rcmt "+VersionString+" which is not the latest, you should upgrade to v%s", res.Current))
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}
