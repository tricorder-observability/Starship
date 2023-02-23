// Copyright (C) 2023 Tricorder Observability
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
	"fmt"
	"time"

	"github.com/tricorder/src/utils/sqlite"
)

const (
	// AgentStatusOnline online
	AgentStatusOnline = 0
	// AgentStatusOffline offline
	AgentStatusOffline = 1
	// AgentStatusTerminated terminated
	AgentStatusTerminated = 2
)

// NodeAgentGORM sqlite gorm storage and response object
type NodeAgentGORM struct {
	Name           string     `gorm:"'name' primarykey" json:"name,omitempty"`
	ID             string     `gorm:"id" json:"id,omitempty"`
	Status         int        `gorm:"status" json:"status,omitempty"`
	CreateTime     *time.Time `gorm:"create_time" json:"create_time,omitempty"`
	LastUpdateTime *time.Time `gorm:"last_update_time" json:"last_update_time,omitempty"`
}

func (NodeAgentGORM) TableName() string {
	return "node_agent"
}

type NodeAgentDao struct {
	Client *sqlite.ORM
}

func (g *NodeAgentDao) SaveAgent(agent *NodeAgentGORM) error {
	agent.LastUpdateTime = &time.Time{}
	*agent.LastUpdateTime = time.Now()
	result := g.Client.Engine.Create(agent)
	return result.Error
}

func (g *NodeAgentDao) UpdateByName(agent *NodeAgentGORM) error {
	if len(agent.Name) == 0 {
		return fmt.Errorf("name is empty")
	}

	agent.LastUpdateTime = &time.Time{}
	*agent.LastUpdateTime = time.Now()
	agent.CreateTime = &time.Time{}
	*agent.CreateTime = time.Now()

	result := g.Client.Engine.Model(&NodeAgentGORM{Name: agent.Name}).Updates(agent)
	return result.Error
}

func (g *NodeAgentDao) UpdateStatusByName(name string, status int) error {
	now := time.Now()
	result := g.Client.Engine.Model(&NodeAgentGORM{Name: name}).Select("status").Updates(NodeAgentGORM{Status: status, LastUpdateTime: &now})
	return result.Error
}

func (g *NodeAgentDao) DeleteByName(name string) error {
	result := g.Client.Engine.Delete(&NodeAgentGORM{Name: name})
	return result.Error
}

func (g *NodeAgentDao) List(query ...string) ([]NodeAgentGORM, error) {
	var nodeList []NodeAgentGORM
	if len(query) == 0 {
		query = []string{"name", "id", "status", "last_update_time"}
	}
	result := g.Client.Engine.
		Select(query).Where("name is not null and name != '' ").
		Order("last_update_time desc").
		Find(&nodeList)
	if result.Error != nil {
		return nil, fmt.Errorf("query node agent list error:%v", result.Error)
	}
	return nodeList, nil
}

func (g *NodeAgentDao) ListByStatus(states int) ([]NodeAgentGORM, error) {
	var nodeList []NodeAgentGORM
	result := g.Client.Engine.Where(&NodeAgentGORM{Status: states}).Order("create_time desc").Find(&nodeList)
	if result.Error != nil {
		return make([]NodeAgentGORM, 0), fmt.Errorf("query node agent list by Status error:%v", result.Error)
	}
	return nodeList, nil
}

func (g *NodeAgentDao) QueryByName(name string) (*NodeAgentGORM, error) {
	node := &NodeAgentGORM{}
	result := g.Client.Engine.Where(&NodeAgentGORM{Name: name}).First(node)
	if result.Error != nil {
		return nil, fmt.Errorf("query node agent by name error:%v", result.Error)
	}
	return node, nil
}

func (g *NodeAgentDao) QueryByID(id string) (*NodeAgentGORM, error) {
	node := &NodeAgentGORM{}
	result := g.Client.Engine.Where(&NodeAgentGORM{ID: id}).First(node)
	if result.Error != nil {
		return nil, fmt.Errorf("query node agent by id error:%v", result.Error)
	}
	return node, nil
}
