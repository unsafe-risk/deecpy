package deecpy

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var ErrUnsupportedType = errors.New("unsupported type")

var IgnoreUnsupportedTypes = true

var buildCache sync.Map

func build(t reflect.Type) ([]op, error) {
	var ops []op

	// Lookup in cache
	if v, ok := buildCache.Load(t); ok {
		ops = v.([]op)
		return ops, nil
	}

	if isValueType(t) {
		ops = append(ops, &opCopyMem{
			Offset: 0,
			Size:   t.Size(),
		})
		buildCache.Store(t, ops)
		return ops, nil
	}

	switch t.Kind() {
	case reflect.Pointer:
		elem := t.Elem()
		elemSize := elem.Size()
		if elemSize == 0 {
			// struct{} Pointer
			return ops, nil
		}
		if isValueType(elem) {
			ops = append(ops,
				&opPtrDupMem{
					Offset: 0,
					Size:   elem.Size(),
				},
			)
			buildCache.Store(t, ops)
			return ops, nil
		}
		subOps, err := build(elem)
		if err != nil {
			return nil, err
		}
		ops = append(ops,
			&opPtrDup{
				Offset:          0,
				Size:            t.Size(),
				SubInstructions: subOps,
			},
		)
		buildCache.Store(t, ops)
		return ops, nil
	case reflect.Array:
		elem := t.Elem()
		elemSize := elem.Size()
		subOps, err := build(elem)
		if err != nil {
			return nil, err
		}
		ops = append(ops, &opArrayCopy{
			Offset:          0,
			ArrayLen:        uintptr(t.Len()),
			ElemSize:        elemSize,
			SubInstructions: subOps,
		})
		buildCache.Store(t, ops)
		return ops, nil
	case reflect.Slice:
		elem := t.Elem()
		elemSize := elem.Size()
		if isValueType(elem) {
			ops = append(ops, &opSliceCopyMem{
				Offset:   0,
				ElemSize: elemSize,
			})
			buildCache.Store(t, ops)
			return ops, nil
		}
		subOps, err := build(elem)
		if err != nil {
			return nil, err
		}
		ops = append(ops, &opSliceCopy{
			Offset:          0,
			ElemSize:        elemSize,
			SubInstructions: subOps,
		})
		buildCache.Store(t, ops)
		return ops, nil
	case reflect.Struct:
		var valueTypes uint
		var nonValueTypes uint
		var structOps []op

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldType := field.Type
			if isValueType(fieldType) {
				valueTypes++
			} else {
				nonValueTypes++
				subOps, err := build(fieldType)
				if err != nil {
					return nil, err
				}
				subOps = append([]op{}, subOps...) // Duplicate subOps

				// Apply offset
				for i := range subOps {
					switch subOps[i].(type) {
					case *opCopyMem:
						subOps[i].(*opCopyMem).Offset += field.Offset
					case *opPtrDup:
						subOps[i].(*opPtrDup).Offset += field.Offset
					case *opPtrDupMem:
						subOps[i].(*opPtrDupMem).Offset += field.Offset
					case *opArrayCopy:
						subOps[i].(*opArrayCopy).Offset += field.Offset
					case *opSliceCopy:
						subOps[i].(*opSliceCopy).Offset += field.Offset
					case *opSliceCopyMem:
						subOps[i].(*opSliceCopyMem).Offset += field.Offset
					case *opMapDup:
						subOps[i].(*opMapDup).Offset += field.Offset
					}
				}
				structOps = append(structOps, subOps...)
			}
		}

		if valueTypes > 0 && nonValueTypes > 0 {
			ops = append(ops, &opCopyMem{
				Offset: 0,
				Size:   t.Size(),
			})
			// Check if optmization is possible
			var hasCopyMem bool
			for i := range structOps {
				if _, ok := structOps[i].(*opCopyMem); ok {
					hasCopyMem = true
					break
				}
			}
			if hasCopyMem {
				// Optimize: remove overlapping CopyMem instructions
				var newStructOps []op = make([]op, 0, len(structOps))
				for i := range structOps {
					if _, ok := structOps[i].(*opCopyMem); ok {
						continue
					}
					newStructOps = append(newStructOps, structOps[i])
				}
				structOps = newStructOps
			}
			ops = append(ops, structOps...)
		} else if valueTypes > 0 && nonValueTypes == 0 {
			// Unreachable code
			panic("unreachable")
		} else if valueTypes == 0 && nonValueTypes > 0 {
			ops = append(ops, structOps...)
		} else if valueTypes == 0 && nonValueTypes == 0 {
			// Unreachable code
			panic("unreachable")
		} else {
			// Unreachable code
			panic("unreachable")
		}

		buildCache.Store(t, ops)
		return ops, nil
	case reflect.Map:
		key := t.Key()
		keySize := key.Size()
		elem := t.Elem()
		elemSize := elem.Size()
		keySubOps, err := build(key)
		if err != nil {
			return nil, err
		}
		elemSubOps, err := build(elem)
		if err != nil {
			return nil, err
		}
		// TODO: Optimize Map Duplication
		ops = append(ops, &opMapDup{
			Offset:               0,
			ReflectType:          t,
			KeySize:              keySize,
			ValueSize:            elemSize,
			KeySubInstructions:   keySubOps,
			ValueSubInstructions: elemSubOps,
		})
		buildCache.Store(t, ops)
		return ops, nil
	case reflect.String:
		if !ConfigCopyString {
			// Unreachable code
			panic("unreachable")
		}
		ops = append(ops, &opCopyString{
			Offset: 0,
		})
	default:
		// Unsupported type
		if IgnoreUnsupportedTypes {
			// Use CopyMem for unsupported types
			ops = append(ops, &opCopyMem{
				Offset: 0,
				Size:   t.Size(),
			})
			buildCache.Store(t, ops)
		}
	}

	if IgnoreUnsupportedTypes {
		return ops, nil
	}
	return ops, ErrUnsupportedType
}

func debugBuild(t reflect.Type) {
	ops, err := build(t)
	if err != nil {
		panic(err)
	}
	for i := range ops {
		fmt.Printf("%v\n", ops[i])
	}
}

//nolint:unused
var _ = debugBuild
