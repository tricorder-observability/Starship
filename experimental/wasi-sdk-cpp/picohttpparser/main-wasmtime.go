package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bytecodealliance/wasmtime-go/v3"
)

const (
	wasmFile = "./pico.wasm"

	debug = false

	// Max field lengths in http header, should be consistent with the c/wasm side definitions.
	// See main.c for the latters.
	maxHttpMethodLen = 8
	maxHttpPathLen   = 120
)

func main() {
	engine := wasmtime.NewEngine()
	module, err := wasmtime.NewModuleFromFile(engine, wasmFile)
	check(err)

	if debug {
		printImportsExports(module)
	}

	// Create a linker with WASI functions defined within it
	fmt.Printf("Go: creating new linker\n")
	linker := wasmtime.NewLinker(engine)
	err = linker.DefineWasi()
	if err != nil {
		log.Fatal(err)
	}

	// Create store
	wasiConfig := wasmtime.NewWasiConfig()
	wasiConfig.InheritEnv()
	wasiConfig.InheritStdout()
	wasiConfig.InheritStdin()
	wasiConfig.InheritStderr()
	if len(os.Args) == 1 {
		fmt.Printf("Go: http request not provided through CLI, using built-in ones")
		wasiConfig.SetArgv(getArgs())
	} else {
		wasiConfig.InheritArgv()
	}

	store := wasmtime.NewStore(engine)
	store.SetWasi(wasiConfig)
	instance, err := linker.Instantiate(store, module)
	if err != nil {
		log.Fatal(err)
	}

	// Run the main() function in sandbox (with input data provided via argument list)
	fmt.Printf("Go: calling into sandbox main() ...\n")
	mainFn := instance.GetFunc(store, "_start")
	if mainFn == nil {
		panic("get function failed")
	}

	_, err = mainFn.Call(store)
	if err != nil {
		log.Fatalf("Call wasm function failed: %s\n", err)
	}

	// Get sandbox results via memory
	if debug {
		inspectSandboxMemory(instance, store)
	}

	fmt.Printf("Go: calling into sandbox get_result_buf() ...\n")
	getResultBufFn := instance.GetFunc(store, "get_result_buf")
	ptr, err := getResultBufFn.Call(store)
	if err != nil {
		log.Fatalf("Go: call wasm function to get result buffer failed: %s\n", err)
	}

	startIndex := int(ptr.(int32))

	fmt.Printf("Go: calling into sandbox get_result_count() ...\n")
	getResultCountFn := instance.GetFunc(store, "get_result_count")
	resultCountVal, err := getResultCountFn.Call(store)
	if err != nil {
		log.Fatalf("Go: call wasm function to get result count failed: %s\n", err)
	}

	resultCount := int(resultCountVal.(int32))

	fmt.Printf("Go: calling into sandbox get_result_struct_size() ...\n")
	getResultStructSizeFn := instance.GetFunc(store, "get_result_struct_size")
	resultStructSizeVal, err := getResultStructSizeFn.Call(store)
	if err != nil {
		log.Fatalf("Go: call wasm function to get result count failed: %s\n", err)
	}

	resultStructSize := int(resultStructSizeVal.(int32))

	fmt.Printf("Go: sandbox processed %d events, result stored in buf addr %d, result struct size %d\n",
		resultCount, startIndex, resultStructSize)

	decodeParsedResult(instance, store, startIndex, resultCount, resultStructSize)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func printImportsExports(module *wasmtime.Module) {
	fmt.Printf("Go: Imports & Exports of wasm module:\n")

	for _, i := range module.Imports() {
		t := i.Type().FuncType()
		if t != nil {
			fmt.Printf("\tImports: %v, %v, %v\n", *i.Name(), t.Params(), t.Results())
			continue
		}

		fmt.Printf("\tImports: %v, %v, %v\n", *i.Name(), nil, nil)
	}

	for _, i := range module.Exports() {
		t := i.Type().FuncType()
		if t != nil {
			fmt.Printf("\tExports: %v, %v, %v\n", i.Name(), t.Params(), t.Results())
			continue
		}

		fmt.Printf("\tExports: %v, %v, %v\n", i.Name(), nil, nil)
	}
}

// getArgs returns an argument list that will be passed to the wasm program,
// where the first item in the list is the "argc" argument, and the remaining
// ones are "argv[]".
func getArgs() []string {
	events := getEvents()

	// About argc/argv[] size:
	// https://ubuntuforums.org/archive/index.php/t-1186554.html
	return append([]string{string(len(events))}, events...)
}

// getEvents return a http event list for test
func getEvents() []string {
	return []string{
		"GET /api/v1/bpf HTTP/1.1\r\nHost: tricorder.dev\r\nCookie: cookie\r\n\r\n",
		"PUT /api/v2/wasm HTTP/1.1\r\nHost: tricorder.dev\r\n\r\n",
		"GET /wp-content/uploads/2010/03/hello-kitty-darth-vader-pink.jpg HTTP/1.1\r\n" +
			"Host: www.kittyhell.com\r\n" +
			"User-Agent: Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10.6; ja-JP-mac; rv:1.9.2.3) Gecko/20100401 Firefox/3.6.3 " +
			"Pathtraq/0.9\r\n" +
			"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\n" +
			"Accept-Language: ja,en-us;q=0.7,en;q=0.3\r\n" +
			"Accept-Encoding: gzip,deflate\r\n" +
			"Accept-Charset: Shift_JIS,utf-8;q=0.7,*;q=0.7\r\n" +
			"Keep-Alive: 115\r\n" +
			"Connection: keep-alive\r\n" +
			"Cookie: wp_ozh_wsa_visits=2; wp_ozh_wsa_visit_lasttime=xxxxxxxxxx; " +
			"__utma=xxxxxxxxx.xxxxxxxxxx.xxxxxxxxxx.xxxxxxxxxx.xxxxxxxxxx.x; " +
			"__utmz=xxxxxxxxx.xxxxxxxxxx.x.x.utmccn=(referral)|utmcsr=reader.livedoor.com|utmcct=/reader/|utmcmd=referral\r\n" +
			"\r\n",
	}
}

func inspectSandboxMemory(instance *wasmtime.Instance, store *wasmtime.Store) {
	mem := instance.GetExport(store, "memory").Memory()
	fmt.Printf("Go: wasm sandbox memory info: size %d, data size %d\n",
		mem.Size(store),
		mem.DataSize(store),
	)

	// Print the entire content of the sandbox memory, if you like
	if false {
		buf := mem.UnsafeData(store)
		fmt.Printf("Go: wasm sandbox memory content: %s\n", string(buf))
	}
}

func decodeParsedResult(instance *wasmtime.Instance, store *wasmtime.Store, startIndex, resultCount, resultStructSize int) error {
	mem := instance.GetExport(store, "memory").Memory()
	buf := mem.UnsafeData(store)

	fmt.Printf("Go: parsed result decoded from sandbox output:\n")
	for i := 0; i < resultCount; i++ {
		start := startIndex + i*resultStructSize
		method := string(buf[start : start+maxHttpMethodLen])

		start += maxHttpMethodLen
		path := string(buf[start : start+maxHttpPathLen])

		fmt.Printf("Method: %s, Path %s\n", method, path)
	}

	return nil
}

// How to copy data to sandbox, save for later reference
// mem0, err := NewMemory(store, NewMemoryType(2, true, 3))
// copy(mem0.UnsafeData(store)[2:3], []byte{100})
