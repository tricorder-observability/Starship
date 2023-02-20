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

package wasm

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"github.com/tricorder/src/agent/wasm/programs/cgo"
	bazelutils "github.com/tricorder/src/testing/bazel"
)

// Tests that the same C struct's object has different memory layout in host environment and WASM.
// This test works by reading a memory buffer from WASM, and cast it to Cgo type with the same definition.
func TestMemLayout(t *testing.T) {
	assert := assert.New(t)

	testFilePath := "src/agent/wasm/programs/struct_test.wasm"

	wasm, err := bazelutils.ReadTestBinFile(testFilePath)
	assert.Nil(err)

	module, err := NewWasiModule(wasm, []string{})
	assert.Nil(err)

	ret, err := module.Run("get_ints_t_size")
	assert.NotNil(ret)
	assert.Nil(err)

	v, err := unpackInt32Intf(ret)
	assert.Nil(err)
	assert.Equal(16, int(v))

	var cInts cgo.CPackedInts
	assert.Equal(16, int(unsafe.Sizeof(cInts)))
}
