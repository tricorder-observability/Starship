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
	"fmt"
	"os"
	"strings"

	"github.com/tricorder/src/utils/errors"
	"github.com/tricorder/src/utils/log"

	// Import sqlite driver.
	_ "github.com/mattn/go-sqlite3"
)

const (
	SqliteDBFileName = "tricorder.db"
)

// PrepareSqliteDbFile will prepare a sqlite db file according the specified dir path.
// sqlite db file absolute path will be returned.
func PrepareSqliteDbFile(dirPath string) (string, error) {
	// check is the dir is existed.
	if _, err := os.Stat(dirPath); errors.Is(err, os.ErrNotExist) {
		log.Warnf("Dir '%s' does not exist, creat it now", dirPath)
		// If it is a multi tier folder, recursively create all folders
		// otherwise an error will be reported if the folder does not exist
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return "", errors.Wrap("preparing SQLite DB file", "create parent directory", err)
		}
	}

	// check the db file is existed, if not, create it.
	var sqliteDbFilePath string
	if strings.HasSuffix(dirPath, "/") {
		sqliteDbFilePath = fmt.Sprintf("%s%s", dirPath, SqliteDBFileName)
	} else {
		sqliteDbFilePath = fmt.Sprintf("%s/%s", dirPath, SqliteDBFileName)
	}

	// db file is existed, return directly.
	return sqliteDbFilePath, nil
}
