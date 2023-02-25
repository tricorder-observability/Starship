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
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/tricorder/src/utils/log"

	"github.com/tricorder/src/cli/internal/outputs"
	modulepb "github.com/tricorder/src/pb/module"
	"github.com/tricorder/src/utils/file"
	http_utils "github.com/tricorder/src/utils/http"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create module by json file.",
	Long: `create module by json file. For example:
$ starship-cli module create -b path/to/bcc_file -m path/to/module_json_file -w path/to/wasm_file`,
	Run: func(cmd *cobra.Command, args []string) {
		bccStr, err := file.Read(bccFilePath)
		if err != nil {
			log.Fatalf("Failed to read --bcc-file-path='%s', error: %v", bccFilePath, err)
		}

		wasmBytes, err := file.ReadBin(wasmFilePath)
		if err != nil {
			log.Fatalf("Failed to read --wasm-file-path='%s', error: %v", wasmFilePath, err)
		}

		moduleReq, err := parseModuleJsonFile(moduleFilePath)
		if err != nil {
			log.Fatalf("Failed to read --module-json-path='%s', error: %v", moduleFilePath, err)
		}
		// override bcc code contet by bcc file
		moduleReq.Ebpf.Code = bccStr
		// override wasm code contet by wasm file
		moduleReq.Wasm.Code = wasmBytes
		url := http_utils.GetAPIUrl(apiAddress, http_utils.API_ROOT, http_utils.ADD_MODULE)
		resp, err := createModule(url, moduleReq)
		if err != nil {
			log.Error(err)
		}

		err = outputs.Output(output, resp)
		if err != nil {
			log.Error(err)
		}
	},
}

// the file path of module in json format flag
var (
	moduleFilePath string
	bccFilePath    string
	wasmFilePath   string
)

func init() {
	createCmd.Flags().StringVarP(&moduleFilePath, "module-json-path", "m",
		moduleFilePath, "The file path of module in json format.")
	createCmd.Flags().StringVarP(&bccFilePath, "bcc-file-path", "b", bccFilePath, "The file path of bcc code.")
	createCmd.Flags().StringVarP(&wasmFilePath, "wasm-file-path", "w", wasmFilePath, "The file path of wasm code.")
}

func parseModuleJsonFile(moduleJsonFilePath string) (*modulepb.Module, error) {
	bytes, err := file.ReadBin(moduleJsonFilePath)
	if err != nil {
		return nil, err
	}
	var moduleReq *modulepb.Module
	err = json.Unmarshal([]byte(bytes), &moduleReq)
	if err != nil {
		return nil, err
	}
	return moduleReq, nil
}

func createModule(url string, moduleReq *modulepb.Module) ([]byte, error) {
	bodyBytes, err := json.Marshal(moduleReq)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
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
