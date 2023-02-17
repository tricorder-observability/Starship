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
