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

package testing

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tricorder/src/utils/file"
)

func TestCreateTmpDir(t *testing.T) {
	assert := assert.New(t)

	dir1 := CreateTmpDir()
	dir2 := CreateTmpDir()
	assert.NotEqual(dir1, dir2)
	assert.Contains(dir1, "tricorder-")
	assert.Contains(dir2, "tricorder-")
}

func TestCreateTmpFile(t *testing.T) {
	assert := assert.New(t)

	file1 := CreateTmpFile()
	assert.True(file.Exists(file1))

	file2 := CreateTmpFile()
	assert.True(file.Exists(file2))

	assert.NotEqual(file1, file2)
}

// Tests that ReadTestFile() and ReadTestBinFile() work as expected.
func TestReadTestFile(t *testing.T) {
	assert := assert.New(t)
	content, err := ReadTestFile("src/testing/bazel/test")
	assert.Nil(err)
	assert.Equal("hello test\n", content)
	binContent, err := ReadTestBinFile("src/testing/bazel/test")
	assert.Nil(err)
	assert.Equal([]byte("hello test\n"), binContent)
}
