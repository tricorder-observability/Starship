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
