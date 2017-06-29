package util

import (
	"reflect"
	"unsafe"
)

// convert string to []byte without memory copy
func String2ByteSlice(s string) []byte {
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	var result []byte
	shdr := (*reflect.SliceHeader)(unsafe.Pointer(&result))
	shdr.Data = hdr.Data
	shdr.Len = hdr.Len
	shdr.Cap = hdr.Len
	return result
}

