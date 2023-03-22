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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tricorder/src/agent/wasm/programs/cgo"
	bazelutils "github.com/tricorder/src/testing/bazel"
	"github.com/tricorder/src/utils/file"
)

// Tests that copy_input_to_output() in WASM module can copy data to the WASM runtime, repetitively.
func TestWasmMemoryIO(t *testing.T) {
	assert := assert.New(t)
	testFilePath := "modules/sample_json/sample_json.wasm"

	wasmFilePath := bazelutils.TestFilePath(testFilePath)
	wasm, err := file.ReadBin(wasmFilePath)
	assert.Nil(err)

	module, err := NewWasiModule(wasm, []string{})
	assert.Nil(err)

	inputBufOffset, err := MallocInputBuf(module, 100)
	assert.Nil(err)
	assert.NotEqual(0, inputBufOffset)

	inputLen, err := getInputBufLen(module)
	assert.Nil(err)
	assert.Equal(int32(0), inputLen)
	capacity, err := GetInputBufCap(module)
	assert.Nil(err)
	assert.Equal(int32(100), capacity)

	outputBufOffset, err := MallocOutputBuf(module, 100)
	assert.Nil(err)
	assert.NotEqual(0, outputBufOffset)

	offset, err := GetInputBuf(module)
	assert.Nil(err)
	assert.NotEqual(100, offset)

	err = CopyToInputBuf(module, []byte("Hello WASM!"))
	assert.Nil(err)

	ret, err := module.Run("copy_input_to_output")
	assert.Nil(ret)
	assert.Nil(err)

	outputLen, err := getOutputBufLen(module)
	assert.Nil(err)
	assert.Equal(int32(11), outputLen)
	capacity, err = getOutputBufCap(module)
	assert.Nil(err)
	assert.Equal(int32(100), capacity)

	data, err := ReadFromOutputBuf(module)
	assert.Nil(err)
	assert.Equal(11, len(data))
	assert.Equal("Hello WASM!", string(data))

	err = FreeInputBuf(module)
	assert.Nil(err)
	err = FreeOutputBuf(module)
	assert.Nil(err)

	inputLen, err = getInputBufLen(module)
	assert.Nil(err)
	assert.Equal(int32(0), inputLen)
	capacity, err = GetInputBufCap(module)
	assert.Nil(err)
	assert.Equal(int32(0), capacity)
}

// Tests that the same C struct's object has different memory layout in host environment and WASM.
// This test works by reading a memory buffer from WASM, and cast it to Cgo type with the same definition.
func TestCStructMemLayout(t *testing.T) {
	assert := assert.New(t)

	testFilePath := "src/agent/wasm/programs/struct_test.wasm"

	wasmFilePath := bazelutils.TestFilePath(testFilePath)
	wasm, err := file.ReadBin(wasmFilePath)
	assert.Nil(err)

	module, err := NewWasiModule(wasm, []string{})
	assert.Nil(err)

	inputBufOffset, err := MallocInputBuf(module, 100)
	assert.Nil(err)
	assert.NotEqual(0, inputBufOffset)

	cInts := cgo.CInts{}
	cInts.A = 10
	cInts.B = 11
	cInts.C = 12
	cInts.D = 13
	cInts.E = 14

	err = CopyToInputBuf(module, cgo.CIntsToBytes(&cInts))
	assert.Nil(err)

	outputBufOffset, err := MallocOutputBuf(module, 100)
	assert.Nil(err)
	assert.NotEqual(0, outputBufOffset)

	// Invoke function to copy a C struct with the same definition as C.Ints_t to output buffer.
	// The C.Ints_t value inside WASM assigned different value than cInts.
	// TODO: Change to compute a new value based on the input, so we can verify the input has been
	// passed correctly to the WASM runtime, right now, we can only verify output has been copied
	// correctly from WASM runtime to the host environment.
	ret, err := module.Run("write_ints_to_output")
	assert.Nil(ret)
	assert.Nil(err)

	data, err := ReadFromOutputBuf(module)
	assert.Nil(err)
	assert.Equal(24, len(data))

	cIntsPtr := cgo.GetCInts(data)
	assert.Equal(cgo.CInts{
		A: 1,
		B: 2,
		C: 3,
		D: 4,
		E: 5,
	}, *cIntsPtr)

	err = FreeOutputBuf(module)
	assert.Nil(err)
}

// Tests that the same C struct's object has different memory layout in host environment and WASM.
// This test works by reading a memory buffer from WASM, and cast it to Cgo type with the same definition.
func TestCStructMemMarshalAndUnmarshal(t *testing.T) {
	assert := assert.New(t)

	testFilePath := "src/agent/wasm/programs/struct_test.wasm"

	wasmFilePath := bazelutils.TestFilePath(testFilePath)
	wasm, err := file.ReadBin(wasmFilePath)
	assert.Nil(err)

	module, err := NewWasiModule(wasm, []string{})
	assert.Nil(err)

	inputBufOffset, err := MallocInputBuf(module, 100)
	assert.Nil(err)
	assert.NotEqual(0, inputBufOffset)

	cEvent := cgo.CEvent{}
	cEvent.X = 1.01
	cEvent.Y = 2.02
	cEvent.Z = 3
	cEvent.Comm[0] = 5

	err = CopyToInputBuf(module, cgo.CEventToBytes(&cEvent))
	assert.Nil(err)

	outputBufOffset, err := MallocOutputBuf(module, 100)
	assert.Nil(err)
	assert.NotEqual(0, outputBufOffset)

	// Invoke function to copy a C struct with the same definition as C.Ints_t to output buffer.
	// The C.Ints_t value inside WASM assigned different value than cInts.
	ret, err := module.Run("write_marshaled_struct_to_output")
	assert.Nil(ret)
	assert.Nil(err)

	data, err := ReadFromOutputBuf(module)
	assert.Nil(err)

	cEventPtr := cgo.GetCEvent(data)
	eventRes := cgo.CEvent{
		A: 1,
		X: 0.03,
		Y: 0.04,
		Z: 5,
	}
	eventRes.Comm[0] = 1
	eventRes.Comm[5] = 10
	assert.Equal(eventRes, *cEventPtr)

	err = FreeOutputBuf(module)
	assert.Nil(err)
}
