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

package cgo

/*
#include "event.h"
*/
import "C"
import "unsafe"

type CEvent = C.Event_t

func GetCEvent(data []byte) *C.Event_t {
	return (*C.Event_t)(unsafe.Pointer(&data[0]))
}

func CEventToBytes(data *C.Event_t) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(data)), unsafe.Sizeof(*data))
}
