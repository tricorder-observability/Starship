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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	pb "github.com/tricorder/src/api-server/pb"
	bazelutils "github.com/tricorder/src/testing/bazel"
	"github.com/tricorder/src/utils/uuid"
)

// test module dao fun
// init sqlit gorm and create table
// test dao.SaveModuleInstance and check save result
// test dao.QueryByID
// test update code status and check update result
func TestModuleInstance(t *testing.T) {
	assert := assert.New(t)

	dirPath := bazelutils.CreateTmpDir()
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
		State:      int(pb.ModuleInstanceState_INIT),
	}

	// save module
	err := ModuleInstanceDao.SaveModuleInstance(moduleInstance)
	assert.Nil(err, "save module instance err %v", err)

	moduleInstance.NodeName = "TestNodeAgent2"
	err = ModuleInstanceDao.SaveModuleInstance(moduleInstance)
	assert.Nil(err, "save module instance upsert err %v", err)
	moduleRes, err := ModuleInstanceDao.QueryByID(id)
	assert.Nil(err, "not query ID=%s data, save module instance err %v", id, err)
	assert.Equal(moduleRes.NodeName, "TestNodeAgent2",
		"save module instance error, moduleInstance.Name != TestNodeAgent2")

	moduleInstance.NodeName = "TestNodeAgent"
	err = ModuleInstanceDao.SaveModuleInstance(moduleInstance)
	assert.Nil(err, "save module instance upsert err %v", err)
	moduleRes, err = ModuleInstanceDao.QueryByID(id)
	assert.Nil(err, "not query ID=%s data, save module instance err %v", id, err)
	assert.Equal(moduleRes.NodeName, "TestNodeAgent", "save module instance error, moduleInstance.Name != TestNodeAgent2")

	// test queryByID
	moduleInstance, err = ModuleInstanceDao.QueryByID(id)
	assert.Nil(err, "not query ID=%s data, save module instance err %v", id, err)
	assert.Equal(moduleInstance.AgentID, id, "save module instance error, moduleInstance.ID !=  "+id)
	assert.Equal(moduleInstance.NodeName, "TestNodeAgent",
		"save module instance error, moduleInstance.Name != TestNodeAgent")

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
	assert.NotEqual(*moduleInstance.LastUpdateTime, lastUpdateTime,
		"update moduleInstance.AgentID=newID error, LastUpdateTime not update")
	assert.Equal(*moduleInstance.CreateTime, createTime,
		"update moduleInstance.AgentID=newID error, can not update CreateTime")

	createTime = *moduleInstance.CreateTime
	lastUpdateTime = *moduleInstance.LastUpdateTime

	// test node.Status
	assert.Equal(moduleInstance.State, int(pb.ModuleInstanceState_INIT),
		"query moduleInstance state error, moduleInstance.Status != pb.ModuleInstanceState_INIT ")

	// test update module status
	err = ModuleInstanceDao.UpdateStatusByID(moduleInstance.ID, int(pb.ModuleInstanceState_SUCCEEDED))
	assert.Nil(err, "update moduleInstance status by ID error: %v", err)

	moduleInstance, err = ModuleInstanceDao.QueryByID(moduleInstance.ID)
	assert.Nil(err, "query moduleInstance by ID error: %v", err)

	// check node status
	assert.Equal(moduleInstance.State, int(pb.ModuleInstanceState_SUCCEEDED),
		"change moduleInstance status by ID error: not change moduleInstance status")
	assert.NotEqual(*moduleInstance.LastUpdateTime, lastUpdateTime,
		"change moduleInstance status by ID error, LastUpdateTime not update")
	assert.Equal(*moduleInstance.CreateTime, createTime,
		"change moduleInstance status by ID error, can not update CreateTime")

	// get moduleInstance list *
	list, err := ModuleInstanceDao.List("*")
	assert.Nil(err, "query moduleInstance list error: %v", err)
	assert.NotEqual(len(list), 0,
		"query moduleInstance list error: not found moduleInstance data")
	assert.Equal(list[0].ID, moduleInstance.ID,
		"query moduleInstance list erro default: not found inserted moduleInstance")

	// get moduleInstance list default
	list, err = ModuleInstanceDao.List()
	assert.Nil(err, "query moduleInstance list default error: %v", err)
	assert.NotEqual(len(list), 0,
		"query moduleInstance list erro default: not found moduleInstance data")
	assert.Equal(list[0].ID, moduleInstance.ID,
		"query moduleInstance list erro default: not found inserted moduleInstance")

	// get moduleInstance list
	list, err = ModuleInstanceDao.List("agent_id", "node_name")
	assert.Nil(err, "query moduleInstance list error: %v", err)
	assert.NotEqual(len(list), 0, "query moduleInstance list error: not found moduleInstance data")
	assert.NotEqual(len(list[0].AgentID), 0,
		"query moduleInstance list erro default: AgentID is empty")
	assert.NotEqual(len(list[0].NodeName), 0,
		"query moduleInstance list erro default: NodeName is empty")

	// get moduleInstance list By State
	list, err = ModuleInstanceDao.ListByState(int(pb.ModuleInstanceState_SUCCEEDED))
	assert.Nil(err, "query moduleInstance list by state error: %v", err)
	assert.NotEqual(len(list), 0, "query moduleInstance list by state error: not found moduleInstance data")
	assert.Equal(list[0].ID, moduleInstance.ID,
		"query moduleInstance list by state erro default: not found inserted moduleInstance")

	// get moduleInstance list By AgentID
	list, err = ModuleInstanceDao.ListByAgentID(moduleInstance.AgentID)
	assert.Nil(err, "query moduleInstance list by agentID error: %v", err)
	assert.NotEqual(len(list), 0,
		"query moduleInstance list by agentID error: not found moduleInstance data")
	assert.Equal(list[0].ID, moduleInstance.ID,
		"query moduleInstance list by agentID erro default: not found inserted moduleInstance")

	// get moduleInstance list By ModuleID
	list, err = ModuleInstanceDao.ListByModuleID(moduleInstance.ModuleID)
	assert.Nil(err, "query moduleInstance list by moduleID error: %v", err)
	assert.NotEqual(len(list), 0,
		"query moduleInstance list by moduleID error: not found moduleInstance data")
	assert.Equal(list[0].ID, moduleInstance.ID,
		"query moduleInstance list by moduleID erro default: not found inserted moduleInstance")

	// get moduleInstance list By nodeName
	list, err = ModuleInstanceDao.ListByNodeName(moduleInstance.NodeName)
	assert.Nil(err, "query moduleInstance list by nodeName error: %v", err)
	assert.NotEqual(len(list), 0,
		"query moduleInstance list by nodeName error: not found moduleInstance data")
	assert.Equal(list[0].ID, moduleInstance.ID,
		"query moduleInstance list by nodeName erro default: not found inserted moduleInstance")

	moduleRes, err = ModuleInstanceDao.QueryByAgentIDAndModuleID(moduleInstance.AgentID, moduleInstance.ModuleID)
	assert.Nil(err)
	assert.Equal(moduleRes.ID, moduleInstance.ID)
	assert.Equal(moduleRes.NodeName, moduleInstance.NodeName)
}

// Tests that CheckModuleDesiredState returns expected values.
func TestCheckModuleDesiredState(t *testing.T) {
	assert := assert.New(t)

	dirPath := bazelutils.CreateTmpDir()
	sqliteClient, _ := InitSqlite(dirPath)
	ModuleInstanceDao := ModuleInstanceDao{
		Client: sqliteClient,
	}

	id := strings.Replace(uuid.New(), "-", "_", -1)
	moduleInstance := &ModuleInstanceGORM{
		ID:          "0",
		ModuleID:    id,
		ModuleName:  "TestModule",
		AgentID:     id,
		NodeName:    "TestNodeAgent",
		DesireState: int(pb.ModuleState_DEPLOYED),
		State:       int(pb.ModuleInstanceState_INIT),
	}

	err := ModuleInstanceDao.SaveModuleInstance(moduleInstance)
	assert.Nil(err, "save module instance err %v", err)

	moduleInstance.ID = "1"
	moduleInstance.AgentID = "agent-0"
	moduleInstance.DesireState = int(pb.ModuleState_DEPLOYED)
	err = ModuleInstanceDao.SaveModuleInstance(moduleInstance)
	assert.Nil(err, "save module instance err %v", err)

	isDesiredState, err := ModuleInstanceDao.CheckModuleDesiredState(moduleInstance.ModuleID,
		int(pb.ModuleState_DEPLOYED))
	assert.Nil(err)
	assert.True(isDesiredState)

	isDesiredState, err = ModuleInstanceDao.CheckModuleDesiredState("non-existent-module-id",
		int(pb.ModuleState_DEPLOYED))
	assert.Nil(err)
	// Because there is no instances.
	assert.False(isDesiredState)

	moduleInstance.ID = "2"
	moduleInstance.AgentID = "agent-0"
	moduleInstance.DesireState = int(pb.ModuleState_UNDEPLOYED)
	err = ModuleInstanceDao.SaveModuleInstance(moduleInstance)
	assert.Nil(err, "save module instance err %v", err)

	isDesiredState, err = ModuleInstanceDao.CheckModuleDesiredState(moduleInstance.ModuleID,
		int(pb.ModuleState_DEPLOYED))
	assert.Nil(err)
	assert.False(isDesiredState)
}

// Tests that CheckModuleInProgress returns expected values.
func TestCheckModuleInProgress(t *testing.T) {
	assert := assert.New(t)

	dirPath := bazelutils.CreateTmpDir()
	sqliteClient, _ := InitSqlite(dirPath)
	ModuleInstanceDao := ModuleInstanceDao{
		Client: sqliteClient,
	}

	id := strings.Replace(uuid.New(), "-", "_", -1)
	moduleInstance := &ModuleInstanceGORM{
		ID:          "0",
		ModuleID:    id,
		ModuleName:  "TestModule",
		AgentID:     id,
		NodeName:    "TestNodeAgent",
		DesireState: int(pb.ModuleState_DEPLOYED),
		State:       int(pb.ModuleInstanceState_INIT),
	}

	err := ModuleInstanceDao.SaveModuleInstance(moduleInstance)
	assert.Nil(err, "save module instance err %v", err)

	moduleInstance.ID = "1"
	moduleInstance.AgentID = "agent-0"
	moduleInstance.DesireState = int(pb.ModuleState_DEPLOYED)
	moduleInstance.State = int(pb.ModuleInstanceState_IN_PROGRESS)
	err = ModuleInstanceDao.SaveModuleInstance(moduleInstance)
	assert.Nil(err, "save module instance err %v", err)

	isInProgress, err := ModuleInstanceDao.CheckModuleInProgress(moduleInstance.ModuleID)
	assert.Nil(err)
	assert.True(isInProgress)

	isInProgress, err = ModuleInstanceDao.CheckModuleInProgress("non-existent-module-id")
	assert.Nil(err)
	// Because there is no instances.
	assert.False(isInProgress)

	moduleInstance.ID = "2"
	moduleInstance.ModuleID = "22"
	moduleInstance.AgentID = "agent-0"
	moduleInstance.DesireState = int(pb.ModuleState_UNDEPLOYED)
	moduleInstance.State = int(pb.ModuleInstanceState_SUCCEEDED)
	err = ModuleInstanceDao.SaveModuleInstance(moduleInstance)
	assert.Nil(err, "save module instance err %v", err)

	isInProgress, err = ModuleInstanceDao.CheckModuleInProgress(moduleInstance.ModuleID)
	assert.Nil(err)
	assert.False(isInProgress)
}
