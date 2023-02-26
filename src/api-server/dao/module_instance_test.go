// Copyright (C) 2023  Tricorder Observability

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package dao

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	bazelutils "github.com/tricorder/src/testing/bazel"
	"github.com/tricorder/src/utils/uuid"
)

// test module dao fun
// init sqlit gorm and create table
// test dao.SaveCode and check save result
// test dao.QueryByID
// test update code status and check update result
func TestModuleInstance(t *testing.T) {
	assert := assert.New(t)

	dirPath := bazelutils.CreateTmpDir()
	defer func() {
		assert.Nil(os.RemoveAll(dirPath))
	}()

	sqliteClient, _ := InitSqlite(dirPath)

	ModuleInstanceDao := ModuleInstanceDao{
		Client: sqliteClient,
	}

	id := strings.Replace(uuid.New(), "-", "_", -1)
	moduleInstance := &ModuleInstanceGORM{
		ID:         id,
		ModuleID:   id,
		ModuleName: "TestModule",
		AgentID:    id,
		NodeName:   "TestNodeAgent",
		State:      ModuleInstanceInit,
	}

	// save module
	err := ModuleInstanceDao.SaveModuleInstance(moduleInstance)
	assert.Nil(err, "save module instance err %v", err)

	// test queryByID
	moduleInstance, err = ModuleInstanceDao.QueryByID(id)
	assert.Nil(err, "not query ID=%s data, save module instance err %v", id, err)
	assert.Equal(moduleInstance.AgentID, id, "save module instance error, moduleInstance.ID !=  "+id)
	assert.Equal(moduleInstance.NodeName, "TestNodeAgent", "save module instance error, moduleInstance.Name != TestNodeAgent")

	createTime := *moduleInstance.CreateTime
	lastUpdateTime := *moduleInstance.LastUpdateTime

	// update ID
	newID := strings.Replace(uuid.New(), "-", "_", -1)
	moduleInstance.AgentID = newID
	err = ModuleInstanceDao.UpdateByID(moduleInstance)
	assert.Nil(err, "update module instance error: %v", err)

	moduleInstance, err = ModuleInstanceDao.QueryByID(id)
	assert.Nil(err, "query module instance by ID error: %v", err)
	assert.Equal(moduleInstance.AgentID, newID, "update moduleInstance.AgentID=newID error")
	assert.NotEqual(*moduleInstance.LastUpdateTime, lastUpdateTime, "update moduleInstance.AgentID=newID error, LastUpdateTime not update")
	assert.Equal(*moduleInstance.CreateTime, createTime, "update moduleInstance.AgentID=newID error, can not update CreateTime")

	createTime = *moduleInstance.CreateTime
	lastUpdateTime = *moduleInstance.LastUpdateTime

	// test node.Status
	assert.Equal(moduleInstance.State, int(ModuleInstanceInit), "query moduleInstance state error, moduleInstance.Status != ModuleInstanceInit ")

	// test update module status
	err = ModuleInstanceDao.UpdateStatusByID(moduleInstance.ID, int(ModuleInstanceSucceeed))
	assert.Nil(err, "update moduleInstance status by ID error: %v", err)

	moduleInstance, err = ModuleInstanceDao.QueryByID(moduleInstance.ID)
	assert.Nil(err, "query moduleInstance by ID error: %v", err)

	// check node status
	assert.Equal(moduleInstance.State, int(ModuleInstanceSucceeed), "change moduleInstance status by ID error: not change moduleInstance status")
	assert.NotEqual(*moduleInstance.LastUpdateTime, lastUpdateTime, "change moduleInstance status by ID error, LastUpdateTime not update")
	assert.Equal(*moduleInstance.CreateTime, createTime, "change moduleInstance status by ID error, can not update CreateTime")

	// get moduleInstance list *
	list, err := ModuleInstanceDao.List("*")
	assert.Nil(err, "query moduleInstance list error: %v", err)
	assert.NotEqual(len(list), 0, "query moduleInstance list error: not found moduleInstance data")
	assert.Equal(list[0].ID, moduleInstance.ID, "query moduleInstance list erro default: not found inserted moduleInstance")

	// get moduleInstance list default
	list, err = ModuleInstanceDao.List()
	assert.Nil(err, "query moduleInstance list default error: %v", err)
	assert.NotEqual(len(list), 0, "query moduleInstance list erro default: not found moduleInstance data")
	assert.Equal(list[0].ID, moduleInstance.ID, "query moduleInstance list erro default: not found inserted moduleInstance")

	// get moduleInstance list
	list, err = ModuleInstanceDao.List("agent_id", "node_name")
	assert.Nil(err, "query moduleInstance list error: %v", err)
	assert.NotEqual(len(list), 0, "query moduleInstance list error: not found moduleInstance data")
	assert.NotEqual(len(list[0].AgentID), 0, "query moduleInstance list erro default: AgentID is empty")
	assert.NotEqual(len(list[0].NodeName), 0, "query moduleInstance list erro default: NodeName is empty")

	// get moduleInstance list By State
	list, err = ModuleInstanceDao.ListByState(ModuleInstanceSucceeed)
	assert.Nil(err, "query moduleInstance list by state error: %v", err)
	assert.NotEqual(len(list), 0, "query moduleInstance list by state error: not found moduleInstance data")
	assert.Equal(list[0].ID, moduleInstance.ID, "query moduleInstance list by state erro default: not found inserted moduleInstance")

	// get moduleInstance list By AgentID
	list, err = ModuleInstanceDao.ListByAgentID(moduleInstance.AgentID)
	assert.Nil(err, "query moduleInstance list by agentID error: %v", err)
	assert.NotEqual(len(list), 0, "query moduleInstance list by agentID error: not found moduleInstance data")
	assert.Equal(list[0].ID, moduleInstance.ID, "query moduleInstance list by agentID erro default: not found inserted moduleInstance")

	// get moduleInstance list By ModuleID
	list, err = ModuleInstanceDao.ListByModuleID(moduleInstance.ModuleID)
	assert.Nil(err, "query moduleInstance list by moduleID error: %v", err)
	assert.NotEqual(len(list), 0, "query moduleInstance list by moduleID error: not found moduleInstance data")
	assert.Equal(list[0].ID, moduleInstance.ID, "query moduleInstance list by moduleID erro default: not found inserted moduleInstance")

	// get moduleInstance list By nodeName
	list, err = ModuleInstanceDao.ListByNodeName(moduleInstance.NodeName)
	assert.Nil(err, "query moduleInstance list by nodeName error: %v", err)
	assert.NotEqual(len(list), 0, "query moduleInstance list by nodeName error: not found moduleInstance data")
	assert.Equal(list[0].ID, moduleInstance.ID, "query moduleInstance list by nodeName erro default: not found inserted moduleInstance")
}
