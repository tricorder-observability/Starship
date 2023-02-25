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

	"github.com/stretchr/testify/assert"

	bazelutils "github.com/tricorder/src/testing/bazel"
	"github.com/tricorder/src/utils/uuid"
)

// test module dao fun
// init sqlit gorm and create table
// test dao.SaveCode and check save result
// test dao.QueryByID
// test update code status and check update result
func TestNodeAgent(t *testing.T) {
	assert := assert.New(t)

	dirPath := bazelutils.CreateTmpDir()
	defer func() {
		assert.Nil(os.RemoveAll(dirPath))
	}()

	sqliteClient, _ := InitSqlite(dirPath)

	nodeAgentDao := NodeAgentDao{
		Client: sqliteClient,
	}

	id := strings.Replace(uuid.New(), "-", "_", -1)
	node := &NodeAgentGORM{
		AgentID:  id,
		NodeName: "TestNodeAgent",
	}

	// save module
	err := nodeAgentDao.SaveAgent(node)
	if err != nil {
		t.Errorf("save node agent err %v", err)
	}

	// test queryByID
	node, err = nodeAgentDao.QueryByID(id)
	if err != nil {
		t.Errorf("not query ID=%s data, save node agent err %v", id, err)
	}
	if node.AgentID != id {
		t.Errorf("save node error, node.ID !=  " + id)
	}

	// if node.NodeName != TestNodeAgent, node save error
	if node.NodeName != "TestNodeAgent" {
		t.Errorf("save agent error, node.Name != TestNodeAgent")
	}

	createTime := *node.CreateTime
	lastUpdateTime := *node.LastUpdateTime

	// update ID
	newID := strings.Replace(uuid.New(), "-", "_", -1)
	node.AgentID = newID
	err = nodeAgentDao.UpdateByName(node)
	if err != nil {
		t.Errorf("update node error: %v", err)
	}

	node, err = nodeAgentDao.QueryByName(node.NodeName)
	if err != nil {
		t.Errorf("query node by Name error: %v", err)
	}
	// check update id result
	if node.AgentID != newID {
		t.Errorf("update node.AgentID=newID error")
	}

	if *node.LastUpdateTime == lastUpdateTime {
		t.Errorf("update node.AgentID=newID error, LastUpdateTime not update")
	}

	if *node.CreateTime != createTime {
		t.Errorf("update node.AgentID=newID error, can not update CreateTime")
	}

	createTime = *node.CreateTime
	lastUpdateTime = *node.LastUpdateTime

	// test node.Status
	if node.State != int(AgentStateOnline) {
		t.Errorf("query node state error, node.Status != AgentStatusOnline ")
	}

	// test update module status
	err = nodeAgentDao.UpdateStatusByName(node.NodeName, int(AgentStateOffline))
	if err != nil {
		t.Errorf("change node state error: %v", err)
	}
	node, err = nodeAgentDao.QueryByID(node.AgentID)
	if err != nil {
		t.Errorf("query node by AgentID error: %v", err)
	}
	// check node status
	if node.State != int(AgentStateOffline) {
		t.Errorf("change node status by AgentID error: not change node status")
	}

	if *node.LastUpdateTime == lastUpdateTime {
		t.Errorf("change node status by AgentID error, LastUpdateTime not update")
	}

	if *node.CreateTime != createTime {
		t.Errorf("change node status by AgentID error, can not update CreateTime")
	}

	// get module list *
	list, err := nodeAgentDao.List("*")
	if err != nil {
		t.Errorf("query node list error: %v", err)
	}
	if len(list) == 0 {
		t.Errorf("query node list error: not found node data")
	}

	if list[0].AgentID != node.AgentID {
		t.Errorf("query node list erro default: not found inserted node")
	}

	// get node list default
	list, err = nodeAgentDao.List()
	if err != nil {
		t.Errorf("query module list default error: %v", err)
	}
	if len(list) == 0 {
		t.Errorf("query module list erro default: not found node data")
	}
	if list[0].AgentID != node.AgentID {
		t.Errorf("query module list erro default: not found inserted node")
	}

	// get module list default
	list, err = nodeAgentDao.List("agent_id", "node_name")
	if err != nil {
		t.Errorf("query node list default error: %v", err)
	}
	if len(list) == 0 {
		t.Errorf("query node list erro default: not found node data")
	}
	if len(list[0].AgentID) == 0 {
		t.Errorf("query node list erro default: AgentID is empty")
	}
	if len(list[0].NodeName) == 0 {
		t.Errorf("query node list erro default: NodeName is empty")
	}
}
