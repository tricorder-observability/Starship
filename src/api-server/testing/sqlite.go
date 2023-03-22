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

package testing

import (
	"time"

	"github.com/tricorder/src/utils/log"
	"github.com/tricorder/src/utils/sqlite"

	pb "github.com/tricorder/src/api-server/pb"

	"github.com/tricorder/src/api-server/http/dao"
)

var ebpfJson = `
#include <linux/ptrace.h>

BPF_PERF_OUTPUT(events);

// Writes a fixed JSON string to perf buffer.
int sample_json(struct bpf_perf_event_data *ctx) {
  const char word[] = "{\"name\":\"John\", \"age\":30}";
  events.perf_submit(ctx, (void *)word, sizeof(word));
  return 0;
}
`

// PrepareTricorderDBData writes test data into a testing database.
func PrepareTricorderDBData(moduleID string, agentID string, moduleInstanceID string, sqliteClient *sqlite.ORM) {
	moduleDao := dao.ModuleDao{
		Client: sqliteClient,
	}

	moduleInstanceDao := dao.ModuleInstanceDao{
		Client: sqliteClient,
	}

	module := &dao.ModuleGORM{
		ID:                 moduleID,
		Ebpf:               ebpfJson,
		Wasm:               []byte("moduleString"),
		CreateTime:         time.Date(2022, 12, 31, 14, 30, 0, 0, time.Local).Format("2006-01-02 15:04:05"),
		DesireState:        int(pb.ModuleState_DEPLOYED),
		Name:               "test-module-foo",
		EbpfFmt:            0,
		EbpfLang:           0,
		EbpfPerfBufferName: "events",

		SchemaName: "out_put_name",
		SchemaAttr: "[{\"name\":\"data\",\"type\":5}]",
		Fn:         "copy_input_to_output",
		WasmFmt:    0,
		WasmLang:   0,
	}

	err := moduleDao.SaveModule(module)
	if err != nil {
		log.Fatalf("While writing data to database for testing, failed to save module data, error: %v", err)
	}

	moduleInstance := &dao.ModuleInstanceGORM{
		ID:          moduleInstanceID,
		ModuleID:    moduleID,
		ModuleName:  module.Name,
		AgentID:     agentID,
		DesireState: module.DesireState,
		State:       int(pb.ModuleInstanceState_INIT),
	}
	err = moduleInstanceDao.SaveModuleInstance(moduleInstance)
	if err != nil {
		log.Fatalf("While writing data to database for testing, failed to save module instance data, error: %v", err)
	}
}
