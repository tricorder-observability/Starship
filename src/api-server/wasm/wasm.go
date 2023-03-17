package wasm

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/tricorder/src/utils/errors"
	"github.com/tricorder/src/utils/uuid"
)

const (
	WASIClang           = "WASI_SDK_PATH=/opt/wasi-sdk /opt/wasi-sdk/bin/clang"
	WASICFlags          = "--sysroot=/opt/wasi-sdk/share/wasi-sysroot -Wall -Wextra -Wl,--export-all"
	WASIStarshipInclude = "/opt/starship/include"
	BuildTmpDir         = "/tmp"
)

func isWasmELF(filePath string) bool {
	if !strings.HasSuffix(filePath, ".wasm") {
		return false
	}

	file, err := os.Open(filePath)
	if err != nil {
		return false
	}

	// Read the first four bytes of the file
	buf := make([]byte, 4)
	n, err := file.Read(buf)
	if err != nil || n < 4 {
		return false
	}

	// Compare with wasm magic number: \x00\x61\x73\x6d (0x6d736100 in little endian)
	wasmMagic := []byte{0x00, 0x61, 0x73, 0x6d}
	isWasm := true

	for i := range buf {
		if buf[i] != wasmMagic[i] {
			isWasm = false
			break
		}
	}

	return isWasm
}

func WASICompileC(code string) ([]byte, error) {
	srcID := strings.Replace(uuid.New(), "-", "_", -1)
	srcFilePath := BuildTmpDir + "/" + srcID + ".c"
	dstFilePath := BuildTmpDir + "/" + srcID + ".wasm"

	// write code string to tmp file
	phase := "write code to " + srcFilePath
	_, err := os.Stat(srcFilePath)
	if errors.Is(err, os.ErrNotExist) {
		content := []byte(code)
		err = ioutil.WriteFile(srcFilePath, content, 0o644)
		if err != nil {
			return nil, errors.Wrap("compile wasm code", phase, err)
		}
	} else if err == nil {
		return nil, errors.New("compile wasm code", phase+" error: File already exists.")
	} else {
		return nil, errors.Wrap("compile wasm code", phase, err)
	}

	// compile code
	phase = "compile " + srcFilePath + " to " + dstFilePath
	cmd := exec.Command(WASIClang, WASICFlags, "-L"+WASIStarshipInclude, srcFilePath, "-o", dstFilePath)
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap("compile wasm code", phase, err)
	}

	if len(out) > 0 {
		return nil, errors.New("compile wasm code", phase+" error: cc output:\n"+string(out))
	}

	// check compiled file if exists
	phase = "check compiled file " + dstFilePath
	_, err = os.Stat(dstFilePath)
	if err != nil {
		return nil, errors.Wrap("compile wasm code", phase, err)
	}

	// check comiled file fmt
	phase = "check compiled file format " + dstFilePath
	if !isWasmELF(dstFilePath) {
		return nil, errors.New("compile wasm code", phase+" error: File is not a wasm file.")
	}

	// read compiled file
	phase = "read compiled file " + dstFilePath
	data, err := ioutil.ReadFile("input.txt")
	if err != nil {
		return nil, errors.Wrap("compile wasm code", phase, err)
	}

	// delete tmp files
	phase = "delete tmp files"
	err = os.Remove(srcFilePath)
	if err != nil {
		return nil, errors.Wrap("compile wasm code", phase+" "+srcFilePath, err)
	}
	err = os.Remove(dstFilePath)
	if err != nil {
		return nil, errors.Wrap("compile wasm code", phase+" "+dstFilePath, err)
	}
	return data, nil
}
