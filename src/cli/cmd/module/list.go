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

	"github.com/tricorder/src/api-server/http/client"
	"github.com/tricorder/src/cli/pkg/output"
	"github.com/tricorder/src/utils/log"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List eBPF+WASM modules",
	Long: "List eBPF+WASM modules. For example:\n" +
		"$ starship-cli module list --api-server=<address>",
	Run: func(cmd *cobra.Command, args []string) {
		client := client.NewClient(apiServerAddress)
		resp, err := client.ListModules(nil)
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
