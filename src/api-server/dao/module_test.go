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

package dao

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	pb "github.com/tricorder/src/api-server/pb"
	bazelutils "github.com/tricorder/src/testing/bazel"
	"github.com/tricorder/src/utils/uuid"
)

// Save test data in this path
const SQLiteFilePath string = "module_test"

// test module dao fun
// init sqlit gorm and create table
// test dao.SaveModule and check save result
// test dao.QueryBiID and dao.QueryByWasmID
// test update module status and check update result
func TestModule(t *testing.T) {
	assert := assert.New(t)

	dirPath := bazelutils.CreateTmpDir()
	defer func() {
		assert.Nil(os.RemoveAll(dirPath))
	}()

	sqliteClient, _ := InitSqlite(dirPath)

	moduleDao := ModuleDao{
		Client: sqliteClient,
	}

	id := strings.Replace(uuid.New(), "-", "_", -1)
	module := &ModuleGORM{
		ID:                 id,
		DesireState:        int(pb.DeploymentState_CREATED),
		Name:               "TestModule",
		Wasm:               []byte("WasmUid"),
		CreateTime:         time.Now().Format("2006-01-02 15:04:05"),
		EbpfPerfBufferName: "events",
	}
	// save module
	err := moduleDao.SaveModule(module)
	if err != nil {
		t.Errorf("save module err %v", err)
	}
	// test queryByID

	module, err = moduleDao.QueryByID(id)
	if err != nil {
		t.Errorf("not query ID=%s data, save module err %v", id, err)
	}
	if module.ID != id {
		t.Errorf("save module error, module.ID !=  " + id)
	}

	// if module.Name != TestModule, module save error
	if module.Name != "TestModule" {
		t.Errorf("save module error, module.Name != TestModule ")
	}

	// update status
	module.Name = "UpdateName"
	err = moduleDao.UpdateByID(module)
	if err != nil {
		t.Errorf("update module error: %v", err)
	}
	module, err = moduleDao.QueryByID(module.ID)
	if err != nil {
		t.Errorf("query module by ID error: %v", err)
	}
	// check update name result
	if module.Name != "UpdateName" {
		t.Errorf("update module.Name=UpdateName error")
	}

	// test module.DesireState
	if module.DesireState != int(pb.DeploymentState_CREATED) {
		t.Errorf("query module status error, module.DesireState != DeploymentState_CREATED ")
	}

	// test update module status
	err = moduleDao.UpdateStatusByID(module.ID, int(pb.DeploymentState_TO_BE_DEPLOYED))
	if err != nil {
		t.Errorf("change module status error: %v", err)
	}
	module, err = moduleDao.QueryByID(module.ID)
	if err != nil {
		t.Errorf("query module by ID error: %v", err)
	}
	// check module DesireState
	if module.DesireState != int(pb.DeploymentState_TO_BE_DEPLOYED) {
		t.Errorf("change module status by ID error: not change module status")
	}
	// get module list *
	list, err := moduleDao.ListModule("*")
	if err != nil {
		t.Errorf("query module list error: %v", err)
	}
	if len(list) == 0 {
		t.Errorf("query module list error: not found module data")
	}

	if len(list[0].Wasm) == 0 {
		t.Errorf("query module list erro default: not found wasm data")
	}

	// get module list default
	list, err = moduleDao.ListModule()
	if err != nil {
		t.Errorf("query module list default error: %v", err)
	}
	if len(list) == 0 {
		t.Errorf("query module list erro default: not found module data")
	}
	if len(list[0].Wasm) != 0 {
		t.Errorf("query module list erro default: not found wasm data")
	}

	// get module list default
	list, err = moduleDao.ListModule("id", "name")
	if err != nil {
		t.Errorf("query module list default error: %v", err)
	}
	if len(list) == 0 {
		t.Errorf("query module list erro default: not found module data")
	}
	if len(list[0].ID) == 0 {
		t.Errorf("query module list erro default: ID is empty")
	}
	if len(list[0].Name) == 0 {
		t.Errorf("query module list erro default: Name is empty")
	}
	if len(list[0].Wasm) != 0 {
		t.Errorf("query module list erro default: Wasm is not empty")
	}
}
