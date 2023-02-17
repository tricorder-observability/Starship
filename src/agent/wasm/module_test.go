package wasm

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	bazelutils "github.com/tricorder/src/testing/bazel"
	sysutils "github.com/tricorder/src/testing/sys"
	"github.com/tricorder/src/utils/file"
)

// getEvents return a http event list for test
func getEvents() []string {
	return []string{
		"GET /api/v1/bpf HTTP/1.1\r\nHost: tricorder.dev\r\nCookie: cookie\r\n\r\n",
		"PUT /api/v2/wasm HTTP/1.1\r\nHost: tricorder.dev\r\n\r\n",
		"GET /wp-content/uploads/2010/03/hello-kitty-darth-vader-pink.jpg HTTP/1.1\r\n" +
			"Host: www.kittyhell.com\r\n" +
			"User-Agent: Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10.6; ja-JP-mac; rv:1.9.2.3) " +
			"Gecko/20100401 Firefox/3.6.3 " +
			"Pathtraq/0.9\r\n" +
			"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\n" +
			"Accept-Language: ja,en-us;q=0.7,en;q=0.3\r\n" +
			"Accept-Encoding: gzip,deflate\r\n" +
			"Accept-Charset: Shift_JIS,utf-8;q=0.7,*;q=0.7\r\n" +
			"Keep-Alive: 115\r\n" +
			"Connection: keep-alive\r\n" +
			"Cookie: wp_ozh_wsa_visits=2; wp_ozh_wsa_visit_lasttime=xxxxxxxxxx; " +
			"__utma=xxxxxxxxx.xxxxxxxxxx.xxxxxxxxxx.xxxxxxxxxx.xxxxxxxxxx.x; " +
			"__utmz=xxxxxxxxx.xxxxxxxxxx.x.x.utmccn=(referral)|utmcsr=reader.livedoor.com|utmcct=" +
			"/reader/|utmcmd=referral" +
			"\r\n" +
			"\r\n",
	}
}

func newWasmModule(watRelPath string, fn func()) (*Module, error) {
	watFilePath := bazelutils.TestFilePath(watRelPath)
	wasm, err := wat2Wasm(watFilePath)
	if err != nil {
		return nil, err
	}
	module, err := NewWasmModule(wasm, fn)
	if err != nil {
		return nil, err
	}
	return module, nil
}

func TestWasmFnPicoWasm(t *testing.T) {
	testFilePath := "src/agent/wasm/programs/pico/pico.wasm"

	wasmFilePath := bazelutils.TestFilePath(testFilePath)
	wasm, err := file.ReadBin(wasmFilePath)
	if err != nil {
		t.Errorf("Failed to read wasm file %s, error: %v", wasmFilePath, err)
	}

	module, err := newBasicModule(wasm)
	if err != nil {
		t.Errorf("failed to create new module, error: %v", err)
	}

	argv := getEvents()
	args := append([]string{strconv.Itoa(len(argv))}, argv...)
	err = module.NewWasiInstance(args)
	if err != nil {
		t.Errorf("failed to create new instance with linker, error: %v", err)
	}

	out := sysutils.CaptureStdout(func() {
		_, err = module.Run("_start")
		if err != nil {
			t.Errorf("Failed to run '_start', error: %v", err)
		}
	})

	if out != "" {
		t.Errorf("Expect to get, got '%s'", out)
	}

	// Repeat creating instance and invoking functions.
	// This reuses the result of creating module.
	err = module.NewWasiInstance(args)
	if err != nil {
		t.Errorf("failed to create new instance with linker, error: %v", err)
	}

	out = sysutils.CaptureStdout(func() {
		_, err = module.Run("_start")
		if err != nil {
			t.Errorf("Failed to run '_start', error: %v", err)
		}
	})

	if out != "" {
		t.Errorf("Expect to get, got '%s'", out)
	}
}

func TestWasmFnMissArgs(t *testing.T) {
	testFilePath := "src/agent/wasm/programs/pico/pico.wasm"

	wasmFilePath := bazelutils.TestFilePath(testFilePath)
	wasm, err := file.ReadBin(wasmFilePath)
	if err != nil {
		t.Errorf("Failed to read wasm file %s, error: %v", wasmFilePath, err)
	}

	module, err := NewWasiModule(wasm, getEvents())
	if err != nil {
		t.Errorf("Failed to create new WASM module, error: %v", err)
	}

	_, err = module.Run("phr_parse_request")
	wantErrMsg := "expected 10 arguments, got 0"
	assert.Containsf(t, err.Error(), wantErrMsg, "expected error containing %q, got %s", wantErrMsg, err)
}

const watRelPath = "src/agent/wasm/programs/hello_world.wat"

func TestWasmFn(t *testing.T) {
	assert := assert.New(t)

	module, err := newWasmModule(watRelPath, func() {
		fmt.Print("Hello from Go!")
	})
	assert.Nil(err)

	out := sysutils.CaptureStdout(func() {
		_, err = module.Run("run")
		if err != nil {
			t.Errorf("Failed to run 'run', error: %v", err)
		}
	})

	if out != "Hello from Go!" {
		t.Errorf("Expect to get, got '%s'", out)
	}
}

// Tests that the debug string lines are as expected.
func TestDebugString(t *testing.T) {
	assert := assert.New(t)

	module, err := newWasmModule(watRelPath, func() {
		fmt.Print("Hello from Go!")
	})
	assert.Nil(err)

	lines := module.DebugString()
	assert.Equal([]string{"Imports: hello, [], []", "Exports: run, [], []"}, lines)
}
