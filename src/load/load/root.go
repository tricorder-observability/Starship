package load

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "starship-load",
	Short: "The CLI (Command Line Interface) for starship",
	Long:  `CLI to manage starship observe modules.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(loadCmd)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.SuggestionsMinimumDistance = 1
}
