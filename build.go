package deecpy

import (
	"errors"
	"reflect"
	"sync"

	"v8.run/go/deecpy/unsafeops"
)

var ErrUnsupportedType = errors.New("unsupported type")

var IgnoreUnsupportedTypes = true

var buildCache sync.Map

func build(t reflect.Type) (*instructions, error) {
	var inst instructions

	// Lookup in cache
	if v, ok := buildCache.Load(t); ok {
		return v.(*instructions), nil
	}

	// Store Cache
	buildCache.Store(t, &inst)

	if isValueType(t) {
		inst.ops = append(inst.ops, &opCopyMem{
			Offset: 0,
			Size:   t.Size(),
		})
		return &inst, nil
	}

	switch t.Kind() {
	case reflect.Pointer:
		elem := t.Elem()
		elemSize := elem.Size()
		if elemSize == 0 {
			// struct{} Pointer
			return &inst, nil
		}
		if isValueType(elem) {
			inst.ops = append(inst.ops,
				&opPtrDupMem{
					Offset:     0,
					Size:       elemSize,
					UnsafeType: unsafeops.UnsafeType(elem),
				},
			)
			return &inst, nil
		}
		subInsts, err := build(elem)
		if err != nil {
			return nil, err
		}
		inst.ops = append(inst.ops,
			&opPtrDup{
				Offset:          0,
				Size:            elemSize,
				UnsafeType:      unsafeops.UnsafeType(elem),
				SubInstructions: subInsts,
			},
		)
		return &inst, nil
	case reflect.Array:
		elem := t.Elem()
		elemSize := elem.Size()
		subInsts, err := build(elem)
		if err != nil {
			return nil, err
		}
		inst.ops = append(inst.ops, &opArrayCopy{
			Offset:          0,
			ArrayLen:        uintptr(t.Len()),
			ElemSize:        elemSize,
			UnsafeElemType:  unsafeops.UnsafeType(elem),
			SubInstructions: subInsts,
		})
		return &inst, nil
	case reflect.Slice:
		elem := t.Elem()
		elemSize := elem.Size()
		if isValueType(elem) {
			inst.ops = append(inst.ops, &opSliceCopyMem{
				Offset:         0,
				ElemSize:       elemSize,
				UnsafeElemType: unsafeops.UnsafeType(elem),
			})
			return &inst, nil
		}
		subInsts, err := build(elem)
		if err != nil {
			return nil, err
		}
		inst.ops = append(inst.ops, &opSliceCopy{
			Offset:          0,
			ElemSize:        elemSize,
			UnsafeElemType:  unsafeops.UnsafeType(elem),
			SubInstructions: subInsts,
		})
		return &inst, nil
	case reflect.Struct:
		var valueTypes uint
		var nonValueTypes uint
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldType := field.Type
			if isValueType(fieldType) {
				valueTypes++
			} else {
				nonValueTypes++
			}
		}
		if valueTypes > 0 {
			inst.ops = append(inst.ops, &opCopyMem{
				Offset: 0,
				Size:   t.Size(),
			})
		}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldType := field.Type
			if isValueType(fieldType) {
				valueTypes++
			} else {
				nonValueTypes++
				subInsts, err := build(fieldType)
				if err != nil {
					return nil, err
				}
				inst.ops = append(inst.ops, &opCopyStruct{
					Offset:          field.Offset,
					Size:            fieldType.Size(),
					SubInstructions: subInsts,
				})
			}
		}

		return &inst, nil
	case reflect.Map:
		Map := reflect.MakeMap(t)
		MapIface := Map.Interface()
		MapType := unsafeops.TypeID(&MapIface)
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
		inst.ops = append(inst.ops, &opMapDup{
			Offset:               0,
			ReflectType:          t,
			MapUnsafeType:        MapType,
			KeySize:              keySize,
			ValueSize:            elemSize,
			KeyUnsafeType:        unsafeops.UnsafeType(key),
			ValueUnsafeType:      unsafeops.UnsafeType(elem),
			KeySubInstructions:   keySubOps,
			ValueSubInstructions: elemSubOps,
		})
		return &inst, nil
	case reflect.String:
		if !ConfigCopyString {
			// Unreachable code
			panic("unreachable")
		}
		inst.ops = append(inst.ops, &opCopyString{
			Offset: 0,
		})
	case reflect.Interface:
		inst.ops = append(inst.ops, &opCopyInterface{
			Offset: 0,
		})
	default:
		// Unsupported type
		if IgnoreUnsupportedTypes {
			// Use CopyMem for unsupported types
			inst.ops = append(inst.ops, &opCopyMem{
				Offset: 0,
				Size:   t.Size(),
			})
			return &inst, nil
		}
	}

	if IgnoreUnsupportedTypes {
		return &inst, nil
	}
	return nil, ErrUnsupportedType
}
