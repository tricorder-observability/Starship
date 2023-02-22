// Copyright (C) 2023 Tricorder Observability
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

package load

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

var (
	bccFilePath    string
	wasmFilePath   string
	moduleFilePath string
	dbFilePath     string
)

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

func loadModule(body *modulepb.Module) error {
	sqliteClient, _ := dao.InitSqlite(dbFilePath)
	codeDao := dao.Module{
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
		Status:             int(0),
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

	err = codeDao.SaveCode(mod)
	return err
}

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "load module",
	Long: `[WARNING] Private command, not for general use. Load module,
	For example: $ starship-load load --bcc-file-path path/to/bcc_file --module-file-path path/to/module_request_json_file
	`,
	Run: func(cmd *cobra.Command, args []string) {
		bccStr, err := file.Read(bccFilePath)
		if err != nil {
			log.Errorf("read bcc file error: %v", err)
		}

		wasmBytes, err := file.ReadBin(wasmFilePath)
		if err != nil {
			log.Errorf("read wasm file error: %v", err)
		}

		moduleReq, err := parseModuleJsonFile(moduleFilePath)
		if err != nil {
			log.Errorf("read module json file error: %v", err)
		}
		// override bcc code contet by bcc file
		moduleReq.Ebpf.Code = bccStr
		// override wasm code contet by wasm file
		moduleReq.Wasm.Code = wasmBytes
		err = loadModule(moduleReq)
		if err != nil {
			log.Errorf("load module request error: %v", err)
		}
	},
}

func init() {
	loadCmd.Flags().StringVarP(&moduleFilePath, "module-file-path", "m",
		moduleFilePath, `The file path of module in json format.`)
	err := loadCmd.MarkFlagRequired("module-file-path")
	if err != nil {
		log.Errorf("set required flag error: %v", err)
	}
	loadCmd.Flags().StringVarP(&bccFilePath, "bcc-file-path", "b", bccFilePath, `The file path of bcc code.`)
	err = loadCmd.MarkFlagRequired("bcc-file-path")
	if err != nil {
		log.Errorf("set required flag error: %v", err)
	}
	loadCmd.Flags().StringVarP(&wasmFilePath, "wasm-file-path", "w", wasmFilePath, `The file path of wasm code.`)
	err = loadCmd.MarkFlagRequired("wasm-file-path")
	if err != nil {
		log.Errorf("set required flag error: %v", err)
	}
	loadCmd.Flags().StringVarP(&dbFilePath, "output", "o", dbFilePath, `The file path of SQLite db.`)
	err = loadCmd.MarkFlagRequired("output")
	if err != nil {
		log.Errorf("set required flag error: %v", err)
	}
}
