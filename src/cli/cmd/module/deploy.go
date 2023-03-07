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

	"github.com/tricorder/src/api-server/http/api"
	"github.com/tricorder/src/utils/log"

	outputs "github.com/tricorder/src/cli/internal/outputs"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy Starship module to the platform.",
	Long: "Deploy a previously-created module to the Starship platform.\n" +
		"The deployed module, if succeeded, will be executed by Tricorder.\n" +
		"For example:\n" +
		"    starship-cli module deploy --id=ce8a4fbe_45db_49bb_9568_6688dd84480b",
	Run: func(cmd *cobra.Command, args []string) {
		url := api.GetURL(apiServerAddress, api.DEPLOY_MODULE_PATH)
		resp, err := deployModule(url, moduleId)
		if err != nil {
			log.Fatalf("Failed to deploy module, id='%s', error: %v", moduleId, err)
		}
		if len(resp) == 0 {
			log.Fatalf("Failed to deploy module, id='%s', Empty response from API Server", moduleId)
		}
		if err := outputs.Output(output, resp); err != nil {
			log.Fatalf("Failed to write output, error: %v", err)
		}
	},
}

func init() {
	deployCmd.Flags().StringVarP(&moduleId, "id", "i", moduleId, "the id of module.")
	_ = deployCmd.MarkFlagRequired("id")
}

func deployModule(url string, moduleId string) ([]byte, error) {
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
