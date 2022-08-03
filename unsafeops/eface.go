package unsafeops

import (
	"unsafe"
)

type Eface struct {
	Type uintptr
	Data unsafe.Pointer
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
	return EfaceOf(NoEscape(v)).Type
}

func DataOf(v *any) unsafe.Pointer {
	return EfaceOf(NoEscape(v)).Data
}

func MakeEface(data unsafe.Pointer, t uintptr) any {
	return *(*any)(unsafe.Pointer(&Eface{
		Type: t,
		Data: data,
	}))
}
