// Copyright (C) 2023  Tricorder Observability
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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tricorder/src/cli/cmd/agent"
	"github.com/tricorder/src/cli/cmd/module"
	"github.com/tricorder/src/utils/log"
)

var rootCmd = &cobra.Command{
	Use:   "starship-cli",
	Short: "The Starship CLI",
	Long:  `The CLI to use the Tricorder Starship Observability platform.`,

	// https://pkg.go.dev/github.com/spf13/cobra#section-readme
	SuggestionsMinimumDistance: 1,
}

var apiServerAddress string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	const apiServerFlagName = "api-server"
	rootCmd.AddCommand(module.ModuleCmd)
	rootCmd.AddCommand(agent.AgentCmd)
	rootCmd.PersistentFlags().StringVar(&apiServerAddress, apiServerFlagName,
		"localhost:8080", "address of Starship API Server.")
	err := viper.BindPFlag(apiServerFlagName, rootCmd.PersistentFlags().Lookup(apiServerFlagName))
	if err != nil {
		log.Fatalf("Could not bind flag: %v", err)
	}
}
