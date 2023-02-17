package cgo

/*
#include "ints.h"
*/
import "C"

import "unsafe"

// testdata.C.Ints_t used outside of `testdata` package, cannot compile.
// This is the only way to workaround the issue.
type CInts = C.Ints_t
type CPackedInts = C.Packed_ints_t

// GetCInts converts a byte slice into a Cgo ints_t object.
func GetCInts(data []byte) *C.Ints_t {
	return (*C.Ints_t)(unsafe.Pointer(&data[0]))
}

// GetCInts converts a byte slice into a Cgo ints_t object.
func CIntsToBytes(data *C.Ints_t) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(data)), unsafe.Sizeof(*data))
}
