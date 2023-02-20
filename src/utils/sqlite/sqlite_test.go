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
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDBFile(t *testing.T) {
	dir, _ := os.Getwd()
	testCases := []struct {
		caseStr        string
		dirPath        string
		wantDBFilePath string
		err            error
	}{
		{
			caseStr:        "successful create db file with dir suffix",
			dirPath:        fmt.Sprintf("%s/", dir),
			wantDBFilePath: fmt.Sprintf("%s/%s", dir, SqliteDBFileName),
			err:            nil,
		},
		{
			caseStr:        "successful create db file without suffix",
			dirPath:        dir,
			wantDBFilePath: fmt.Sprintf("%s/%s", dir, SqliteDBFileName),
			err:            nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.caseStr, func(t *testing.T) {
			dbFilePath, err := PrepareSqliteDbFile(tc.dirPath)
			if err != nil {
				assert.Equal(t, true, strings.Contains(err.Error(), tc.err.Error()))
			}
			assert.Equal(t, tc.wantDBFilePath, dbFilePath)
			// clean up created file.
			_ = os.Remove(dbFilePath)
		})
	}
}
