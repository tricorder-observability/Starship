package load

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

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
		err = loadModule(moduleReq)
		if err != nil {
			log.Error(err)
		}
	},
}

func init() {
	loadCmd.Flags().StringVarP(&moduleFilePath, "module-file-path", "m",
		moduleFilePath, "The file path of module in json format. If missing a flag would cause command to fail, mark this flag as required")
	loadCmd.Flags().StringVarP(&bccFilePath, "bcc-file-path", "b", bccFilePath, "The file path of bcc code. If missing a flag would cause command to fail, mark this flag as required")
	loadCmd.Flags().StringVarP(&wasmFilePath, "wasm-file-path", "w", wasmFilePath, "The file path of wasm code. If missing a flag would cause command to fail, mark this flag as required")
	loadCmd.Flags().StringVarP(&dbFilePath, "output", "o", dbFilePath, "The file path of SQLite db. If missing a flag would cause command to fail, mark this flag as required")
}
