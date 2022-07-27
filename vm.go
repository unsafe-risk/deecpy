package deecpy

import (
	"reflect"
	"unsafe"

	"github.com/unsafe-risk/deecpy/unsafeops"
)

type ptrmap struct {
	k unsafe.Pointer
	v unsafe.Pointer
}

func ptrmapSearch(m *[]ptrmap, k unsafe.Pointer) (unsafe.Pointer, bool) {
	for _, v := range *m {
		if v.k == k {
			return v.v, true
		}
	}
	return nil, false
}

func exec(dst, src unsafe.Pointer, inst *instructions, nocopy bool, pmap *[]ptrmap) {
	*pmap = append(*pmap, ptrmap{k: src, v: dst})
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
			if val, ok := ptrmapSearch(pmap, srcPtr); ok {
				*(*unsafe.Pointer)(unsafe.Add(dst, v.Offset)) = val
				continue
			}
			*(*unsafe.Pointer)(unsafe.Add(dst, v.Offset)) = unsafeops.NewObject(v.UnsafeType)
			exec(*(*unsafe.Pointer)(unsafe.Add(dst, v.Offset)), srcPtr, v.SubInstructions, false, pmap)
		case *opPtrDupMem:
			srcPtr := *(*unsafe.Pointer)(unsafe.Add(src, v.Offset))
			if srcPtr == nil {
				*(*unsafe.Pointer)(unsafe.Add(dst, v.Offset)) = nil
				continue L
			}
			if val, ok := ptrmapSearch(pmap, srcPtr); ok {
				*(*unsafe.Pointer)(unsafe.Add(dst, v.Offset)) = val
				continue
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
					pmap,
				)
			}
		case *opSliceCopy:
			s := (*reflect.SliceHeader)(unsafe.Add(src, v.Offset))
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
					pmap,
				)
			}
		case *opSliceCopyMem:
			s := (*reflect.SliceHeader)(unsafe.Add(src, v.Offset))
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
				var newKeyPtr unsafe.Pointer
				if val, ok := ptrmapSearch(pmap, oldKeyPtr); ok {
					newKeyPtr = val
				} else {
					newKeyPtr = unsafeops.NewObject(v.KeyUnsafeType)
					exec(newKeyPtr, oldKeyPtr, v.KeySubInstructions, false, pmap)
				}
				newKeyIface := unsafeops.MakeEface(newKeyPtr, oldKeyType)
				newKeyValue := reflect.ValueOf(newKeyIface)

				oldValue := srcMapReflectValue.MapIndex(oldKey)
				oldValueIface := oldValue.Interface()
				oldValuePtr := unsafeops.DataOf(&oldValueIface)
				oldValueType := unsafeops.TypeID(&oldValueIface)
				var newValuePtr unsafe.Pointer
				if val, ok := ptrmapSearch(pmap, oldValuePtr); ok {
					newValuePtr = val
				} else {
					newValuePtr = unsafeops.NewObject(v.ValueUnsafeType)
					exec(newValuePtr, oldValuePtr, v.ValueSubInstructions, false, pmap)
				}
				newValueIface := unsafeops.MakeEface(newValuePtr, oldValueType)
				newValueValue := reflect.ValueOf(newValueIface)

				newMap.SetMapIndex(newKeyValue, newValueValue)
			}
			newMapPtr := newMap.UnsafePointer()
			*(*unsafe.Pointer)(unsafe.Add(dst, v.Offset)) = newMapPtr
		case *opCopyString:
			s := (*reflect.StringHeader)(unsafe.Add(src, v.Offset))
			buffer := make([]byte, s.Len)
			(*reflect.StringHeader)(unsafe.Add(dst, v.Offset)).Data = uintptr(unsafe.Pointer(&buffer[0]))
			(*reflect.StringHeader)(unsafe.Add(dst, v.Offset)).Len = s.Len
			unsafeops.MemMove(
				unsafe.Pointer(&buffer[0]),
				unsafe.Pointer(s.Data),
				uintptr(s.Len),
			)
		case *opCopyStruct:
			newDst := unsafe.Add(dst, v.Offset)
			newSrc := unsafe.Add(src, v.Offset)
			exec(newDst, newSrc, v.SubInstructions, true, pmap)
		default:
			// Unreachable
			panic("unreachable")
		}
	}
}
