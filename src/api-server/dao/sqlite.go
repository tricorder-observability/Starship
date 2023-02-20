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
	"fmt"

	"github.com/tricorder/src/utils/sqlite"

	"github.com/tricorder/src/utils/log"
)

// InitSqlite prepares sqlite db file and setup the initial condition.
func InitSqlite(dbPath string) (*sqlite.ORM, error) {
	log.Infof("Opening SQLite database file at %s", dbPath)

	fullDbPath, err := sqlite.PrepareSqliteDbFile(dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite datasource error %v", err)
	}
	engine, err := sqlite.NewORM(fullDbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite datasource error %v", err)
	}
	err = engine.CreateTable(&ModuleGORM{})
	if err != nil {
		return nil, fmt.Errorf("create code table error %v", err)
	}
	err = engine.CreateTable(&GrafanaAPIKeyGORM{})
	if err != nil {
		return nil, fmt.Errorf("create grafana_api table error %v", err)
	}
	return engine, nil
}
