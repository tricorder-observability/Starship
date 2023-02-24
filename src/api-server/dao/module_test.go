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
const SQLiteFilePath string = "code_test"

// test code dao fun
// init sqlit gorm and create table
// test dao.SaveCode and check save result
// test dao.QueryBiID and dao.QueryByWasmID
// test update code status and check update result
func TestModule(t *testing.T) {
	assert := assert.New(t)

	dirPath := bazelutils.CreateTmpDir()
	defer func() {
		assert.Nil(os.RemoveAll(dirPath))
	}()

	sqliteClient, _ := InitSqlite(dirPath)

	codeDao := Module{
		Client: sqliteClient,
	}

	id := strings.Replace(uuid.New(), "-", "_", -1)
	code := &ModuleGORM{
		ID:                 id,
		Status:             int(pb.DeploymentState_CREATED),
		Name:               "TestCode",
		Wasm:               []byte("WasmUid"),
		CreateTime:         time.Now().Format("2006-01-02 15:04:05"),
		EbpfPerfBufferName: "events",
	}
	// save code module
	err := codeDao.SaveCode(code)
	if err != nil {
		t.Errorf("save code err %v", err)
	}
	// test queryByID

	code, err = codeDao.QueryByID(id)
	if err != nil {
		t.Errorf("not query ID=%s data, save code err %v", id, err)
	}
	if code.ID != id {
		t.Errorf("save code error, code.ID !=  " + id)
	}

	// if code.Name != TestCode, code save error
	if code.Name != "TestCode" {
		t.Errorf("save code error, code.Name != TestCode ")
	}

	// update status
	code.Name = "UpdateName"
	err = codeDao.UpdateByID(code)
	if err != nil {
		t.Errorf("update code error: %v", err)
	}
	code, err = codeDao.QueryByID(code.ID)
	if err != nil {
		t.Errorf("query code by ID error: %v", err)
	}
	// check update name result
	if code.Name != "UpdateName" {
		t.Errorf("update code.Name=UpdateName error")
	}

	// test code.Status
	if code.Status != int(pb.DeploymentState_CREATED) {
		t.Errorf("query code status error, code.Status != DeploymentState_CREATED ")
	}

	// test update code status
	err = codeDao.UpdateStatusByID(code.ID, int(pb.DeploymentState_TO_BE_DEPLOYED))
	if err != nil {
		t.Errorf("change code status error: %v", err)
	}
	code, err = codeDao.QueryByID(code.ID)
	if err != nil {
		t.Errorf("query code by ID error: %v", err)
	}
	// check code status
	if code.Status != int(pb.DeploymentState_TO_BE_DEPLOYED) {
		t.Errorf("change code status by ID error: not change code status")
	}
	// get code list *
	list, err := codeDao.ListCode("*")
	if err != nil {
		t.Errorf("query code list error: %v", err)
	}
	if len(list) == 0 {
		t.Errorf("query code list error: not found code data")
	}

	if len(list[0].Wasm) == 0 {
		t.Errorf("query code list erro default: not found wasm data")
	}

	// get code list default
	list, err = codeDao.ListCode()
	if err != nil {
		t.Errorf("query code list default error: %v", err)
	}
	if len(list) == 0 {
		t.Errorf("query code list erro default: not found code data")
	}
	if len(list[0].Wasm) != 0 {
		t.Errorf("query code list erro default: not found wasm data")
	}

	// get code list default
	list, err = codeDao.ListCode("id", "name")
	if err != nil {
		t.Errorf("query code list default error: %v", err)
	}
	if len(list) == 0 {
		t.Errorf("query code list erro default: not found code data")
	}
	if len(list[0].ID) == 0 {
		t.Errorf("query code list erro default: ID is empty")
	}
	if len(list[0].Name) == 0 {
		t.Errorf("query code list erro default: Name is empty")
	}
	if len(list[0].Wasm) != 0 {
		t.Errorf("query code list erro default: Wasm is not empty")
	}
}
