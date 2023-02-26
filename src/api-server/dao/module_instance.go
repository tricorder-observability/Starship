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
	// Init means the module have insert module deploy agent record,
	// but not trigger the module deploy agent.
	ModuleInstanceInit = 0
	// Waiting means the module deploy command have send to module deploy agent,
	// but not receive the module deploy agent response.
	ModuleInstanceWaiting = 1
	// Success means the module deploy command have send to module deploy agent,
	// and receive the module deploy agent response success.
	ModuleInstanceSucceeed = 2
	// Failed means the module deploy command have send to module deploy agent,
	// and receive the module deploy agent response failed.
	ModuleInstanceFailed = 3
)

// ModuleInstanceGORM sqlite gorm storage and response object
type ModuleInstanceGORM struct {
	ID             string     `gorm:"'id' primarykey" json:"id,omitempty"`
	ModuleID       string     `gorm:"module_id" json:"module_id,omitempty"`
	ModuleName     string     `gorm:"module_name" json:"module_name,omitempty"`
	NodeName       string     `gorm:"node_name" json:"node_name,omitempty"`
	AgentID        string     `gorm:"agent_id" json:"agent_id,omitempty"`
	State          int        `gorm:"state" json:"state,omitempty"`
	DesireState    int        `gorm:"desire_state" json:"desire_state,omitempty"`
	CreateTime     *time.Time `gorm:"create_time" json:"create_time,omitempty"`
	LastUpdateTime *time.Time `gorm:"last_update_time" json:"last_update_time,omitempty"`
}

func (ModuleInstanceGORM) TableName() string {
	return "module_instance"
}

type ModuleInstanceDao struct {
	Client *sqlite.ORM
}

func (g *ModuleInstanceDao) SaveModuleInstance(module *ModuleInstanceGORM) error {
	module.CreateTime = &time.Time{}
	*module.CreateTime = time.Now()
	module.LastUpdateTime = &time.Time{}
	*module.LastUpdateTime = time.Now()
	result := g.Client.Engine.Create(module)
	return result.Error
}

func (g *ModuleInstanceDao) UpdateByID(module *ModuleInstanceGORM) error {
	if len(module.NodeName) == 0 {
		return fmt.Errorf("name is empty")
	}

	module.LastUpdateTime = &time.Time{}
	*module.LastUpdateTime = time.Now()
	result := g.Client.Engine.Model(&ModuleInstanceGORM{}).Where("id", module.ID).Updates(module)
	return result.Error
}

func (g *ModuleInstanceDao) UpdateStatusByID(ID string, statue int) error {
	module := ModuleInstanceGORM{}

	module.LastUpdateTime = &time.Time{}
	*module.LastUpdateTime = time.Now()
	module.State = statue

	result := g.Client.Engine.Model(&ModuleInstanceGORM{}).Where("id", ID).Updates(module)
	return result.Error
}

func (g *ModuleInstanceDao) DeleteByID(ID string) error {
	result := g.Client.Engine.Delete(&ModuleInstanceGORM{ID: ID})
	return result.Error
}

func (g *ModuleInstanceDao) List(query ...string) ([]ModuleInstanceGORM, error) {
	var moduleList []ModuleInstanceGORM
	if len(query) == 0 {
		query = []string{
			"id", "module_id", "module_name", "node_name", "agent_id", "state",
			"desire_state", "create_time", "last_update_time",
		}
	}
	result := g.Client.Engine.
		Select(query).Where("id is not null and id != '' ").
		Order("last_update_time desc").
		Find(&moduleList)
	if result.Error != nil {
		return nil, fmt.Errorf("query module instance list error:%v", result.Error)
	}
	return moduleList, nil
}

func (g *ModuleInstanceDao) ListByState(state int) ([]ModuleInstanceGORM, error) {
	var moduleList []ModuleInstanceGORM
	result := g.Client.Engine.Where(&ModuleInstanceGORM{State: state}).Order("create_time desc").Find(&moduleList)
	if result.Error != nil {
		return make([]ModuleInstanceGORM, 0), fmt.Errorf("query module instance list by Status error:%v", result.Error)
	}
	return moduleList, nil
}

func (g *ModuleInstanceDao) ListByNodeName(nodeName string) ([]ModuleInstanceGORM, error) {
	var moduleList []ModuleInstanceGORM
	result := g.Client.Engine.Where(&ModuleInstanceGORM{NodeName: nodeName}).Order("create_time desc").Find(&moduleList)
	if result.Error != nil {
		return make([]ModuleInstanceGORM, 0), fmt.Errorf("query module instance list by nodeName error:%v", result.Error)
	}
	return moduleList, nil
}

func (g *ModuleInstanceDao) ListByModuleID(moduleID string) ([]ModuleInstanceGORM, error) {
	var moduleList []ModuleInstanceGORM
	result := g.Client.Engine.Where(&ModuleInstanceGORM{ModuleID: moduleID}).Order("create_time desc").Find(&moduleList)
	if result.Error != nil {
		return make([]ModuleInstanceGORM, 0), fmt.Errorf("query module instance list by nodeName error:%v", result.Error)
	}
	return moduleList, nil
}

func (g *ModuleInstanceDao) ListByAgentID(agentID string) ([]ModuleInstanceGORM, error) {
	var moduleList []ModuleInstanceGORM
	result := g.Client.Engine.Where(&ModuleInstanceGORM{AgentID: agentID}).Order("create_time desc").Find(&moduleList)
	if result.Error != nil {
		return make([]ModuleInstanceGORM, 0), fmt.Errorf("query module instance list by AgentID error:%v", result.Error)
	}
	return moduleList, nil
}

func (g *ModuleInstanceDao) QueryByID(ID string) (*ModuleInstanceGORM, error) {
	module := &ModuleInstanceGORM{}
	result := g.Client.Engine.Where(&ModuleInstanceGORM{ID: ID}).First(module)
	if result.Error != nil {
		return nil, fmt.Errorf("query module instance by id error:%v", result.Error)
	}
	return module, nil
}
