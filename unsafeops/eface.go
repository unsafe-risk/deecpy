package unsafeops

import (
	"unsafe"
)

type Eface struct {
	_type uintptr
	data  unsafe.Pointer
}

func EfaceOf(ep *any) *Eface {
	return (*Eface)(unsafe.Pointer(ep))
}

//nolint:uintptr
func NoEscape[T any](x *T) *T {
	v := uintptr(unsafe.Pointer(x))
	//nolint:staticcheck
	return (*T)(unsafe.Pointer(v ^ 0))
}

func TypeID(v *any) uintptr {
	return EfaceOf(NoEscape(v))._type
}

func DataOf(v *any) unsafe.Pointer {
	return EfaceOf(NoEscape(v)).data
}

func MakeEface(data unsafe.Pointer, t uintptr) any {
	return *(*any)(unsafe.Pointer(&Eface{
		_type: t,
		data:  data,
	}))
}
