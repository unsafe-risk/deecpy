package unsafeops

import (
	"reflect"
	"unsafe"
)

func UnsafeType(typ reflect.Type) uintptr {
	eface := (*Eface)(unsafe.Pointer(NoEscape(&typ)))
	return uintptr(eface.data)
}
