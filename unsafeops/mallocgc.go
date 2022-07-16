package unsafeops

import (
	_ "reflect"
	"unsafe"
)

//go:linkname MallocGC runtime.mallocgc
func MallocGC(size uintptr, typ uintptr, needzero bool) unsafe.Pointer

//go:linkname NewObject reflect.unsafe_New
func NewObject(typ uintptr) unsafe.Pointer

//go:linkname NewArray runtime.newarray
func NewArray(typ uintptr, n int) unsafe.Pointer
