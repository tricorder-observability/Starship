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
	"fmt"

	"github.com/spf13/cobra"
)

var ModuleCmd = &cobra.Command{
	Use:   "module",
	Short: "manage module",
	Long: `manage module. For example:
	1. create module:
	$ starship-cli module create --bcc-file-path path/to/bcc_file --module-file-path path/to/module_request_json_file
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("module called")
	},
}

var (
	apiAddress string
	moduleId   string
	output     string
)

func init() {
	// Here you will define your flags and configuration settings.
	ModuleCmd.PersistentFlags().StringVar(&apiAddress, "api-address", "", "address of starship api server.")
	_ = ModuleCmd.MarkPersistentFlagRequired("api-address")

	ModuleCmd.PersistentFlags().StringVarP(&output, "output", "o", "yaml",
		"the style(json,yaml,table) of output, yaml is default.")

	ModuleCmd.AddCommand(listCmd)
	ModuleCmd.AddCommand(uploadCmd)
	ModuleCmd.AddCommand(createCmd)
	ModuleCmd.AddCommand(deployCmd)
	ModuleCmd.AddCommand(deleteCmd)
	ModuleCmd.AddCommand(undeployCmd)
}
