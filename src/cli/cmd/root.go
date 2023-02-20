// Copyright (C) 2023  tricorder-observability
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

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
