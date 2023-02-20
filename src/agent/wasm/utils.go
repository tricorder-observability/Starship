// Copyright (C) 2023  tricorder-observability
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package wasm

import (
	"fmt"

	"github.com/bytecodealliance/wasmtime-go/v3"

	"github.com/tricorder/src/utils/file"
)

// wat2Wasm reads the content of a file with WASM text-format code, compiles and returns the generated WASM bytecode.
func wat2Wasm(filePath string) ([]byte, error) {
	wat, err := file.Read(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open WASM source file '%s', error: %v", filePath, err)
	}
	return wasmtime.Wat2Wasm(wat)
}

// TODO(yzhao): Use generic.
func unpackInt32Intf(intf interface{}) (int32, error) {
	val, ok := intf.(int32)
	if !ok {
		return 0, fmt.Errorf("expect int32, but got %v", intf)
	}
	return val, nil
}

// Run a function with signature `func() int32`
func runVoidInt32(module *Module, fnName string) (int32, error) {
	intf, err := module.Run(fnName)
	if err != nil {
		return 0, err
	}
	i32Val, ok := intf.(int32)
	if !ok {
		return 0, fmt.Errorf("expect %s() to return int32, but got %v", fnName, intf)
	}
	return i32Val, nil
}

// callU32I32 wraps a call to a func (int32) int32
func callU32I32(module *Module, fn string, i32 int32) (int32, error) {
	ret, err := module.Run1(fn, i32)
	if err != nil {
		return 0, err
	}
	return unpackInt32Intf(ret)
}
