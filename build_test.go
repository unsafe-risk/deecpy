package deecpy

import (
	"io"
	"reflect"
	"testing"
)

func Test_build(t *testing.T) {
	a_type := reflect.TypeOf(a)
	debugBuild(a_type, io.Discard)
}
