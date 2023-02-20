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
