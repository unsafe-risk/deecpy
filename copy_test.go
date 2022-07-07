package deecpy

import (
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
