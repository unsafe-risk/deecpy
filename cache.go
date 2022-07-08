package deecpy

import "sync"

var cache sync.Map

func getOps(typ uintptr) (*instructions, bool) {
	v, ok := cache.Load(typ)
	if !ok {
		return nil, false
	}
	return v.(*instructions), ok
}

func setOps(typ uintptr, ops *instructions) {
	cache.Store(typ, ops)
}
