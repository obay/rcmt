package cmd

import (
	"github.com/spf13/cobra"

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
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}
