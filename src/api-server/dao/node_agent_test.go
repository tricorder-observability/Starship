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
		ID:   id,
		Name: "TestNodeAgent",
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
	if node.ID != id {
		t.Errorf("save node error, node.ID !=  " + id)
	}

	// if node.Name != TestNodeAgent, node save error
	if node.Name != "TestNodeAgent" {
		t.Errorf("save agent error, node.Name != TestNodeAgent")
	}

	listA, _ := nodeAgentDao.List()
	t.Errorf("all list %v", listA)

	// update ID
	newID := strings.Replace(uuid.New(), "-", "_", -1)
	node.ID = newID
	err = nodeAgentDao.UpdateByName(node)
	if err != nil {
		t.Errorf("update node error: %v", err)
	}
	node, err = nodeAgentDao.QueryByName(node.Name)
	if err != nil {
		t.Errorf("query node by Name error: %v", err)
	}
	// check update id result
	if node.ID != newID {
		t.Errorf("update node.ID=newID error")
	}

	// test node.Status
	if node.Status != int(AgentStatusOnline) {
		t.Errorf("query node status error, node.Status != AgentStatusOnline ")
	}

	// test update module status
	err = nodeAgentDao.UpdateStatusByName(node.Name, int(AgentStatusOffline))
	if err != nil {
		t.Errorf("change node status error: %v", err)
	}
	node, err = nodeAgentDao.QueryByID(node.ID)
	if err != nil {
		t.Errorf("query node by ID error: %v", err)
	}
	// check node status
	if node.Status != int(AgentStatusOffline) {
		t.Errorf("change node status by ID error: not change node status")
	}
	// get module list *
	list, err := nodeAgentDao.List("*")
	if err != nil {
		t.Errorf("query node list error: %v", err)
	}
	if len(list) == 0 {
		t.Errorf("query node list error: not found node data")
	}

	if list[0].ID != node.ID {
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
	if list[0].ID != node.ID {
		t.Errorf("query module list erro default: not found inserted node")
	}

	// get module list default
	list, err = nodeAgentDao.List("id", "name")
	if err != nil {
		t.Errorf("query node list default error: %v", err)
	}
	if len(list) == 0 {
		t.Errorf("query node list erro default: not found node data")
	}
	if len(list[0].ID) == 0 {
		t.Errorf("query node list erro default: ID is empty")
	}
	if len(list[0].Name) == 0 {
		t.Errorf("query node list erro default: Name is empty")
	}
}
