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
	"gorm.io/gorm/clause"

	"github.com/tricorder/src/utils/sqlite"
)

// ModuleGORM sqlite gorm storage and response object
type ModuleGORM struct {
	// tag schema https://gorm.io/docs/models.html#Fields-Tags
	ID                 string `gorm:"column:id;primaryKey" json:"id,omitempty"`
	Name               string `gorm:"column:name" json:"name,omitempty"`
	DesireState        int    `gorm:"column:desire_state" json:"desire_state,omitempty"`
	CreateTime         string `gorm:"column:create_time" json:"create_time,omitempty"`
	Ebpf               string `gorm:"column:ebpf" json:"ebpf,omitempty"`
	EbpfFmt            int    `gorm:"column:ebpf_fmt" json:"ebpf_fmt,omitempty"`
	EbpfLang           int    `gorm:"column:ebpf_lang" json:"ebpf_lang,omitempty"`
	EbpfPerfBufferName string `gorm:"column:ebpf_perf_name" json:"ebpf_perf_name,omitempty"`
	EbpfProbes         string `gorm:"column:ebpf_probes" json:"ebpf_probes,omitempty"`
	// wasm store the whole wasm file content
	Wasm       []byte `gorm:"column:wasm" json:"wasm,omitempty"`
	SchemaName string `gorm:"column:schema_name" json:"schema_name,omitempty"`
	SchemaAttr string `gorm:"column:schema_attr" json:"schema_attr,omitempty"`
	Fn         string `gorm:"column:fn" json:"fn,omitempty"`
	WasmFmt    int    `gorm:"column:wasm_fmt" json:"wasm_fmt,omitempty"`
	WasmLang   int    `gorm:"column:wasm_lang" json:"wasm_lang,omitempty"`
}

func (ModuleGORM) TableName() string {
	return "module"
}

type ModuleDao struct {
	Client *sqlite.ORM
}

func (g *ModuleDao) SaveModule(mod *ModuleGORM) error {
	result := g.Client.Engine.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(mod)
	return result.Error
}

func (g *ModuleDao) UpdateByID(mod *ModuleGORM) error {
	result := g.Client.Engine.Model(&ModuleGORM{ID: mod.ID}).Updates(mod)
	return result.Error
}

func (g *ModuleDao) UpdateStatusByID(id string, status int) error {
	result := g.Client.Engine.Model(&ModuleGORM{ID: id}).Select("desire_state").Updates(ModuleGORM{DesireState: status})
	return result.Error
}

func (g *ModuleDao) DeleteByID(id string) error {
	result := g.Client.Engine.Delete(&ModuleGORM{ID: id})
	return result.Error
}

func (g *ModuleDao) ListModule(fields []string) ([]ModuleGORM, error) {
	moduleList := make([]ModuleGORM, 0)
	result := g.Client.Engine.
		Select(fields).Where("name is not null and name != '' ").
		Order("create_time desc").
		Find(&moduleList)
	if result.Error != nil {
		return nil, result.Error
	}
	return moduleList, nil
}

func (g *ModuleDao) ListModuleByStatus(status int) ([]ModuleGORM, error) {
	moduleList := make([]ModuleGORM, 0)
	result := g.Client.Engine.Where(&ModuleGORM{DesireState: status}).Order("create_time desc").Find(&moduleList)
	if result.Error != nil {
		return nil, result.Error
	}
	return moduleList, nil
}

// queryAtMostOneRecord returns at most one record, basically simulate gorm First() without producing error logging.
func (g *ModuleDao) queryAtMostOneRecord(m *ModuleGORM) (*ModuleGORM, error) {
	module := &ModuleGORM{}
	result := g.Client.Engine.Where(m).Find(module)
	if result.RowsAffected == 0 {
		return nil, result.Error
	}
	return module, result.Error
}

func (g *ModuleDao) QueryByName(name string) (*ModuleGORM, error) {
	return g.queryAtMostOneRecord(&ModuleGORM{Name: name})
}

func (g *ModuleDao) QueryByID(id string) (*ModuleGORM, error) {
	return g.queryAtMostOneRecord(&ModuleGORM{ID: id})
}
