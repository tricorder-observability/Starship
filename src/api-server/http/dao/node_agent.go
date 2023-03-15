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

	"gorm.io/gorm/clause"

	"github.com/tricorder/src/utils/sqlite"
)

// NodeAgentGORM sqlite gorm storage and response object
type NodeAgentGORM struct {
	// tag schema https://gorm.io/docs/models.html#Fields-Tags
	AgentID        string     `gorm:"column:agent_id;primaryKey" json:"agent_id,omitempty"`
	NodeName       string     `gorm:"column:node_name" json:"node_name,omitempty"`
	AgentPodID     string     `gorm:"column:agent_pod_id" json:"agent_pod_id,omitempty"`
	State          int        `gorm:"column:state" json:"state,omitempty"`
	CreateTime     *time.Time `gorm:"column:create_time" json:"create_time,omitempty"`
	LastUpdateTime *time.Time `gorm:"column:last_update_time" json:"last_update_time,omitempty"`
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
	result := g.Client.Engine.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(agent)
	return result.Error
}

func (g *NodeAgentDao) UpdateByID(agent *NodeAgentGORM) error {
	if len(agent.NodeName) == 0 {
		return fmt.Errorf("name shuold not be empty")
	}

	agent.LastUpdateTime = &time.Time{}
	*agent.LastUpdateTime = time.Now()
	result := g.Client.Engine.Model(&NodeAgentGORM{}).Where("agent_id", agent.AgentID).Updates(agent)
	return result.Error
}

func (g *NodeAgentDao) UpdateStateByID(agentID string, statue int) error {
	agent := NodeAgentGORM{}

	agent.LastUpdateTime = &time.Time{}
	*agent.LastUpdateTime = time.Now()
	agent.State = statue

	result := g.Client.Engine.Model(&NodeAgentGORM{}).Where("agent_id", agentID).
		Select("state", "last_update_time").Updates(agent)
	return result.Error
}

func (g *NodeAgentDao) DeleteByID(agentID string) error {
	result := g.Client.Engine.Delete(&NodeAgentGORM{AgentID: agentID})
	return result.Error
}

func (g *NodeAgentDao) List(query []string) ([]NodeAgentGORM, error) {
	nodeList := make([]NodeAgentGORM, 0)
	if len(query) == 0 {
		query = []string{"agent_id", "node_name", "agent_pod_id", "state", "create_time", "last_update_time"}
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
		return nil, fmt.Errorf("query node agent list by Status error:%v", result.Error)
	}
	return nodeList, nil
}

func (g *NodeAgentDao) ListByNodeName(nodeName string) ([]NodeAgentGORM, error) {
	nodeList := make([]NodeAgentGORM, 0)
	result := g.Client.Engine.Where(&NodeAgentGORM{NodeName: nodeName}).Order("create_time desc").Find(&nodeList)
	if result.Error != nil {
		return nil, fmt.Errorf("query node agent list by Name error:%v", result.Error)
	}
	return nodeList, nil
}

func (g *NodeAgentDao) QueryByID(agentID string) (*NodeAgentGORM, error) {
	node := &NodeAgentGORM{}
	result := g.Client.Engine.Where(&NodeAgentGORM{AgentID: agentID}).First(node)
	if result.Error != nil {
		return nil, fmt.Errorf("query node agent by id error:%v", result.Error)
	}
	return node, nil
}

func (g *NodeAgentDao) QueryByPodID(agentPodID string) (*NodeAgentGORM, error) {
	node := &NodeAgentGORM{}
	result := g.Client.Engine.Where(&NodeAgentGORM{AgentPodID: agentPodID}).First(node)
	if result.Error != nil {
		return nil, fmt.Errorf("query node agent by pod id error:%v", result.Error)
	}
	return node, nil
}
