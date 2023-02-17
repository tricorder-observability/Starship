package main

import (
	"fmt"
	"os"

	"github.com/bytecodealliance/wasmtime-go/v3"
)

func main() {
	// Operations in wasmtime require a contextual `store` argument to share
	store := wasmtime.NewStore(wasmtime.NewEngine())

	// Convert WebAssembly text format to the binary format.
	source, err := wasmSourceCodeFromFile("./hello-world.wat")
	check(err)
	wasm, err := wasmtime.Wat2Wasm(source)
	check(err)

	// Compile binary wasm into a `*Module` which represents compiled JIT code.
	module, err := wasmtime.NewModule(store.Engine, wasm)
	check(err)

	// `hello.wat` file imports one item, so we create that function here.
	item := wasmtime.WrapFunc(store, func() {
		fmt.Println("Hello from Go!")
	})

	// Instantiate a module to link in all wasm imports. We've got one import so we pass that in here.
	instance, err := wasmtime.NewInstance(store, module, []wasmtime.AsExtern{item})
	check(err)

	// Lookup the `run()` function defined in `hello.wat` and call it.
	run := instance.GetFunc(store, "run")
	if run == nil {
		panic("not a function")
	}

	_, err = run.Call(store)
	check(err)
}

func wasmSourceCodeFromFile(file string) (string, error) {
	rawBytes, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("failed to open WASM source file: %s\n", err)
	}

	return string(rawBytes), nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
