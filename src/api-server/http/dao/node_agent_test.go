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
	"github.com/stretchr/testify/require"

	pb "github.com/tricorder/src/api-server/pb"
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
	require := require.New(t)

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
	require.Nil(err, "save node agent err %v", err)

	node.NodeName = "TestNodeAgent2"
	err = nodeAgentDao.SaveAgent(node)
	require.Nil(err, "save node agent upsert err %v", err)
	nodeRes, err := nodeAgentDao.QueryByID(node.AgentID)
	assert.Nil(err, "not query ID=%s data, save upsert node agent err %v", id, err)
	assert.Equal(nodeRes.NodeName, "TestNodeAgent2", "save upsert node agent err %v", err)

	node.NodeName = "TestNodeAgent"
	err = nodeAgentDao.SaveAgent(node)
	require.Nil(err, "save node agent upsert err %v", err)
	nodeRes, err = nodeAgentDao.QueryByID(id)
	assert.Equal(nodeRes.NodeName, "TestNodeAgent", "save upsert node agent err %v", err)

	// test queryByID
	node, err = nodeAgentDao.QueryByID(id)
	require.Nil(err, "not query ID=%s data, save node agent err %v", id, err)
	assert.Equal(node.AgentID, id, "save node error, node.ID !=  "+id)
	assert.Equal(node.NodeName, "TestNodeAgent", "save agent error, node.Name != TestNodeAgent")

	createTime := *node.CreateTime
	lastUpdateTime := *node.LastUpdateTime

	// update
	newNodeName := "NewTestNodeAgent"
	node.NodeName = newNodeName
	err = nodeAgentDao.UpdateByID(node)
	assert.Nil(err, "update node error: %v", err)

	node, err = nodeAgentDao.QueryByID(node.AgentID)
	require.Nil(err, "not query ID=%s data, save node agent err %v", id, err)
	assert.Equal(node.NodeName, newNodeName, "update agent error, node.Name != ", newNodeName)
	assert.NotEqual(*node.LastUpdateTime, lastUpdateTime, "update node error, LastUpdateTime not update")
	assert.Equal(*node.CreateTime, createTime, "update node error, can not update CreateTime")

	createTime = *node.CreateTime
	lastUpdateTime = *node.LastUpdateTime

	// test node.Status
	assert.Equal(node.State, int(pb.AgentState_ONLINE), "query node state error, node.Status != pb.AgentState_ONLINE ")

	// test update module status
	err = nodeAgentDao.UpdateStateByID(node.AgentID, int(pb.AgentState_OFFLINE))
	require.Nil(err, "change node state error: %v", err)

	node, err = nodeAgentDao.QueryByID(node.AgentID)
	require.Nil(err, "not query ID=%s data, save node agent err %v", id, err)
	assert.Equal(node.State, int(pb.AgentState_OFFLINE), "change node state error, node.Status != pb.AgentState_OFFLINE ")
	assert.NotEqual(*node.LastUpdateTime, lastUpdateTime, "change node state error, LastUpdateTime not update")
	assert.Equal(*node.CreateTime, createTime, "change node state error, can not update CreateTime")
	assert.Equal(node.NodeName, newNodeName, "change node state error, node.NodeName !=  "+newNodeName)

	// get module list *
	list, err := nodeAgentDao.List([]string{"*"})
	require.Nil(err, "query node list error: %v", err)
	assert.NotEqual(len(list), 0, "query node list error: not found node data")
	assert.Equal(list[0].AgentID, node.AgentID, "query node list erro default: not found inserted node")

	// get node list default
	list, err = nodeAgentDao.List([]string{})
	require.Nil(err, "query node list default error: %v", err)
	assert.NotEqual(len(list), 0, "query node list erro default: not found node data")
	assert.Equal(list[0].AgentID, node.AgentID, "query node list erro default: not found inserted node")

	// get module list default
	list, err = nodeAgentDao.List([]string{"agent_id", "node_name"})
	require.Nil(err, "query node list default error: %v", err)
	assert.NotEqual(len(list), 0, "query node list erro default: not found node data")
	assert.Equal(list[0].AgentID, node.AgentID, "query node list erro default: not found inserted node")
	assert.NotEqual(len(list[0].AgentID), 0, "query node list erro default: AgentID is not empty")
	assert.NotEqual(len(list[0].NodeName), 0, "query node list erro default: NodeName is not empty")

	// get module list by state
	list, err = nodeAgentDao.ListByState(int(pb.AgentState_OFFLINE))
	require.Nil(err, "query node ListByState error: %v", err)
	assert.NotEqual(len(list), 0, "query node ListByState error: not found node data")
	assert.Equal(list[0].AgentID, node.AgentID, "query node ListByState error: not found inserted node")
	assert.NotEqual(len(list[0].AgentID), 0, "query node ListByState error: AgentID is not empty")
	assert.NotEqual(len(list[0].NodeName), 0, "query node ListByState error: NodeName is not empty")

	// get module list by node name
	list, err = nodeAgentDao.ListByNodeName(newNodeName)
	require.Nil(err, "query node ListByName error: %v", err)
	assert.NotEqual(len(list), 0, "query node ListByName error: not found node data")
	assert.Equal(list[0].AgentID, node.AgentID, "query node ListByName error: not found inserted node")
	assert.NotEqual(len(list[0].AgentID), 0, "query node ListByName error: AgentID is not empty")
	assert.NotEqual(len(list[0].NodeName), 0, "query node ListByName error: NodeName is not empty")
}
