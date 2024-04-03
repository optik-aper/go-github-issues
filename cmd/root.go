package cmd

import (
	"fmt"
	"os"

	"github.com/optik-aper/go-github-issue-cli/v1/cmd/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "ghcli",
	Short: "ghcli is a command line interface to check github stuff",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("config", "", "path to the config file")
	if err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config")); err != nil {
		fmt.Printf("error binding root pflag 'config': %v", err)
	}

	initConfig()

	rootCmd.AddCommand(
		list.NewCmdList(),
	)
}

func initConfig() {
	configFilePath := viper.GetString("config")

	if configFilePath != "" {
		viper.SetConfigFile(configFilePath)
	} else {
		viper.SetConfigName(".ghcli")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME")
	}

	fmt.Println(configFilePath)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("error reading config : %v", err)
		os.Exit(1)
	}
}
