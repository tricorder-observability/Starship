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
	"time"

	"github.com/tricorder/src/utils/log"

	"github.com/tricorder/src/api-server/http/api"
	"github.com/tricorder/src/cli/pkg/outputs"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List eBPF+WASM modules",
	Long: "List eBPF+WASM modules. For example:\n" +
		"$ starship-cli module list --api-server=<address>",
	Run: func(cmd *cobra.Command, args []string) {
		url := api.GetURL(apiServerAddress, api.LIST_MODULE_PATH)
		resp, err := listModules(url)
		if err != nil {
			log.Error(err)
		}

		err = outputs.Output(output, resp)
		if err != nil {
			log.Error(err)
		}
	},
}

func listModules(url string) ([]byte, error) {
	c := http.Client{Timeout: time.Duration(3) * time.Second}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?fields=%s", url, "id,name,desire_state,create_time,"+
		"ebpf_fmt,ebpf_lang,schema_name,fn,schema_attr"), nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return body, nil
}
