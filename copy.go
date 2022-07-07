package deecpy

import (
	"reflect"
	"unsafe"

	"github.com/unsafe-risk/deecpy/unsafeops"
)

func Copy[T any](dst, src *T) error {
	sAny := any(src)
	typID := unsafeops.TypeID(&sAny)

	// Lookup the type in the cache
	ops, ok := getOps(typID)
	if !ok {
		var err error
		ops, err = build(reflect.TypeOf(src).Elem())
		if err != nil {
			return err
		}
	}

	exec(
		unsafe.Pointer(dst),
		unsafe.Pointer(src),
		ops,
	)
	return nil
}

func Duplicate[T any](src *T) (*T, error) {
	var dst T
	err := Copy(&dst, src)
	if err != nil {
		return nil, err
	}
	return &dst, nil
}
