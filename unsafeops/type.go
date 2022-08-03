package unsafeops

import (
	"reflect"
	"unsafe"
)

var rtType uintptr

func init() {
	t := reflect.TypeOf(unsafe.Pointer(nil))
	rtEface := (*Eface)(unsafe.Pointer(NoEscape(&t)))
	rtType = uintptr(rtEface.Type)
}

func UnsafeType(typ reflect.Type) uintptr {
	eface := (*Eface)(unsafe.Pointer(NoEscape(&typ)))
	return uintptr(eface.Data)
}

func ReflectType(typ uintptr) reflect.Type {
	eface := MakeEface(unsafe.Pointer(typ), rtType)
	return *(*reflect.Type)(unsafe.Pointer(&eface))
}

//go:linkname IfaceIndir reflect.ifaceIndir
func IfaceIndir(t uintptr) bool
