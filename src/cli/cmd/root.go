package cmd

import (
	"os"

	"github.com/spf13/viper"

	"github.com/tricorder/src/utils/log"

	"github.com/tricorder/src/cli/cmd/module"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "starship-cli",
	Short: "The CLI (Command Line Interface) for starship",
	Long:  `CLI to manage starship observe modules.`,
}

var apiAddress string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(module.ModuleCmd)
	rootCmd.PersistentFlags().StringVar(&apiAddress, "api-address", "localhost:8080", "address of starship api server.")
	err := viper.BindPFlag("api-address", rootCmd.PersistentFlags().Lookup("api-address"))
	if err != nil {
		log.Errorf("could not bind flag: %v", err)
	}
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.SuggestionsMinimumDistance = 1
}
