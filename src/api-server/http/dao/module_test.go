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
		DesireState:        int(pb.ModuleState_CREATED_),
		Name:               "TestModule",
		Wasm:               []byte("WasmUid"),
		CreateTime:         time.Now().Format("2006-01-02 15:04:05"),
		EbpfPerfBufferName: "events",
	}
	// save module
	err := moduleDao.SaveModule(module)
	assert.Nil(err, "save module err %v", err)

	module.Name = "TestModule2"
	err = moduleDao.SaveModule(module)
	assert.Nil(err, "save module upsert err %v", err)
	moduleRes, err := moduleDao.QueryByID(module.ID)
	assert.Nil(err, "not query ID=%s data, save module err %v", id, err)
	assert.Equal(moduleRes.Name, "TestModule2", "save module error, module.Name != TestModule2 ")

	module.Name = "TestModule"
	err = moduleDao.SaveModule(module)
	assert.Nil(err, "save module upsert err %v", err)
	moduleRes, err = moduleDao.QueryByID(module.ID)
	assert.Nil(err, "not query ID=%s data, save module err %v", id, err)
	assert.Equal(moduleRes.Name, "TestModule", "save module error, module.Name != TestModule ")

	// test queryByID
	module, err = moduleDao.QueryByID(id)
	assert.Nil(err, "not query ID=%s data, save module err %v", id, err)
	assert.Equal(module.ID, id, "save module error, module.ID !=  "+id)

	// if module.Name != TestModule, module save error
	assert.Equal(module.Name, "TestModule", "save module error, module.Name != TestModule ")

	// update status
	module.Name = "UpdateName"
	err = moduleDao.UpdateByID(module)
	assert.Nil(err, "update module error: %v", err)

	module, err = moduleDao.QueryByID(module.ID)
	assert.Nil(err, "query module by ID error: %v", err)
	assert.Equal(module.Name, "UpdateName", "update module.Name=UpdateName error")

	// test module.DesireState
	assert.Equal(module.DesireState, int(pb.ModuleState_CREATED_),
		"query module status error, module.DesireState != ModuleState_CREATED_ ")

	// test update module status
	err = moduleDao.UpdateStatusByID(module.ID, int(pb.ModuleState_DEPLOYED))
	assert.Nil(err, "change module status error: %v", err)

	module, err = moduleDao.QueryByID(module.ID)
	assert.Nil(err, "query module by ID error: %v", err)

	// check module DesireState
	assert.Equal(module.DesireState, int(pb.ModuleState_DEPLOYED),
		"change module status error: not change module status")

	// get module list *
	list, err := moduleDao.ListModule("*")
	assert.Nil(err, "query module list error: %v", err)
	assert.NotEqual(len(list), 0, "query module list error: not found module data")
	assert.NotEqual(len(list[0].Wasm), 0, "query module list error: not found wasm data")

	// get module list default
	list, err = moduleDao.ListModule()
	assert.Nil(err, "query module list default error: %v", err)
	assert.NotEqual(len(list), 0, "query module list erro default: not found module data")
	assert.Equal(len(list[0].Wasm), 0, "query module list erro default: not found wasm data")

	// get module list default
	list, err = moduleDao.ListModule("id", "name")
	assert.Nil(err, "query module list default error: %v", err)
	assert.NotEqual(len(list), 0, "query module list erro default: not found module data")
	assert.NotEqual(len(list[0].ID), 0, "query module list erro default: ID is empty")
	assert.NotEqual(len(list[0].Name), 0, "query module list erro default: Name is empty")
	assert.Equal(len(list[0].Wasm), 0, "query module list erro default: Wasm is not empty")
}
