package deecpy

import (
	"reflect"
)

type op interface {
	Op()
}

type instructions struct {
	ops []op
}

type opCopyMem struct {
	Offset uintptr
	Size   uintptr
}

func (o *opCopyMem) Op() {}

type opPtrDup struct {
	Offset          uintptr
	Size            uintptr
	UnsafeType      uintptr
	SubInstructions *instructions
}

func (o *opPtrDup) Op() {}

type opPtrDupMem struct {
	Offset     uintptr
	Size       uintptr
	UnsafeType uintptr
}

func (o *opPtrDupMem) Op() {}

type opArrayCopy struct {
	Offset          uintptr
	ArrayLen        uintptr
	ElemSize        uintptr
	UnsafeElemType  uintptr
	SubInstructions *instructions
}

func (o *opArrayCopy) Op() {}

type opSliceCopy struct {
	Offset          uintptr
	ElemSize        uintptr
	UnsafeElemType  uintptr
	SubInstructions *instructions
}

func (o *opSliceCopy) Op() {}

type opSliceCopyMem struct {
	Offset         uintptr
	ElemSize       uintptr
	UnsafeElemType uintptr
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

func (o *opMapDup) Op() {}

type opCopyString struct {
	Offset uintptr
}

func (o *opCopyString) Op() {}

type opCopyStruct struct {
	Offset          uintptr
	Size            uintptr
	SubInstructions *instructions
}

func (o *opCopyStruct) Op() {}

type opCopyInterface struct {
	Offset uintptr
}

func (o *opCopyInterface) Op() {}

var _ op = (*opCopyMem)(nil)
var _ op = (*opPtrDup)(nil)
var _ op = (*opPtrDupMem)(nil)
var _ op = (*opArrayCopy)(nil)
var _ op = (*opSliceCopy)(nil)
var _ op = (*opSliceCopyMem)(nil)
var _ op = (*opMapDup)(nil)
var _ op = (*opCopyString)(nil)
var _ op = (*opCopyStruct)(nil)
var _ op = (*opCopyInterface)(nil)
