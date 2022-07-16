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
	inst, ok := getOps(typID)
	if !ok {
		var err error
		inst, err = build(reflect.TypeOf(src).Elem())
		if err != nil {
			return err
		}
		setOps(typID, inst)
	}

	exec(
		unsafe.Pointer(dst),
		unsafe.Pointer(src),
		inst,
		false,
	)
	return nil
}

func Duplicate[T any](src T) (T, error) {
	var dst T
	err := Copy(unsafeops.NoEscape(&dst), unsafeops.NoEscape(&src))
	if err != nil {
		return dst, err
	}
	return dst, nil
}
