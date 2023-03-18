package wasm

import (
	"testing"

	"github.com/stretchr/testify/assert"

	bazelutils "github.com/tricorder/src/testing/bazel"
)

func TestWASMBUILDC(t *testing.T) {
	assert := assert.New(t)

	wasiBazelFilePath := "src/api-server/cmd/"
	wasiSDKPath := bazelutils.TestFilePath(wasiBazelFilePath) + "/wasi-sdk-19.0"
	wasiBazelIncludeFilePath := "modules/common"
	wasmStarshipIncudePath := bazelutils.TestFilePath(wasiBazelIncludeFilePath)
	tmpBuildDir := "/tmp"
	// Compare with wasm magic number: \x00\x61\x73\x6d (0x6d736100 in little endian)
	wasmMagic := []byte{0x00, 0x61, 0x73, 0x6d}

	// c hello function with main function
	testWASMCode1 := `
int hello() {
	return 0;
}
int main() { return 0; }`

	wasiCompiler := NewWASICompiler(wasiSDKPath, wasmStarshipIncudePath, tmpBuildDir)
	wasmELF, err := wasiCompiler.BuildC(testWASMCode1)
	assert.Nil(err)
	assert.Equal(wasmELF[:4], wasmMagic)

	// c hello function without main function
	testWASMCode2 := `
	int hello() {
		return 0;
	}		
	`
	wasiCompiler = NewWASICompiler(wasiSDKPath, wasmStarshipIncudePath, tmpBuildDir)
	_, err = wasiCompiler.BuildC(testWASMCode2)
	assert.NotNil(err)

	testWASMCode3 := "aaaaa"
	wasiCompiler = NewWASICompiler(wasiSDKPath, wasmStarshipIncudePath, tmpBuildDir)
	_, err = wasiCompiler.BuildC(testWASMCode3)
	assert.NotNil(err)

	// contain starship common header
	testWASMCode4 := `
	#include <string.h>

	#include "io.h"
	
	// A simple function to copy entire input buf to output buffer.
	void copy_input_to_output() {
	  malloc_output_buf(input_buf.capacity);
	  if (can_write_to_output_buf(input_buf.length)) {
		copy_to_output(input_buf.data, input_buf.length);
	  }
	}
	
	int main() { return 0; }
	`
	wasiCompiler = NewWASICompiler(wasiSDKPath, wasmStarshipIncudePath, tmpBuildDir)
	wasmELF, err = wasiCompiler.BuildC(testWASMCode4)
	assert.Nil(err)
	assert.Equal(wasmELF[:4], wasmMagic)
}
