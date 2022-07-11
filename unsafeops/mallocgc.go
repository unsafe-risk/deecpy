package unsafeops

import "unsafe"

//go:linkname NewObject runtime.newobject
func NewObject(typ uintptr) unsafe.Pointer

//go:linkname NewArray runtime.newarray
func NewArray(typ uintptr, n int) unsafe.Pointer
