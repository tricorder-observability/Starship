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

package agent

import (
	"encoding/json"

	"github.com/tricorder/src/utils/log"

	apiserver "github.com/tricorder/src/api-server/http"
	"github.com/tricorder/src/cli/pkg/output"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List agents",
	Long: "List agents. For example:\n" +
		"$ starship-cli agent list --api-server=<address>",
	Run: func(cmd *cobra.Command, args []string) {
		client := apiserver.NewClient(apiServerAddress)
		resp, err := client.ListAgents(nil)
		if err != nil {
			log.Error(err)
			return
		}

		// TODO(jun): refactor output to delete this hack
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
