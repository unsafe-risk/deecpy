package deecpy

import (
	"reflect"
	"unsafe"

	"github.com/unsafe-risk/deecpy/unsafeops"
)

func exec(dst, src unsafe.Pointer, inst *instructions) {
L:
	for i := range inst.ops {
		switch v := inst.ops[i].(type) {
		case *opCopyMem:
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
			obj := make([]byte, v.Size)
			*(*unsafe.Pointer)(unsafe.Add(dst, v.Offset)) = unsafe.Pointer(&obj[0])
			exec(*(*unsafe.Pointer)(unsafe.Add(dst, v.Offset)), srcPtr, v.SubInstructions)
		case *opPtrDupMem:
			srcPtr := *(*unsafe.Pointer)(unsafe.Add(src, v.Offset))
			if srcPtr == nil {
				*(*unsafe.Pointer)(unsafe.Add(dst, v.Offset)) = nil
				continue L
			}
			obj := make([]byte, v.Size)
			*(*unsafe.Pointer)(unsafe.Add(dst, v.Offset)) = unsafe.Pointer(&obj[0])
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
				)
			}
		case *opSliceCopy:
			s := *(*reflect.SliceHeader)(unsafe.Add(src, v.Offset))
			(*reflect.SliceHeader)(unsafe.Add(dst, v.Offset)).Cap = s.Cap
			(*reflect.SliceHeader)(unsafe.Add(dst, v.Offset)).Len = s.Len
			if s.Cap == 0 {
				continue L
			}
			sliceBuffer := make([]byte, uintptr(s.Cap)*v.ElemSize)
			(*reflect.SliceHeader)(unsafe.Add(dst, v.Offset)).Data = uintptr(unsafe.Pointer(&sliceBuffer[0]))
			for i := uintptr(0); i < uintptr(s.Cap); i++ {
				exec(
					unsafe.Add(dst, v.Offset+uintptr(i)*v.ElemSize),
					unsafe.Add(src, v.Offset+uintptr(i)*v.ElemSize),
					v.SubInstructions,
				)
			}
		case *opSliceCopyMem:
			s := *(*reflect.SliceHeader)(unsafe.Add(src, v.Offset))
			(*reflect.SliceHeader)(unsafe.Add(dst, v.Offset)).Cap = s.Cap
			(*reflect.SliceHeader)(unsafe.Add(dst, v.Offset)).Len = s.Len
			if s.Cap == 0 {
				continue L
			}
			sliceBuffer := make([]byte, uintptr(s.Cap)*v.ElemSize)
			(*reflect.SliceHeader)(unsafe.Add(dst, v.Offset)).Data = uintptr(unsafe.Pointer(&sliceBuffer[0]))
			unsafeops.MemMove(
				unsafe.Pointer(&sliceBuffer[0]),
				unsafe.Pointer(s.Data),
				uintptr(s.Cap)*v.ElemSize,
			)
		case *opMapDup:
			// TODO: implement
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
			exec(newDst, newSrc, v.SubInstructions)
		default:
			// Unreachable
			panic("unreachable")
		}
	}
}
