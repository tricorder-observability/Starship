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
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tricorder/src/testing/bazel"
	"github.com/tricorder/src/utils/file"
)

func TestInitDBFile(t *testing.T) {
	assert := assert.New(t)

	dir := path.Join(bazel.CreateTmpDir(), "test", "sqlite")

	testCases := []struct {
		caseStr        string
		dirPath        string
		wantDBFilePath string
		err            error
	}{
		{
			caseStr:        "successful create db file with dir suffix",
			dirPath:        fmt.Sprintf("%s//", dir),
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
		dbFilePath, err := PrepareSQLiteDBDir(tc.dirPath)
		assert.NoError(err)
		assert.True(file.Exists(tc.dirPath))
		assert.Equal(tc.wantDBFilePath, dbFilePath)
	}
}
