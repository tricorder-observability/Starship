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

	"github.com/tricorder/src/utils/log"

	"github.com/spf13/cobra"

	"github.com/tricorder/src/api-server/http/api"
	"github.com/tricorder/src/cli/internal/outputs"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete module from staship",
	Long: `Delete module from staship. For example:

$ starship-cli module delete --id 2a339411_7dd8_46ba_9581_e9d41286b564
`,
	Run: func(cmd *cobra.Command, args []string) {
		url := api.GetURL(apiServerAddress, api.DELETE_MODULE_PATH)
		resp, err := deleteModule(url, moduleId)
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
	deleteCmd.Flags().StringVarP(&moduleId, "id", "i", moduleId, "the id of module.")
	_ = deleteCmd.MarkFlagRequired("id")
}

func deleteModule(url string, moduleId string) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("%s?id=%s", url, moduleId))
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
