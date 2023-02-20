// Copyright (C) 2023  tricorder-observability
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

package sqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"

	bazelutils "github.com/tricorder/src/testing/bazel"
)

type SchemaDemo struct {
	ID   uint `gorm:"primarykey"`
	Name string
	Ebpf string
	Wasm string
}

// TestClient tests sqlite.ORM's main APIs.
func TestClient(t *testing.T) {
	assert := assert.New(t)

	dbFile := bazelutils.CreateTmpFile()

	client, err := NewORM(dbFile)
	assert.Nil(err)

	// will create table schema_demos
	err = client.CreateTable(&SchemaDemo{})
	assert.Nil(err)

	schemaDemo := SchemaDemo{Name: "schema", Ebpf: "Ebpf", Wasm: "Wasm"}

	// save data and id autecreate from 1
	// INSERT INTO `schema_demos` (`name`,`ebpf`, `wasm`) VALUES ("schema", "Ebpf", "Wasm"));
	result := client.Engine.Create(&schemaDemo)
	assert.Nil(result.Error)

	schema := SchemaDemo{ID: 1}

	// find first data by id schema.ID=1
	// select * from schema_demos where id = 1 order by id limit 1
	searchResult := client.Engine.First(&schema)
	assert.Nil(searchResult.Error)
	assert.Equal(uint(1), schema.ID)

	// select * from schema_demos where name = 'schema' order by limit 1
	searchByName := client.Engine.Where(&SchemaDemo{Name: "schema"}).First(&schema)
	assert.Nil(searchByName.Error)
	assert.Equal("schema", schema.Name)

	// update data where id = 1
	// update schema_demos set name = 'hello', ebpf = 'update' where id = 1
	updateResult := client.Engine.Model(&SchemaDemo{ID: 1}).Updates(SchemaDemo{Name: "hello", Ebpf: "update"})
	assert.Nil(updateResult.Error)

	// delete data where id = 1
	// DELETE FROM `schema_demos` WHERE 1=1
	deleteResult := client.Engine.Delete(&SchemaDemo{}, 1)
	assert.Nil(deleteResult.Error)
}
