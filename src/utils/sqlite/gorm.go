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
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// A Object–relational mapping (https://en.wikipedia.org/wiki/Object-relational_mapping)
// type for intermediating between Golang types and Sqlite DB.
type ORM struct {
	// The underlying engine that provides all APIs.
	Engine *gorm.DB
}

// NewROM Returns a new ORM object.
func NewORM(dbfile string) (*ORM, error) {
	client := new(ORM)
	engine, err := gorm.Open(sqlite.Open(dbfile),
		// See https://gorm.io/docs/gorm_config.html for detailed configurations.
		&gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not connect to database at '%s', error: %v", dbfile, err)
	}
	client.Engine = engine
	return client, nil
}

// CreateTable wraps gorm.DB.AutoMigrate()
func (g *ORM) CreateTable(schema interface{}) error {
	return g.Engine.AutoMigrate(schema)
}
