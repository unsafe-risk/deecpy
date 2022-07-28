package deecpy

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_build(t *testing.T) {
	a_type := reflect.TypeOf(a)
	ops, err := build(a_type)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ops)
}
