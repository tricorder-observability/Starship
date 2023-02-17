package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bytecodealliance/wasmtime-go/v3"
)

func main() {
	// Configure our `Store`, but be sure to use a `Config` that enables the
	// wasm multi-value feature since it's not stable yet.
	config := wasmtime.NewConfig()
	config.SetWasmMultiValue(true)
	store := wasmtime.NewStore(wasmtime.NewEngineWithConfig(config))

	source, err := wasmSourceCodeFromFile("./hello-multi-args.wat")
	if err != nil {
		log.Fatal(err)
	}

	wasm, err := wasmtime.Wat2Wasm(source)
	if err != nil {
		log.Fatal(err)
	}

	module, err := wasmtime.NewModule(store.Engine, wasm)
	if err != nil {
		log.Fatal(err)
	}

	callback := wasmtime.WrapFunc(store, func(a int32, b int64) (int64, int32) {
		return b + 1, a + 1
	})

	instance, err := wasmtime.NewInstance(store, module, []wasmtime.AsExtern{callback})
	if err != nil {
		log.Fatal(err)
	}

	g := instance.GetFunc(store, "g")

	results, err := g.Call(store, 1, 3)
	if err != nil {
		log.Fatal(err)
	}
	arr := results.([]wasmtime.Val)
	a := arr[0].I64()
	b := arr[1].I32()
	fmt.Printf("> %d %d\n", a, b)

	if a != 4 {
		log.Fatal("unexpected value for a")
	}
	if b != 2 {
		log.Fatal("unexpected value for b")
	}

	roundTripMany := instance.GetFunc(store, "round_trip_many")
	results, err = roundTripMany.Call(store, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	if err != nil {
		log.Fatal(err)
	}
	arr = results.([]wasmtime.Val)

	for i := 0; i < len(arr); i++ {
		fmt.Printf(" %d", arr[i].Get())
		if arr[i].I64() != int64(i) {
			log.Fatal("unexpected value for arr[i]")
		}
	}
}

func wasmSourceCodeFromFile(file string) (string, error) {
	rawBytes, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("failed to open WASM source file: %s\n", err)
	}

	return string(rawBytes), nil
}
