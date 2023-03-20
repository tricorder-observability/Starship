package wasm

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	bazelutils "github.com/tricorder/src/testing/bazel"
	"github.com/tricorder/src/utils/exec"
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

	stdout, stderr, err := exec.Run("ls", "-l", path.Join(wasiSDKPath, "bin", "clang"))
	assert.NoError(err)
	assert.Equal("", stdout)
	assert.Equal("", stderr)

	testWASMCode3 := "aaaaa"
	wasiCompiler = NewWASICompiler(wasiSDKPath, wasmStarshipIncudePath, tmpBuildDir)
	_, err = wasiCompiler.BuildC(testWASMCode3)
	assert.NotNil(err)

	// contain starship common header
	testWASMCode4 := `
	#include "cJSON.h"
	#include "io.h"
	#include <assert.h>
	#include <stdint.h>
	#include <string.h>

	struct detectionPackets {
	  unsigned long long nb_ddos_packets;
	} __attribute__((packed));

	static_assert(sizeof(struct detectionPackets) == 8,
				  "Size of detectionPackets is not 8");

	// A simple function to copy entire input buf to output buffer.
	// Return 0 if succeeded.
	// Return 1 if failed to malloc output buffer.
	int write_events_to_output() {
	  struct detectionPackets *detection_packet = get_input_buf();

	  cJSON *root = cJSON_CreateObject();

	  cJSON_AddNumberToObject(root, "nb_ddos_packets",
							  detection_packet->nb_ddos_packets);

	  char *json = NULL;
	  json = cJSON_Print(root);
	  cJSON_Delete(root);

	  int json_size = strlen(json);
	  void *buf = malloc_output_buf(json_size);
	  if (buf == NULL) {
		return 1;
	  }
	  copy_to_output(json, json_size);
	  // Free allocated memory from JSON_print().
	  free(json);
	  return 0;
	}

	// Do nothing
	// TODO(yaxiong): Investigate how to remove this and build wasi module without
	// main().
	int main() { return 0; }
	`
	wasiCompiler = NewWASICompiler(wasiSDKPath, wasmStarshipIncudePath, tmpBuildDir)
	wasmELF, err = wasiCompiler.BuildC(testWASMCode4)
	assert.Nil(err)
	assert.Equal(wasmELF[:4], wasmMagic)
}
