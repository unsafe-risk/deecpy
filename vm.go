package deecpy

import (
	"reflect"
	"unsafe"

	"github.com/unsafe-risk/deecpy/unsafeops"
)

func exec(dst, src unsafe.Pointer, inst *instructions, nocopy bool) {
L:
	for i := range inst.ops {
		switch v := inst.ops[i].(type) {
		case *opCopyMem:
			if nocopy {
				continue L
			}
			unsafeops.MemMove(
				unsafe.Add(dst, v.Offset),
				unsafe.Add(src, v.Offset),
				v.Size,
			)
		case *opPtrDup:
			srcPtr := *(*unsafe.Pointer)(unsafe.Add(src, v.Offset))
			if srcPtr == nil {
				*(*unsafe.Pointer)(unsafe.Add(dst, v.Offset)) = nil
				continue
			}
			*(*unsafe.Pointer)(unsafe.Add(dst, v.Offset)) = unsafeops.NewObject(v.UnsafeType)
			exec(*(*unsafe.Pointer)(unsafe.Add(dst, v.Offset)), srcPtr, v.SubInstructions, false)
		case *opPtrDupMem:
			srcPtr := *(*unsafe.Pointer)(unsafe.Add(src, v.Offset))
			if srcPtr == nil {
				*(*unsafe.Pointer)(unsafe.Add(dst, v.Offset)) = nil
				continue L
			}
			*(*unsafe.Pointer)(unsafe.Add(dst, v.Offset)) = unsafeops.NewObject(v.UnsafeType)
			if v.Size == 0 {
				continue L
			}
			unsafeops.MemMove(
				*(*unsafe.Pointer)(unsafe.Add(dst, v.Offset)),
				srcPtr,
				v.Size,
			)
		case *opArrayCopy:
			for i := uintptr(0); i < v.ArrayLen; i++ {
				exec(
					unsafe.Add(dst, v.Offset+i*v.ElemSize),
					unsafe.Add(src, v.Offset+i*v.ElemSize),
					v.SubInstructions,
					nocopy,
				)
			}
		case *opSliceCopy:
			s := *(*reflect.SliceHeader)(unsafe.Add(src, v.Offset))
			(*reflect.SliceHeader)(unsafe.Add(dst, v.Offset)).Cap = s.Cap
			(*reflect.SliceHeader)(unsafe.Add(dst, v.Offset)).Len = s.Len
			if uintptr(s.Cap)*v.ElemSize == 0 {
				continue L
			}
			(*reflect.SliceHeader)(unsafe.Add(dst, v.Offset)).Data = uintptr(unsafeops.NewArray(v.UnsafeElemType, s.Cap))

			sliceSrcData := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Add(src, v.Offset)).Data)
			sliceDstData := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Add(dst, v.Offset)).Data)

			for i := uintptr(0); i < uintptr(s.Cap); i++ {
				exec(
					unsafe.Add(sliceSrcData, uintptr(i)*v.ElemSize),
					unsafe.Add(sliceDstData, uintptr(i)*v.ElemSize),
					v.SubInstructions,
					false,
				)
			}
		case *opSliceCopyMem:
			s := *(*reflect.SliceHeader)(unsafe.Add(src, v.Offset))
			(*reflect.SliceHeader)(unsafe.Add(dst, v.Offset)).Cap = s.Cap
			(*reflect.SliceHeader)(unsafe.Add(dst, v.Offset)).Len = s.Len
			if s.Cap == 0 {
				continue L
			}
			(*reflect.SliceHeader)(unsafe.Add(dst, v.Offset)).Data = uintptr(unsafeops.NewArray(v.UnsafeElemType, s.Cap))
			unsafeops.MemMove(
				unsafe.Pointer((*reflect.SliceHeader)(unsafe.Add(dst, v.Offset)).Data),
				unsafe.Pointer(s.Data),
				uintptr(s.Cap)*v.ElemSize,
			)
		case *opMapDup:
			srcMap := (*unsafe.Pointer)(unsafe.Add(src, v.Offset))
			srcMapIface := unsafeops.MakeEface(*srcMap, v.MapUnsafeType)
			srcMapReflectValue := reflect.ValueOf(srcMapIface)
			keys := srcMapReflectValue.MapKeys()
			newMap := reflect.MakeMapWithSize(v.ReflectType, len(keys))
			for i := range keys {
				oldKey := keys[i]
				oldKeyIface := oldKey.Interface()
				oldKeyPtr := unsafeops.DataOf(&oldKeyIface)
				oldKeyType := unsafeops.TypeID(&oldKeyIface)
				newKeyPtr := unsafeops.NewObject(v.KeyUnsafeType)
				exec(newKeyPtr, oldKeyPtr, v.KeySubInstructions, false)
				newKeyIface := unsafeops.MakeEface(newKeyPtr, oldKeyType)
				newKeyValue := reflect.ValueOf(newKeyIface)

				oldValue := srcMapReflectValue.MapIndex(oldKey)
				oldValueIface := oldValue.Interface()
				oldValuePtr := unsafeops.DataOf(&oldValueIface)
				oldValueType := unsafeops.TypeID(&oldValueIface)
				newValuePtr := unsafeops.NewObject(v.ValueUnsafeType)
				exec(newValuePtr, oldValuePtr, v.ValueSubInstructions, false)
				newValueIface := unsafeops.MakeEface(newValuePtr, oldValueType)
				newValueValue := reflect.ValueOf(newValueIface)

				newMap.SetMapIndex(newKeyValue, newValueValue)
			}
			newMapPtr := newMap.UnsafePointer()
			*(*unsafe.Pointer)(unsafe.Add(dst, v.Offset)) = newMapPtr
		case *opCopyString:
			s := *(*reflect.StringHeader)(unsafe.Add(src, v.Offset))
			len := s.Len
			data := s.Data
			buffer := make([]byte, len)
			(*reflect.StringHeader)(unsafe.Add(dst, v.Offset)).Data = uintptr(unsafe.Pointer(&buffer[0]))
			(*reflect.StringHeader)(unsafe.Add(dst, v.Offset)).Len = len
			unsafeops.MemMove(
				unsafe.Pointer(&buffer[0]),
				unsafe.Pointer(data),
				uintptr(len),
			)
		case *opCopyStruct:
			newDst := unsafe.Add(dst, v.Offset)
			newSrc := unsafe.Add(src, v.Offset)
			exec(newDst, newSrc, v.SubInstructions, true)
		default:
			// Unreachable
			panic("unreachable")
		}
	}
}
