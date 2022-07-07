package deecpy

import (
	"fmt"
	"reflect"
)

type op interface {
	String() string
	Op()
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
	SubInstructions []op
}

func (o *opPtrDup) String() string {
	return fmt.Sprintf("ptrdup(offset: %x, size: %x, subinstructions: %v)", o.Offset, o.Size, o.SubInstructions)
}

func (o *opPtrDup) Op() {}

type opPtrDupMem struct {
	Offset uintptr
	Size   uintptr
}

func (o *opPtrDupMem) String() string {
	return fmt.Sprintf("ptrdupmem(offset: %x, size: %x)", o.Offset, o.Size)
}

func (o *opPtrDupMem) Op() {}

type opArrayCopy struct {
	Offset          uintptr
	ArrayLen        uintptr
	ElemSize        uintptr
	SubInstructions []op
}

func (o *opArrayCopy) String() string {
	return fmt.Sprintf("arraycopy(offset: %x, arraylen: %x, elemsize: %x, subinstructions: %v)", o.Offset, o.ArrayLen, o.ElemSize, o.SubInstructions)
}

func (o *opArrayCopy) Op() {}

type opSliceCopy struct {
	Offset          uintptr
	ElemSize        uintptr
	SubInstructions []op
}

func (o *opSliceCopy) String() string {
	return fmt.Sprintf("slicecopy(offset: %x, elemsize: %x, subinstructions: %v)", o.Offset, o.ElemSize, o.SubInstructions)
}

func (o *opSliceCopy) Op() {}

type opSliceCopyMem struct {
	Offset   uintptr
	ElemSize uintptr
}

func (o *opSliceCopyMem) String() string {
	return fmt.Sprintf("slicecopymem(offset: %x, elemsize: %x)", o.Offset, o.ElemSize)
}

func (o *opSliceCopyMem) Op() {}

type opMapDup struct {
	Offset uintptr

	ReflectType          reflect.Type
	KeySize              uintptr
	KeySubInstructions   []op
	ValueSize            uintptr
	ValueSubInstructions []op
}

func (o *opMapDup) String() string {
	return fmt.Sprintf("mapdup(offset: %x, keysize: %x, valuesize: %x, keysubinstructions: %v, valuesubinstructions: %v)", o.Offset, o.KeySize, o.ValueSize, o.KeySubInstructions, o.ValueSubInstructions)
}

func (o *opMapDup) Op() {}

type opCopyString struct {
	Offset uintptr
}

func (o *opCopyString) String() string {
	return fmt.Sprintf("copystring(offset: %x)", o.Offset)
}

func (o *opCopyString) Op() {}

var _ op = &opCopyMem{}
var _ op = &opPtrDup{}
var _ op = &opPtrDupMem{}
var _ op = &opArrayCopy{}
var _ op = &opSliceCopy{}
var _ op = &opSliceCopyMem{}
var _ op = &opMapDup{}
var _ op = &opCopyString{}
