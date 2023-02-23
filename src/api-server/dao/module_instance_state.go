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

package dao

import (
	"fmt"
	"time"

	"github.com/tricorder/src/utils/sqlite"
)

const (
	// ModuleStatusSuccess success
	ModuleStatusSuccess = 0
	// ModuleStatusFailed failed
	ModuleStatusFailed = 1
	// ModuleStatusInProgress in progress
	ModuleStatusInProgress = 2
)

// NodeAgentGORM sqlite gorm storage and response object
type ModuleInstanceORM struct {
	ID             string     `gorm:"id" json:"id,omitempty"`
	NodeName       string     `gorm:"node_name" json:"node_name,omitempty"`
	AgentID        string     `gorm:"agent_id" json:"agent_id,omitempty"`
	DesiredStatus  int        `gorm:"desired_status" json:"desired_status,omitempty"`
	Status         int        `gorm:"status" json:"status,omitempty"`
	CreateTime     *time.Time `gorm:"create_time" json:"create_time,omitempty"`
	LastUpdateTime *time.Time `gorm:"last_update_time" json:"last_update_time,omitempty"`
}

func (ModuleInstanceORM) TableName() string {
	return "module_instance_state"
}

type ModuleInstanceDao struct {
	Client *sqlite.ORM
}

func (g *ModuleInstanceDao) SaveModuleInstance(moduleInstance *ModuleInstanceORM) error {
	moduleInstance.LastUpdateTime = &time.Time{}
	*moduleInstance.LastUpdateTime = time.Now()
	result := g.Client.Engine.Create(moduleInstance)
	return result.Error
}

func (g *ModuleInstanceDao) UpdateByID(moduleInstance *ModuleInstanceORM) error {
	if len(moduleInstance.ID) == 0 {
		return fmt.Errorf("id is empty")
	}

	moduleInstance.LastUpdateTime = &time.Time{}
	*moduleInstance.LastUpdateTime = time.Now()
	moduleInstance.CreateTime = &time.Time{}
	*moduleInstance.CreateTime = time.Now()

	result := g.Client.Engine.Model(&ModuleInstanceORM{ID: moduleInstance.ID}).Updates(moduleInstance)
	return result.Error
}

func (g *ModuleInstanceDao) UpdateStatusByID(id string, status int) error {
	now := time.Now()
	result := g.Client.Engine.Model(&ModuleInstanceORM{ID: id}).Select("status").Updates(ModuleInstanceORM{Status: status, LastUpdateTime: &now})
	return result.Error
}

func (g *ModuleInstanceDao) UpdateDesiredStatusByID(id string, desiredStatus int) error {
	now := time.Now()
	result := g.Client.Engine.Model(&ModuleInstanceORM{ID: id}).Select("desired_status").Updates(ModuleInstanceORM{DesiredStatus: desiredStatus, LastUpdateTime: &now})
	return result.Error
}

func (g *ModuleInstanceDao) GetModuleInstanceByID(id string) (*ModuleInstanceORM, error) {
	moduleInstance := &ModuleInstanceORM{}
	result := g.Client.Engine.Where("id = ?", id).First(moduleInstance)
	return moduleInstance, result.Error
}

func (g *ModuleInstanceDao) Delete(id string) error {
	result := g.Client.Engine.Delete(&ModuleInstanceORM{ID: id})
	return result.Error
}

func (g *ModuleInstanceDao) DeleteByNodeName(nodeName string) error {
	result := g.Client.Engine.Where("node_name = ?", nodeName).Delete(&ModuleInstanceORM{})
	return result.Error
}

func (g *ModuleInstanceDao) DeleteByAgentID(agentID string) error {
	result := g.Client.Engine.Where("agent_id = ?", agentID).Delete(&ModuleInstanceORM{})
	return result.Error
}

func (g *ModuleInstanceDao) List(query ...string) ([]ModuleInstanceORM, error) {
	var moduleList []ModuleInstanceORM
	if len(query) == 0 {
		query = []string{"id", "node_name", "agent_id", "name", "version", "status", "desired_status", "last_update_time", "create_time"}
	}
	result := g.Client.Engine.
		Select(query).Where("id is not null and id != '' ").
		Order("last_update_time desc").
		Find(&moduleList)
	if result.Error != nil {
		return nil, fmt.Errorf("query module agent list error:%v", result.Error)
	}
	return moduleList, nil
}

func (g *ModuleInstanceDao) ListByNodeName(nodeName string, query ...string) ([]ModuleInstanceORM, error) {
	var moduleList []ModuleInstanceORM
	if len(query) == 0 {
		query = []string{"id", "node_name", "agent_id", "name", "version", "status", "desired_status", "last_update_time", "create_time"}
	}
	result := g.Client.Engine.
		Select(query).Where("node_name = ?", nodeName).
		Order("last_update_time desc").
		Find(&moduleList)
	if result.Error != nil {
		return nil, fmt.Errorf("query module agent list error:%v", result.Error)
	}
	return moduleList, nil
}

func (g *ModuleInstanceDao) ListByAgentID(agentID string, query ...string) ([]ModuleInstanceORM, error) {
	var moduleList []ModuleInstanceORM
	if len(query) == 0 {
		query = []string{"id", "node_name", "agent_id", "name", "version", "status", "desired_status", "last_update_time", "create_time"}
	}
	result := g.Client.Engine.
		Select(query).Where("agent_id = ?", agentID).
		Order("last_update_time desc").
		Find(&moduleList)
	if result.Error != nil {
		return nil, fmt.Errorf("query module agent list error:%v", result.Error)
	}
	return moduleList, nil
}

func (g *ModuleInstanceDao) ListByStatus(status int, query ...string) ([]ModuleInstanceORM, error) { //nolint:dupl
	var moduleList []ModuleInstanceORM
	if len(query) == 0 {
		query = []string{"id", "node_name", "agent_id", "name", "version", "status", "desired_status", "last_update_time", "create_time"}
	}
	result := g.Client.Engine.
		Select(query).Where("status = ?", status).
		Order("last_update_time desc").
		Find(&moduleList)
	if result.Error != nil {
		return nil, fmt.Errorf("query module agent list error:%v", result.Error)
	}
	return moduleList, nil
}

func (g *ModuleInstanceDao) QueryByID(id string, query ...string) (*ModuleInstanceORM, error) {
	var module ModuleInstanceORM
	if len(query) == 0 {
		query = []string{"id", "node_name", "agent_id", "name", "version", "status", "desired_status", "last_update_time", "create_time"}
	}
	result := g.Client.Engine.Select(query).Where("id = ?", id).First(&module)
	if result.Error != nil {
		return nil, fmt.Errorf("query module agent list error:%v", result.Error)
	}
	return &module, nil
}

func (g *ModuleInstanceDao) QueryByNodeName(nodeName string, query ...string) (*ModuleInstanceORM, error) {
	var module ModuleInstanceORM
	if len(query) == 0 {
		query = []string{"id", "node_name", "agent_id", "name", "version", "status", "desired_status", "last_update_time", "create_time"}
	}
	result := g.Client.Engine.Select(query).Where("node_name = ?", nodeName).First(&module)
	if result.Error != nil {
		return nil, fmt.Errorf("query module agent list error:%v", result.Error)
	}
	return &module, nil
}
