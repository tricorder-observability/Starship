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

// NodeAgentGORM sqlite gorm storage and response object
type NodeAgentGORM struct {
	NodeName       string     `gorm:"'node_name' primarykey" json:"node_name,omitempty"`
	AgentID        string     `gorm:"agent_id" json:"agent_id,omitempty"`
	State          int        `gorm:"state" json:"state,omitempty"`
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
	agent.CreateTime = &time.Time{}
	*agent.CreateTime = time.Now()
	agent.LastUpdateTime = &time.Time{}
	*agent.LastUpdateTime = time.Now()
	result := g.Client.Engine.Create(agent)
	return result.Error
}

func (g *NodeAgentDao) UpdateByName(agent *NodeAgentGORM) error {
	if len(agent.NodeName) == 0 {
		return fmt.Errorf("name shuold not be empty")
	}

	agent.LastUpdateTime = &time.Time{}
	*agent.LastUpdateTime = time.Now()
	result := g.Client.Engine.Model(&NodeAgentGORM{}).Where("node_name", agent.NodeName).Updates(agent)
	return result.Error
}

func (g *NodeAgentDao) UpdateStateByName(nodeName string, statue int) error {
	agent := NodeAgentGORM{}

	agent.LastUpdateTime = &time.Time{}
	*agent.LastUpdateTime = time.Now()
	agent.State = statue

	result := g.Client.Engine.Model(&NodeAgentGORM{}).Where("node_name", nodeName).Updates(agent)
	return result.Error
}

func (g *NodeAgentDao) DeleteByName(nodeName string) error {
	result := g.Client.Engine.Delete(&NodeAgentGORM{NodeName: nodeName})
	return result.Error
}

func (g *NodeAgentDao) List(query ...string) ([]NodeAgentGORM, error) {
	nodeList := make([]NodeAgentGORM, 0)
	if len(query) == 0 {
		query = []string{"node_name", "agent_id", "state", "create_time", "last_update_time"}
	}

	result := g.Client.Engine.
		Select(query).Where("node_name is not null and node_name != '' ").
		Order("last_update_time desc").
		Find(&nodeList)
	if result.Error != nil {
		return nil, fmt.Errorf("query node agent list error:%v", result.Error)
	}
	return nodeList, nil
}

func (g *NodeAgentDao) ListByState(state int) ([]NodeAgentGORM, error) {
	nodeList := make([]NodeAgentGORM, 0)
	result := g.Client.Engine.Where(&NodeAgentGORM{State: state}).Order("create_time desc").Find(&nodeList)
	if result.Error != nil {
		return make([]NodeAgentGORM, 0), fmt.Errorf("query node agent list by Status error:%v", result.Error)
	}
	return nodeList, nil
}

func (g *NodeAgentDao) QueryByName(nodeName string) (*NodeAgentGORM, error) {
	node := &NodeAgentGORM{}
	result := g.Client.Engine.Where(&NodeAgentGORM{NodeName: nodeName}).First(node)
	if result.Error != nil {
		return nil, fmt.Errorf("query node agent by name error:%v", result.Error)
	}
	return node, nil
}

func (g *NodeAgentDao) QueryByID(agentID string) (*NodeAgentGORM, error) {
	node := &NodeAgentGORM{}
	result := g.Client.Engine.Where(&NodeAgentGORM{AgentID: agentID}).First(node)
	if result.Error != nil {
		return nil, fmt.Errorf("query node agent by id error:%v", result.Error)
	}
	return node, nil
}
