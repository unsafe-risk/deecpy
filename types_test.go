package deecpy

import (
	"reflect"
	"testing"
)

func Test_isValue(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Byte Array",
			args{reflect.TypeOf([16]byte{})},
			true,
		},
		{
			"Struct Array",
			args{reflect.TypeOf([256]struct {
				A int
				B string
				C float32
			}{})},
			true,
		},
		{
			"Mixed Slice",
			args{reflect.TypeOf([]struct {
				A   int
				Map map[string]int
				B   string
				Ptr *int
				C   float32
			}{})},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValueType(tt.args.t); got != tt.want {
				t.Errorf("isValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
