// Demo: passing data from host to wasm module with host side allocated wasm memory
//
// Status:
//   2023-01-14: not working (tooling chain not ready, see README.md for more info)
package main

import (
	"fmt"
	"log"

	"github.com/bytecodealliance/wasmtime-go/v3"
)

const (
	wasmFile = "./libpico3.wasm"

	// Functions exported by libpico.wasm
	fnNameAllocateInputOutputBufs = "allocate_input_output_bufs"
	fnNameGetInputBuf             = "get_input_buf"
	fnNameGetOutputBuf            = "get_output_buf"
	fnNameGetOutputItemCount      = "get_output_item_count"
	fnNameGetOutputItemSize       = "get_output_item_size"
	fnNamePicoParseEvents         = "pico_parse_events"

	debug = true

	// Max event (req) size, should be consistent with the c/wasm side definitions.
	maxEventSize = 4096

	// Max field lengths in http header, should be consistent with the c/wasm side definitions.
	// See main.c for the latters.
	maxHttpMethodLen = 8
	maxHttpPathLen   = 120
)

func main() {
	// Create wasm module
	config := wasmtime.NewConfig()
	config.SetWasmMultiMemory(true)
	// config.SetWasmMemory64(true)
	engine := wasmtime.NewEngineWithConfig(config)
	module, err := wasmtime.NewModuleFromFile(engine, wasmFile)
	if err != nil {
		panic(err)
	}

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
	wasiConfig.InheritStdout()
	wasiConfig.InheritStdin()
	wasiConfig.InheritStderr()

	fmt.Printf("Go: creating new store\n")
	store := wasmtime.NewStore(engine)
	store.SetWasi(wasiConfig)

	// Test passing data via memory
	// Create new memory for passing data input wasm sandbox
	fmt.Printf("Go: creating new memory for input data\n")

	var minPageSize uint32 = 1 // 64KB per page
	var maxPageSize uint32 = 3
	hasMaxPageSize := true

	input_mem, err := wasmtime.NewMemory(store, wasmtime.NewMemoryType(minPageSize, hasMaxPageSize, maxPageSize))
	if err != nil {
		log.Fatal(err)
	}

	// TODO: the first parameter is "module string", but I don't know what "env"
	// stands for here, it just works.
	linker.FuncWrap("env", "get_input_buf", func() int32 {
		// Get host side memory address.
		// There is no API to get the offset of the memory address for wasm module:
		// https://docs.rs/wasmtime/0.17.0/wasmtime/struct.Memory.html
		data := input_mem.Data(store)

		// Also note that `data` is a int64 memory address, on return a i32
		// address to wasm module, the uppper 32bit will be discarded, leaving
		// the return vale non-sense.
		// Note: will int64 memory work? wasi-sdk/clang doesn't support wasm64 yet.
		fmt.Printf("Go: input memory address: host side %+v\n", data)

		// Conclusion: we can pass the host side memory address to wasm module,
		// but it doesn't make sense to the latter, that's why it will crash
		// on reading data from that address.
		return int32(uintptr(data))
	})

	output_mem, err := wasmtime.NewMemory(store, wasmtime.NewMemoryType(minPageSize, hasMaxPageSize, maxPageSize))
	if err != nil {
		log.Fatal(err)
	}

	linker.FuncWrap("env", "get_output_buf", func() int32 {
		data := output_mem.Data(store)
		fmt.Printf("Go: output memory address: host side %+v\n", data)
		return int32(uintptr(data))
	})

	// Create instance
	fmt.Printf("Go: instantiating sandbox\n")
	instance, err := linker.Instantiate(store, module)
	if err != nil {
		log.Fatal(err)
	}

	// // Copy to be processed events into sandbox
	numEvents := 0
	{
		startIndex := 0
		events := getEvents()
		numEvents = len(events)
		if err := copyInput(input_mem, instance, store, startIndex, events); err != nil {
			log.Fatal("Go: copy input data into sandbox memory failed: %s\n", err)
		}
	}

	// Call pico http parser in the sandbox (with input data provided via memory)
	fmt.Printf("Go: calling wasm func %s\n", fnNamePicoParseEvents)
	if _, err := instance.GetFunc(store, fnNamePicoParseEvents).Call(store, numEvents); err != nil {
		log.Fatalf("Call wasm func %s failed: %s\n", fnNamePicoParseEvents, err)
	}

	if debug {
		inspectSandboxMemory(instance, store)
	}

	// Get output memory address
	fmt.Printf("Go: calling wasm func %s ...\n", fnNameGetOutputBuf)
	ptr, err := instance.GetFunc(store, fnNameGetOutputBuf).Call(store)
	if err != nil {
		log.Fatalf("Go: call wasm func %s failed: %s\n", fnNameGetOutputBuf, err)
	}

	startIndex := int(ptr.(int32))

	// Get output item count
	fmt.Printf("Go: calling wasm func %s ...\n", fnNameGetOutputItemCount)
	outputItemCountVal, err := instance.GetFunc(store, fnNameGetOutputItemCount).Call(store)
	if err != nil {
		log.Fatalf("Go: call wasm func %s failed: %s\n", fnNameGetOutputItemCount, err)
	}

	outputItemCount := int(outputItemCountVal.(int32))

	// Get output item size
	fmt.Printf("Go: calling wasm func %s ...\n", fnNameGetOutputItemSize)
	outputItemSizeVal, err := instance.GetFunc(store, fnNameGetOutputItemSize).Call(store)
	if err != nil {
		log.Fatalf("Go: call wasm func %s failed: %s\n", fnNameGetOutputItemSize, err)
	}

	outputItemSize := int(outputItemSizeVal.(int32))

	// Decode output items
	fmt.Printf("Go: sandbox processed %d events, result stored in buf addr %d, result struct size %d\n",
		outputItemCount, startIndex, outputItemSize)

	decodeParsedResult(instance, store, startIndex, outputItemCount, outputItemSize)
}

// Print the import and export functions of the wasm module
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

// getEvents returns a http event list for testing
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

// For debugging and introspection
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

// Copy input data (events) into sandbox memory
func copyInput(mem *wasmtime.Memory, instance *wasmtime.Instance, store *wasmtime.Store, startIndex int, events []string) error {
	reqCount := len(events)
	buf := mem.UnsafeData(store)

	fmt.Printf("Go: copying events into sandbox memory\n")
	for i := 0; i < reqCount; i++ {
		event := events[i]
		start := startIndex + i*maxEventSize
		size := len(event)
		if size > maxEventSize {
			size = maxEventSize
			copy(buf[start:start+maxEventSize], event[:maxEventSize])
		} else {
			copy(buf[start:start+size], event[:])
		}
	}

	return nil
}

// Copy/decode output data from sandbox memory
func decodeParsedResult(instance *wasmtime.Instance, store *wasmtime.Store, startIndex, outputItemCount, outputItemSize int) error {
	mem := instance.GetExport(store, "memory").Memory()
	buf := mem.UnsafeData(store)

	fmt.Printf("Go: wasm outputs decoded:\n")
	for i := 0; i < outputItemCount; i++ {
		start := startIndex + i*outputItemSize
		method := string(buf[start : start+maxHttpMethodLen])

		start += maxHttpMethodLen
		path := string(buf[start : start+maxHttpPathLen])

		fmt.Printf("Method: %s, Path %s\n", method, path)
	}

	return nil
}
