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
	"io"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/tricorder/src/utils/log"

	"github.com/tricorder/src/cli/internal/outputs"
	http_utils "github.com/tricorder/src/utils/http"
)

var undeployCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "Undeploy module",
	Long: `Undeploy module command. For example:

$ starship-cli module undeploy --id ce8a4fbe_45db_49bb_9568_6688dd84480b
`,
	Run: func(cmd *cobra.Command, args []string) {
		url := http_utils.GetAPIUrl(apiAddress, http_utils.API_ROOT, http_utils.UN_DEPLOY)
		resp, err := undeployModule(url, moduleId)
		if err != nil {
			log.Error(err)
		}

		err = outputs.Output(output, resp)
		if err != nil {
			log.Error(err)
		}
	},
}

func init() {
	undeployCmd.Flags().StringVarP(&moduleId, "id", "i", moduleId, "the id of module.")
	_ = undeployCmd.MarkFlagRequired("id")
}

func undeployModule(url string, moduleId string) ([]byte, error) {
	resp, err := http.Post(fmt.Sprintf("%s?id=%s", url, moduleId), "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
