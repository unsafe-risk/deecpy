package deecpy

import (
	"reflect"
	"sync"
)

var isValueTypeCache sync.Map

var ConfigCopyString = false

func isValueType(t reflect.Type) bool {
	// Fast path
	ivC, ok := isValueTypeCache.Load(t)
	if ok {
		return ivC.(bool)
	}
	// Slow path
	switch t.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Bool,
		reflect.Uintptr, reflect.UnsafePointer: // DO NOT DUPLICATE
		isValueTypeCache.Store(t, true)
		return true
	case reflect.String: // String: immutable
		return !ConfigCopyString
	case reflect.Array:
		isValueTypeCache.Store(t, isValueType(t.Elem()))
		return isValueType(t.Elem())
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			if !isValueType(t.Field(i).Type) {
				isValueTypeCache.Store(t, false)
				return false
			}
		}
		isValueTypeCache.Store(t, true)
		return true
	}
	isValueTypeCache.Store(t, false)
	return false
}
