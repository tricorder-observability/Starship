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

	testutils "github.com/tricorder/src/testing/bazel"
)

// Tests that Load() can load a .tar file.
func TestLoad(t *testing.T) {
	assert := assert.New(t)

	var cli CLI
	tarFilePath := testutils.TestFilePath("src/testing/docker/testdata/test_image.tar")
	assert.Equal("bazel/src/testing/docker/testdata:test_image", cli.Load(tarFilePath))
}
