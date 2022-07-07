package unsafeops

import "unsafe"

type Eface struct {
	_type uintptr
	data  unsafe.Pointer
}

func EfaceOf(ep *any) *Eface {
	return (*Eface)(unsafe.Pointer(ep))
}

func noescape[T any](x *T) *T {
	v := uintptr(unsafe.Pointer(x))
	return (*T)(unsafe.Pointer(v ^ 0))
}

func TypeID(v *any) uintptr {
	return EfaceOf(noescape(v))._type
}
