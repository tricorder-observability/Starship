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
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/tricorder/src/utils/log"

	"github.com/tricorder/src/api-server/dao"
	modulepb "github.com/tricorder/src/pb/module"
	"github.com/tricorder/src/utils/file"
	"github.com/tricorder/src/utils/uuid"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init module by json file.",
	Long: `init module by json file. For example:
$ starship-cli module init -b path/to/bcc_file -m path/to/module_json_file -w path/to/wasm_file`,
	Run: func(cmd *cobra.Command, args []string) {
		bccStr, err := file.Read(bccFilePath)
		if err != nil {
			log.Error(err)
		}

		wasmBytes, err := file.ReadBin(wasmFilePath)
		if err != nil {
			log.Error(err)
		}

		moduleReq, err := parseModuleJsonFile(moduleFilePath)
		if err != nil {
			log.Error(err)
		}
		// override bcc code contet by bcc file
		moduleReq.Ebpf.Code = bccStr
		// override wasm code contet by wasm file
		moduleReq.Wasm.Code = wasmBytes
		err = initModule(moduleReq)
		if err != nil {
			log.Error(err)
		}
	},
}

var dbFilePath string

func init() {
	initCmd.Flags().StringVarP(&moduleFilePath, "module-json-path", "m",
		moduleFilePath, "The file path of module in json format.")
	initCmd.Flags().StringVarP(&bccFilePath, "bcc-file-path", "b", bccFilePath, "The file path of bcc code.")
	initCmd.Flags().StringVarP(&wasmFilePath, "wasm-file-path", "w", wasmFilePath, "The file path of wasm code.")
	initCmd.Flags().StringVarP(&dbFilePath, "db-file-path", "d", dbFilePath, "The file path of sqlit db.")
}

func initModule(body *modulepb.Module) error {
	sqliteClient, _ := dao.InitSqlite(dbFilePath)
	codeDao := dao.ModuleDao{
		Client: sqliteClient,
	}

	ebpfProbes, err := json.Marshal(body.Ebpf.Probes)
	if err != nil {
		return err
	}

	schemaAttr, err := json.Marshal(body.Wasm.OutputSchema.Fields)
	if err != nil {
		return err
	}

	mod := &dao.ModuleGORM{
		ID:                 strings.Replace(uuid.New(), "-", "_", -1),
		Name:               body.Name,
		CreateTime:         time.Now().Format("2006-01-02 15:04:05"),
		DesiredState:       int(0),
		Ebpf:               body.Ebpf.Code,
		EbpfFmt:            int(body.Ebpf.Fmt),
		EbpfLang:           int(body.Ebpf.Lang),
		EbpfPerfBufferName: body.Ebpf.PerfBufferName,
		EbpfProbes:         string(ebpfProbes),
		Wasm:               body.Wasm.Code,
		SchemaAttr:         string(schemaAttr),
		Fn:                 body.Wasm.FnName,
		WasmFmt:            int(body.Wasm.Fmt),
		WasmLang:           int(body.Wasm.Lang),
	}

	mod.SchemaName = fmt.Sprintf("%s_%s", "tricorder_code", mod.ID)

	err = codeDao.SaveModule(mod)
	return err
}
