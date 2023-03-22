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

package sqlite

import (
	"os"
	"path"

	"github.com/tricorder/src/utils/errors"
	"github.com/tricorder/src/utils/file"
	"github.com/tricorder/src/utils/log"

	// Import sqlite driver.
	_ "github.com/mattn/go-sqlite3"
)

const (
	SqliteDBFileName = "tricorder.db"
)

// PrepareSQLiteDBDir will prepare a sqlite db file according the specified dir path.
// sqlite db file absolute path will be returned.
func PrepareSQLiteDBDir(dirPath string) (string, error) {
	if file.Exists(dirPath) {
		return path.Join(dirPath, SqliteDBFileName), nil
	}
	log.Warnf("Dir '%s' does not exist, creat it now", dirPath)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return "", errors.Wrap("preparing SQLite DB file", "create parent directory", err)
	}
	return path.Join(dirPath, SqliteDBFileName), nil
}
