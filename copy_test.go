package deecpy

import (
	"math/big"
	"reflect"
	"testing"
)

type F struct {
	g uint64
	h string
	I *byte
	J *struct{}

	Interface interface{}

	EmptySlice []struct{}
}
type A struct {
	B       uint
	C       string
	D       []int
	E       F
	G       [4]F
	H       []F
	I       *A
	J_IFace interface{}
	K_IFace interface{}
}

var a = A{
	B: 1,
	C: "2",
	D: []int{3, 4, 5},
	E: F{g: 6, h: "7", I: new(byte), Interface: new(byte), J: new(struct{}), EmptySlice: []struct{}{}},
	G: [4]F{
		{g: 8, h: "9"},
		{g: 10, h: "11"},
		{g: 12, h: "13"},
		{g: 14, h: "15"},
	},
	H: []F{
		{g: 16, h: "17"},
		{g: 18, h: "19"},
		{g: 20, h: "21"},
		{g: 22, h: "23"},
	},
}

var _ = func() int {
	a.I = &a
	a.J_IFace = &a
	a.K_IFace = &a
	return 0
}()

func BenchmarkCopy(b *testing.B) {
	//nolint:errcheck
	Copy(&A{}, &a) // warmup
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			var b A
			err := Copy(&b, &a)
			if err != nil {
				panic(err)
			}
		}
	})
}

func TestDuplicate(t *testing.T) {
	a := &a
	b, err := Duplicate(a)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(b, a) {
		t.Errorf("expected %v, got %v", a, b)
	}
	if b.I != b {
		t.Errorf("recursive copy failed")
	}
}

func TestCopyBigInt(t *testing.T) {
	var a = big.NewInt(1)
	var b big.Int
	err := Copy(&b, a)
	if err != nil {
		t.Error(err)
	}
	if b.Cmp(a) != 0 {
		t.Errorf("expected %v, got %v", a, b)
	}
	if !reflect.DeepEqual(&b, a) {
		t.Errorf("expected %v, got %v", a, b)
	}
}

func TestDuplicateBigInt(t *testing.T) {
	var a = big.NewInt(1)
	var b, err = Duplicate(a)
	if err != nil {
		t.Error(err)
	}
	if b.Cmp(a) != 0 {
		t.Errorf("expected %v, got %v", a, b)
	}
	if !reflect.DeepEqual(b, a) {
		t.Errorf("expected %v, got %v", a, b)
	}
}

func TestDuplicateRecursive(t *testing.T) {
	type AB struct {
		F *AB
	}
	var a = &AB{F: &AB{}}
	var b, err = Duplicate(a)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(b, a) {
		t.Errorf("expected %v, got %v", a, b)
	}
}

func TestOpPtrDupMem(t *testing.T) {
	type Bravo struct {
		U64  uint64
		U32  uint32
		U16  uint16
		U8   uint8
		I64  int64
		I32  int32
		I16  int16
		I8   int8
		F64  float64
		F32  float32
		C64  complex64
		C128 complex128
		B    bool
	}
	type Alfa struct {
		Charlie *Bravo
	}
	var a = &Alfa{Charlie: &Bravo{
		U64:  1,
		U32:  2,
		U16:  3,
		U8:   4,
		I64:  5,
		I32:  6,
		I16:  7,
		I8:   8,
		F64:  3.14,
		F32:  3.141592,
		C64:  complex(1, 2),
		C128: complex(3.14, 2.718),
		B:    true,
	}}
	var b, err = Duplicate(a)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(b, a) {
		t.Errorf("expected %v, got %v", a, b)
	}
}

func TestCopyMapStringUInt64(t *testing.T) {
	mapA := &map[string]uint64{
		"foo": 1,
		"bar": 2,
	}
	mapB, err := Duplicate(mapA)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(mapB, mapA) {
		t.Errorf("expected %v, got %v", mapA, mapB)
	}
}
