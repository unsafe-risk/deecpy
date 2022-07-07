package deecpy

import "sync"

type cacheEntry struct {
	typ uintptr
	ops []op
}

const cacheSize = 1 << 12

var cache = make([][]cacheEntry, cacheSize)
var mu sync.RWMutex

func getOps(typ uintptr) ([]op, bool) {
	mu.RLock()
	defer mu.RUnlock()

	hash := typ & cacheSize

	if cache[hash] == nil {
		return nil, false
	}

	for i := range cache[hash] {
		if cache[hash][i].typ == typ {
			return cache[hash][i].ops, true
		}
	}

	return nil, false
}

func setOps(typ uintptr, ops []op) {
	mu.Lock()
	defer mu.Unlock()

	hash := typ & 4096

	if cache[hash] == nil {
		cache[hash] = make([]cacheEntry, 8)
	}

	for i := range cache[hash] {
		if cache[hash][i].typ == typ {
			cache[hash][i].ops = ops
			return
		}
	}

	cache[hash] = append(cache[hash], cacheEntry{typ, ops})
}
