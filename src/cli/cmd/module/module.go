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

package module

import (
	"github.com/spf13/cobra"

	"github.com/tricorder/src/cli/pkg/kubernetes"
	"github.com/tricorder/src/utils/log"
)

var (
	// Address of API Server, specified from --api-server flag.
	apiServerAddress string

	// ID of a previousely-created eBPF+WASM module.
	moduleId string

	// The format of the output.
	outputFormat string
)

var ModuleCmd = &cobra.Command{
	Use:   "module",
	Short: "Manage eBPF+WASM modules",
	Long:  "Create, deploy, undeploy, delete, list eBPF+WASM modules",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// If Starship apiServerAddress is not set, try to get it from kubernetes
		if apiServerAddress == "" {
			newApiAddress, err := kubernetes.GetStarshipAPIAddress()
			if err != nil {
				log.Fatal("Failed to connect to Kubernetes API Server, " +
					"please manually set --api-server to the correct API Server address.")
			}
			apiServerAddress = newApiAddress
		}
	},
}

func init() {
	// Here you will define your flags and configuration settings.
	ModuleCmd.PersistentFlags().StringVar(&apiServerAddress, "api-server", "", "address of the Starship API Server.")
	ModuleCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "yaml", "format (yaml,json,table) of output.")

	ModuleCmd.AddCommand(listCmd)
	ModuleCmd.AddCommand(createCmd)
	ModuleCmd.AddCommand(deployCmd)
	ModuleCmd.AddCommand(deleteCmd)
	ModuleCmd.AddCommand(undeployCmd)
}
