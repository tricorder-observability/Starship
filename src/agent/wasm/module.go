// Package wasm wraps Wasmtime-Go's Golang binding of the C/C++ Wasmtime API.
package wasm

import (
	"fmt"
	"strconv"

	wasmtime "github.com/bytecodealliance/wasmtime-go/v3"
)

// Module wraps the data structures of wasmtime, and corresponds to a piece of wasm byte code.
// The wasm byte code could be:
//   - WASI style, which utilizes system interfaces, for example they need to call libc:
//     https://github.com/WebAssembly/wasi-libc
//   - WASM style, or non-WASI style, which does not utilize system interface.
type Module struct {
	// The raw wasm binary-format byte code.
	wasm []byte

	// Following are basic data structure for wrapping a piece of WASM byte code.
	engine *wasmtime.Engine
	store  *wasmtime.Store
	module *wasmtime.Module

	// Can be recreated for repeated invocation.
	instance *wasmtime.Instance

	// Only needed for WASI-style wasm module.
	linker     *wasmtime.Linker
	wasiConfig *wasmtime.WasiConfig
}

func (module *Module) newWasiConfig(args []string) error {
	err := module.linker.DefineWasi()
	if err != nil {
		return fmt.Errorf("failed to define WASI, error: %v", err)
	}

	wasiConfig := wasmtime.NewWasiConfig()
	wasiConfig.InheritEnv()
	wasiConfig.InheritStdout()
	wasiConfig.InheritStdin()
	wasiConfig.InheritStderr()

	// In the argument list passed to the wasm program, the first item in the list is the "argc" argument, and the
	// remaining ones are "argv[]".
	args = append([]string{strconv.Itoa(len(args))}, args...)
	wasiConfig.SetArgv(args)

	module.wasiConfig = wasiConfig

	module.store.SetWasi(wasiConfig)

	return nil
}

// NewWasmInstance creates an WASM instance that has a linker to link in all wasm imports.
// The linker can pass command line style arguments to the `_start` function of the WASM module.
func (module *Module) NewWasmInstance(fn func()) error {
	instance, err := wasmtime.NewInstance(module.store, module.module, []wasmtime.AsExtern{module.wrapFn(fn)})
	if err != nil {
		return fmt.Errorf("failed to create new instance, error: %v", err)
	}
	module.instance = instance
	return nil
}

// NewWasiInstance creates an instance for WASI style invocation
// `argv` must be a command line style array of arguments, i.e., 1st arg is the number of argv, the rest are actual
// args, like main(argc, argv).
func (module *Module) NewWasiInstance(argv []string) error {
	module.linker = wasmtime.NewLinker(module.engine)

	err := module.newWasiConfig(argv)
	if err != nil {
		return fmt.Errorf("failed to create WASI config, error: %v", err)
	}
	instance, err := module.linker.Instantiate(module.store, module.module)
	if err != nil {
		return fmt.Errorf("failed to create new instance with linker, error: %v", err)
	}

	module.instance = instance
	return nil
}

// newBasicModule returns a Module that has basic fields initialized.
func newBasicModule(wasm []byte) (*Module, error) {
	module := new(Module)

	module.wasm = wasm
	module.engine = wasmtime.NewEngine()
	module.store = wasmtime.NewStore(module.engine)
	// Compile binary wasm into a `*Module` which represents compiled JIT code.
	wasmModule, err := wasmtime.NewModule(module.engine, wasm)
	if err != nil {
		return nil, fmt.Errorf("failed to create new module, error: %v", err)
	}
	module.module = wasmModule

	return module, nil
}

// NewWasiModule returns a Module for a WASI style invocation for the `_start` function.
// `_start` function corresponds to the usual main() function or other entry point.
// It seems args can only be passed to the wasmModule from the beginning, which makes it quite inefficient to change to
// use different args to the `_start` function.
func NewWasiModule(wasm []byte, argv []string) (*Module, error) {
	module, err := newBasicModule(wasm)
	if err != nil {
		return nil, fmt.Errorf("failed to create new module, error: %v", err)
	}

	err = module.NewWasiInstance(argv)
	if err != nil {
		return nil, fmt.Errorf("failed to create new instance with linker, error: %v", err)
	}

	return module, nil
}

// NewWasmModule returns a Module for a piece of non-WASI byte code.
// Such byte code has no use of system interfaces.
func NewWasmModule(wasm []byte, fn func()) (*Module, error) {
	module, err := newBasicModule(wasm)
	if err != nil {
		return nil, fmt.Errorf("failed to create new module, error: %v", err)
	}

	err = module.NewWasmInstance(fn)
	if err != nil {
		return nil, fmt.Errorf("failed to create new instance, error: %v", err)
	}
	return module, nil
}

func (module *Module) wrapFn(fn func()) *wasmtime.Func {
	return wasmtime.WrapFunc(module.store, fn)
}

// memorySlice returns a byte slide that represents the unsafe memory of the instance.
func (module *Module) memorySlice() []byte {
	const memoryExportName = "memory"
	mem := module.instance.GetExport(module.store, memoryExportName).Memory()
	return mem.UnsafeData(module.store)
}

// Run invokes the WASM function with the input name.
func (module *Module) Run(fnName string) (interface{}, error) {
	run := module.instance.GetFunc(module.store, fnName)
	if run == nil {
		return 0, fmt.Errorf("could not find function '%s'", fnName)
	}
	return run.Call(module.store)
}

func (module *Module) Run1(fnName string, arg int32) (interface{}, error) {
	run := module.instance.GetFunc(module.store, fnName)
	if run == nil {
		return 0, fmt.Errorf("could not find function '%s'", fnName)
	}
	return run.Call(module.store, arg)
}

// DebugString() returns a string that describes the module's important information.
func (module *Module) DebugString() []string {
	lines := make([]string, 0, 10)
	for _, i := range module.module.Imports() {
		t := i.Type().FuncType()
		if t != nil {
			desc := fmt.Sprintf("Imports: %v, %v, %v", *i.Name(), t.Params(), t.Results())
			lines = append(lines, desc)
			continue
		}
		desc := fmt.Sprintf("Imports: %v, %v, %v", *i.Name(), nil, nil)
		lines = append(lines, desc)
	}
	for _, e := range module.module.Exports() {
		t := e.Type().FuncType()
		if t != nil {
			desc := fmt.Sprintf("Exports: %v, %v, %v", e.Name(), t.Params(), t.Results())
			lines = append(lines, desc)
			continue
		}
		desc := fmt.Sprintf("Exports: %v, %v, %v", e.Name(), nil, nil)
		lines = append(lines, desc)
	}
	return lines
}
