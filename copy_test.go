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
}
type A struct {
	B uint
	C string
	D []int
	E F
}

var a = A{B: 1, C: "2", D: []int{3, 4, 5}, E: F{g: 6, h: "7", I: new(byte)}}

func BenchmarkCopy(b *testing.B) {
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
