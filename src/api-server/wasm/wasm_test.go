package wasm

import (
	"testing"

	"github.com/stretchr/testify/assert"

	testuitls "github.com/tricorder/src/testing/bazel"
)

const (
	testdataGoodWASM = "src/api-server/wasm/testdata/good.wasm"
	testdataBadWASM1 = "src/api-server/wasm/testdata/bad_fmt.wasm"
	testdataBadWASM2 = "src/api-server/wasm/testdata/bad_magic_num.wasm"
	testdataBadWASM3 = "src/api-server/wasm/testdata/bad_suffix.wa"
)

func TestIsWasmELF(t *testing.T) {
	assert := assert.New(t)
	goodWASM := testuitls.TestFilePath(testdataGoodWASM)
	badWASM1 := testuitls.TestFilePath(testdataBadWASM1)
	badWASM2 := testuitls.TestFilePath(testdataBadWASM2)
	badWASM3 := testuitls.TestFilePath(testdataBadWASM3)

	assert.True(isWasmELF(goodWASM))
	assert.False(isWasmELF(badWASM1))
	assert.False(isWasmELF(badWASM2))
	assert.False(isWasmELF(badWASM3))
}
