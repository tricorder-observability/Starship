package wasm

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/tricorder/src/utils/errors"
	"github.com/tricorder/src/utils/file"
	"github.com/tricorder/src/utils/uuid"
)

const (
	DefaultWASISDKPath         = "/opt/tricorder/wasm/wasi-sdk"
	DefaultWASIClang           = DefaultWASISDKPath + "/bin/clang"
	DefaultWASICFlags          = "--sysroot=" + DefaultWASISDKPath + "/share/wasi-sysroot"
	DefaultWASIStarshipInclude = "/opt/tricorder/wasm/include"
	DefaultBuildTmpDir         = "/tmp"
)

type WASICompiler struct {
	WASIClang           string
	WASICFlags          string
	WASIStarshipInclude string
	BuildTmpDir         string
}

func NewWASICompiler(WASISDKPath string,
	WASIStarshipInclude string, BuildTmpDir string,
) *WASICompiler {
	if WASISDKPath == "" {
		WASISDKPath = DefaultWASISDKPath
	}

	WASIClang := WASISDKPath + "/bin/clang"
	WASICFlags := "--sysroot=" + WASISDKPath + "/share/wasi-sysroot"

	if WASIStarshipInclude == "" {
		WASIStarshipInclude = DefaultWASIStarshipInclude
	}

	if BuildTmpDir == "" {
		BuildTmpDir = DefaultBuildTmpDir
	}

	return &WASICompiler{
		WASIClang:           WASIClang,
		WASICFlags:          WASICFlags,
		WASIStarshipInclude: WASIStarshipInclude,
		BuildTmpDir:         BuildTmpDir,
	}
}

func (w *WASICompiler) BuildC(code string) ([]byte, error) {
	srcID := strings.Replace(uuid.New(), "-", "_", -1)
	srcFilePath := w.BuildTmpDir + "/" + srcID + ".c"
	dstFilePath := w.BuildTmpDir + "/" + srcID + ".wasm"

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
	cmd := exec.Command(w.WASIClang, w.WASICFlags, w.WASIStarshipInclude+"/cJSON.c", "-I"+w.WASIStarshipInclude, srcFilePath,
		"-Wl,--export-all", "-Wall", "-Wextra", "-o", dstFilePath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap("compile wasm code", phase+" error cc output:\n"+stderr.String(), err)
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
	if !file.IsWasmELF(dstFilePath) {
		return nil, errors.New("compile wasm code", phase+" error: File is not a wasm file.")
	}

	// read compiled file
	phase = "read compiled file " + dstFilePath
	data, err := ioutil.ReadFile(dstFilePath)
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
