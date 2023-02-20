// Copyright (C) 2023  Tricorder Observability
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
)

// These function names are defined in libs/c/io.h
// WASM code must include io.h to be able to work with these APIs.
// The model is host-environment-driven WASM memory management:
// * Memory is entirely managed by functions inside WASM runtime
// * All APIs are defined inside WASM runtime
// * But the process is driven by host environment

// TODOs:
// * Need a way to signal failure when calling malloc_{input,output}_buf() multiple times.

const (
	MallocInputBufFn  = "malloc_input_buf"
	MallocOutputBufFn = "malloc_output_buf"

	FreeInputBufFn  = "free_input_buf"
	FreeOutputBufFn = "free_output_buf"

	GetInputBufFn  = "get_input_buf"
	GetOutputBufFn = "get_output_buf"

	GetInputBufLenFn = "get_input_buf_len"
	GetInputBufCapFn = "get_input_buf_cap"

	GetOutputBufLenFn = "get_output_buf_len"
	GetOutputBufCapFn = "get_output_buf_cap"

	SetInputBufLenFn  = "set_input_buf_len"
	SetOutputBufLenFn = "set_output_buf_len"
)

func MallocInputBuf(module *Module, capacity int32) (int32, error) {
	return callU32I32(module, MallocInputBufFn, capacity)
}

func MallocOutputBuf(module *Module, capacity int32) (int32, error) {
	return callU32I32(module, MallocOutputBufFn, capacity)
}

func FreeInputBuf(module *Module) error {
	_, err := module.Run(FreeInputBufFn)
	return err
}

func FreeOutputBuf(module *Module) error {
	_, err := module.Run(FreeOutputBufFn)
	return err
}

func GetInputBuf(module *Module) (int32, error) {
	return runVoidInt32(module, GetInputBufFn)
}

func GetInputBufCap(module *Module) (int32, error) {
	return runVoidInt32(module, GetInputBufCapFn)
}

func setInputBufLen(module *Module, length int32) error {
	_, err := module.Run1(SetInputBufLenFn, length)
	return err
}

func getInputBufLen(module *Module) (int32, error) {
	return runVoidInt32(module, GetInputBufLenFn)
}

// ClearInputBuf resets length so that the entire buffer is available for writing.
func ClearInputBuf(module *Module) error {
	return setInputBufLen(module, int32(0))
}

func CopyToInputBuf(module *Module, data []byte) error {
	wasmMem := module.memorySlice()
	inputBufOffset, err := GetInputBuf(module)
	if err != nil {
		return fmt.Errorf("while copying data to input buffer, failed to get input buffer, error: %v", err)
	}
	inputBufCapacity, err := GetInputBufCap(module)
	if err != nil {
		return fmt.Errorf("while copying data to input buffer, failed to get input buffer capacity, error: %v", err)
	}
	dataLen := int32(len(data))
	if dataLen > inputBufCapacity {
		return fmt.Errorf("input data size larger than input buffer capacity %d vs %d", dataLen, inputBufCapacity)
	}
	copy(wasmMem[inputBufOffset:inputBufOffset+dataLen], data)
	return setInputBufLen(module, dataLen)
}

func getOutputBuf(module *Module) (int32, error) {
	return runVoidInt32(module, GetOutputBufFn)
}

func getOutputBufLen(module *Module) (int32, error) {
	return runVoidInt32(module, GetOutputBufLenFn)
}

func setOutputBufLen(module *Module, length int32) error {
	_, err := module.Run1(SetOutputBufLenFn, length)
	return err
}

func getOutputBufCap(module *Module) (int32, error) {
	return runVoidInt32(module, GetOutputBufCapFn)
}

// ClearOutputBuf resets length so that the entire buffer is available for writing.
func ClearOutputBuf(module *Module) error {
	return setOutputBufLen(module, int32(0))
}

func ReadFromOutputBuf(module *Module) ([]byte, error) {
	wasmMem := module.memorySlice()
	outputBufOffset, err := getOutputBuf(module)
	if err != nil {
		return nil, fmt.Errorf("while reading from putput buffer, failed to get output buffer, error: %v", err)
	}
	outputBufLen, err := getOutputBufLen(module)
	if err != nil {
		return nil, fmt.Errorf("while reading from putput buffer, failed to get output buffer length, error: %v", err)
	}
	data := make([]byte, outputBufLen)
	copy(data, wasmMem[outputBufOffset:outputBufOffset+outputBufLen])
	return data, nil
}
