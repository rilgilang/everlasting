/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"everlasting/bootstrap"
	"everlasting/src/infrastructure/http"
	"everlasting/src/infrastructure/pkg"

	"github.com/sarulabs/di"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var container di.Container
var config = new(pkg.Config)
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Run AICare Marketplace Rest API",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		http.RunDashboardAPI(container, config)
	},
}

func Execute() {

	cmd := []*cobra.Command{
		{
			Use:   "dashboard",
			Short: "Run HTTP Server For Partner App",
			Run: func(cmd *cobra.Command, args []string) {
				http.RunDashboardAPI(container, config)
			},
		},
	}

	rootCmd.AddCommand(cmd...)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".cmd" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigType("env")
		viper.SetConfigName(".env")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		err := viper.Unmarshal(config)
		if err != nil {
			panic("error loading config")
		}
		container = bootstrap.InitializeContainer(config)
	}
}
