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

package tar

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	testuitls "github.com/tricorder/src/testing/bazel"

	"github.com/tricorder/src/utils/file"
)

func TestGZExtract(t *testing.T) {
	assert := assert.New(t)

	tmpDir := testuitls.CreateTmpDir()
	assert.Nil(GZExtract("testdata/test.tar.gz", tmpDir))
	helloPath := path.Join(tmpDir, "hello.txt")
	assert.Equal(file.Exists(helloPath), true)

	tmpDir = testuitls.CreateTmpDir()
	assert.NotNil(GZExtract("testdata/wrong_file_format.tar.gz", tmpDir))
}
