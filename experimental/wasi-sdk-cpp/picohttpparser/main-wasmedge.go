package main

import (
	"fmt"
	"os"

	"github.com/second-state/WasmEdge-go/wasmedge"
)

const (
	wasmFile = "./pico.wasm"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <http header string>\n", os.Args[0])
		return
	}

	/// Set not to print debug info
	wasmedge.SetLogErrorLevel()

	/// Create configure
	var conf = wasmedge.NewConfigure(wasmedge.WASI)

	/// Create VM with configure
	var vm = wasmedge.NewVMWithConfig(conf)

	/// Init WASI
	argc := "2"                  // hardcode for c
	argv := os.Args[1]           // "GET /eee HTTP/1.1\r\nHost: www.abc.com\r\n\r\n"
	args := []string{argc, argv} // final args passed to the main() in wasm program
	var wasi = vm.GetImportModule(wasmedge.WASI)
	wasi.InitWasi(
		args,            /// The args
		os.Environ(),    /// The envs
		[]string{".:."}, /// The mapping preopens
	)

	/// Run WASM file
	fmt.Printf("Go: spawn wasm sandbox to run http parser ...\n\n")
	vm.RunWasmFile(wasmFile, "_start")
	fmt.Printf("\nGo: returned from wasm sandbox\n")

	exitcode := wasi.WasiGetExitCode()
	if exitcode != 0 {
		fmt.Println("Go: Run WASM failed, exit code:", exitcode)
	}

	vm.Release()
	conf.Release()
}
