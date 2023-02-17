// Demo code to test wasmtime externref and host side new memory feature with WAT code.
package main

import (
	"fmt"
	"log"

	"github.com/bytecodealliance/wasmtime-go/v3"
)

func main() {
	TestMultiMemoryImported()

	// ExampleVal_Externref()

	fmt.Printf("Finish\n")
}

func multiMemoryStore() *wasmtime.Store {
	config := wasmtime.NewConfig()
	config.SetWasmMultiMemory(true)
	return wasmtime.NewStore(wasmtime.NewEngineWithConfig(config))
}

// Seems there is no instruction to get the memory address of a given "memory":
//
// * WAT instructions: https://developer.mozilla.org/en-US/docs/WebAssembly/Reference
// * https://coinexsmartchain.medium.com/wasm-introduction-part-3-memory-7426f19c9624

// WASM linear memory:
// https://github.com/sunfishcode/wasm-reference-manual/blob/master/WebAssembly.md#linear-memories

// Status:
// We could read/write with memory ID specified, but could not get the absolute address for this memory.
// So, C program could not get the addressable pointer for reading/writing data without being aware of the multi-memory fact.

// NewMemory returns a new linear memory space, which is different from the default one:
// * default memory space: [0, 4GB]
// * input memory space:   [0, 4GB]
// * output memory space:  [0, 4GB]
//
// To make C programs to be able to manipulate multiple memory spaces, the C
// programs themselves must be multi-memory aware: they must know that they
// will be compiled into wasm code; or, relying on the clang/llvm tool chain
// to transparently convert the related functions/syscalls (open,read,write,...)
// to the memory space aware version, but there is no such thing right now.

func TestMultiMemoryImported() {
	wasm, err := wasmtime.Wat2Wasm(`
    (module
      (import "" "m0" (memory 1))
      (import "" "m1" (memory $m 1))
      (func (export "load1") (result i32)
        i32.const 2
        i32.load8_s $m
      )
      (func (export "load2") (result i32)
	    i32.const 0
		return
      )
    )`)
	if err != nil {
		log.Fatal(err)
	}
	store := multiMemoryStore()

	mem0, err := wasmtime.NewMemory(store, wasmtime.NewMemoryType(1, true, 3))
	mem1, err := wasmtime.NewMemory(store, wasmtime.NewMemoryType(2, true, 4))

	module, err := wasmtime.NewModule(store.Engine, wasm)
	instance, err := wasmtime.NewInstance(store, module, []wasmtime.AsExtern{mem0, mem1})

	copy(mem1.UnsafeData(store)[2:3], []byte{100})

	res, err := instance.GetFunc(store, "load1").Call(store)
	fmt.Printf("res: %v\n", res.(int32))

	res, err = instance.GetFunc(store, "load2").Call(store)
	fmt.Printf("res: %v\n", res.(int32))
}

func ExampleVal_Externref() {
	config := wasmtime.NewConfig()
	config.SetWasmReferenceTypes(true)
	store := wasmtime.NewStore(wasmtime.NewEngineWithConfig(config))
	wasm, err := wasmtime.Wat2Wasm(`
	(module
	  (table $table (export "table") 10 externref)

	  (global $global (export "global") (mut externref) (ref.null extern))

	  (func (export "func") (param externref) (result externref)
	    local.get 0
	  )
	)
	`)
	if err != nil {
		log.Fatal(err)
	}
	module, err := wasmtime.NewModule(store.Engine, wasm)
	if err != nil {
		log.Fatal(err)
	}
	instance, err := wasmtime.NewInstance(store, module, []wasmtime.AsExtern{})
	if err != nil {
		log.Fatal(err)
	}
	// Create a new `externref` value.
	value := wasmtime.ValExternref("Hello, World!")
	// The `externref`'s wrapped data should be the string "Hello, World!".
	externRef := value.Externref()
	if externRef != "Hello, World!" {
		log.Fatal("unexpected value")
	}
	// Lookup the `table` export.
	table := instance.GetExport(store, "table").Table()
	// Set `table[3]` to our `externref`.
	err = table.Set(store, 3, value)
	if err != nil {
		log.Fatal(err)
	}
	// `table[3]` should now be our `externref`.
	tableValue, err := table.Get(store, 3)
	if err != nil {
		log.Fatal(err)
	}
	if tableValue.Externref() != externRef {
		log.Fatal("unexpected value in table")
	}
	// Lookup the `global` export.
	global := instance.GetExport(store, "global").Global()
	// Set the global to our `externref`.
	err = global.Set(store, value)
	if err != nil {
		log.Fatal(err)
	}
	// Get the global, and it should return our `externref` again.
	globalValue := global.Get(store)
	if globalValue.Externref() != externRef {
		log.Fatal("unexpected value in global")
	}
	// Lookup the `func` export.
	fn := instance.GetFunc(store, "func")
	// And call it!
	result, err := fn.Call(store, value)
	if err != nil {
		log.Fatal(err)
	}
	// `func` returns the same reference we gave it, so `results` should be
	// our `externref`.
	if result != externRef {
		log.Fatal("unexpected value from func")
	}
	// Output:
	//
	fmt.Println(result)
}
