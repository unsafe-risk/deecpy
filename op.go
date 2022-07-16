package deecpy

import (
	"fmt"
	"reflect"
)

type op interface {
	String() string
	Op()
}

type instructions struct {
	ops []op
}

type opCopyMem struct {
	Offset uintptr
	Size   uintptr
}

func (o *opCopyMem) String() string {
	return fmt.Sprintf("copymem(offset: %x, size: %x)", o.Offset, o.Size)
}

func (o *opCopyMem) Op() {}

type opPtrDup struct {
	Offset          uintptr
	Size            uintptr
	UnsafeType      uintptr
	SubInstructions *instructions
}

func (o *opPtrDup) String() string {
	return fmt.Sprintf("ptrdup(offset: %x, size: %x)", o.Offset, o.Size)
}

func (o *opPtrDup) Op() {}

type opPtrDupMem struct {
	Offset     uintptr
	Size       uintptr
	UnsafeType uintptr
}

func (o *opPtrDupMem) String() string {
	return fmt.Sprintf("ptrdupmem(offset: %x, size: %x)", o.Offset, o.Size)
}

func (o *opPtrDupMem) Op() {}

type opArrayCopy struct {
	Offset          uintptr
	ArrayLen        uintptr
	ElemSize        uintptr
	UnsafeElemType  uintptr
	SubInstructions *instructions
}

func (o *opArrayCopy) String() string {
	return fmt.Sprintf("arraycopy(offset: %x, arraylen: %x, elemsize: %x)", o.Offset, o.ArrayLen, o.ElemSize)
}

func (o *opArrayCopy) Op() {}

type opSliceCopy struct {
	Offset          uintptr
	ElemSize        uintptr
	UnsafeElemType  uintptr
	SubInstructions *instructions
}

func (o *opSliceCopy) String() string {
	return fmt.Sprintf("slicecopy(offset: %x, elemsize: %x)", o.Offset, o.ElemSize)
}

func (o *opSliceCopy) Op() {}

type opSliceCopyMem struct {
	Offset         uintptr
	ElemSize       uintptr
	UnsafeElemType uintptr
}

func (o *opSliceCopyMem) String() string {
	return fmt.Sprintf("slicecopymem(offset: %x, elemsize: %x)", o.Offset, o.ElemSize)
}

func (o *opSliceCopyMem) Op() {}

type opMapDup struct {
	Offset uintptr

	ReflectType          reflect.Type
	MapUnsafeType        uintptr
	KeySize              uintptr
	KeyUnsafeType        uintptr
	KeySubInstructions   *instructions
	ValueSize            uintptr
	ValueUnsafeType      uintptr
	ValueSubInstructions *instructions
}

func (o *opMapDup) String() string {
	return fmt.Sprintf("mapdup(offset: %x, reflecttype: %s, mapunsafetype: %x, keysize: %x, valuesize: %x)", o.Offset, o.ReflectType, o.MapUnsafeType, o.KeySize, o.ValueSize)
}

func (o *opMapDup) Op() {}

type opCopyString struct {
	Offset uintptr
}

func (o *opCopyString) String() string {
	return fmt.Sprintf("copystring(offset: %x)", o.Offset)
}

func (o *opCopyString) Op() {}

type opCopyStruct struct {
	Offset          uintptr
	Size            uintptr
	SubInstructions *instructions
}

func (o *opCopyStruct) String() string {
	return fmt.Sprintf("copystruct(offset: %x, size: %x)", o.Offset, o.Size)
}

func (o *opCopyStruct) Op() {}

var _ op = &opCopyMem{}
var _ op = &opPtrDup{}
var _ op = &opPtrDupMem{}
var _ op = &opArrayCopy{}
var _ op = &opSliceCopy{}
var _ op = &opSliceCopyMem{}
var _ op = &opMapDup{}
var _ op = &opCopyString{}
