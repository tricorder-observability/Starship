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
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/tricorder/src/api-server/http/client"
	"github.com/tricorder/src/cli/pkg/output"
	"github.com/tricorder/src/utils/log"
)

var undeployCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "Undeploy an previously-deployed eBPF+WASM module",
	Long: "Undeploy an previously-deployed eBPF+WASM module. For example:\n" +
		"$ starship-cli module undeploy --api-server=<address> --id ce8a4fbe_45db_49bb_9568_6688dd84480b",
	Run: func(cmd *cobra.Command, args []string) {
		client := client.NewClient(apiServerAddress)
		resp, err := client.UndeployModule(moduleId)
		if err != nil {
			log.Error(err)
			return
		}

		// TODO(jun): refactor output to delete this hack
		// we can upgrade golang version and introduce generic code
		// to provide a generic interface to output
		respByte, err := json.Marshal(resp)
		if err != nil {
			log.Error(err)
			return
		}

		err = output.Print(outputFormat, respByte)
		if err != nil {
			log.Error(err)
		}
	},
}

func init() {
	undeployCmd.Flags().StringVarP(&moduleId, "id", "i", moduleId, "the ID of module.")
	_ = undeployCmd.MarkFlagRequired("id")
}
