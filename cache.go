package deecpy

import "sync"

type cacheEntry struct {
	typ uintptr
	ops *instructions
}

const cacheSize = 1 << 12
const cacheIndexMask = cacheSize - 1

var cache = make([][]cacheEntry, cacheSize)
var mu sync.RWMutex

func getOps(typ uintptr) (*instructions, bool) {
	mu.RLock()
	defer mu.RUnlock()

	hash := typ & cacheIndexMask

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

func setOps(typ uintptr, ops *instructions) {
	mu.Lock()
	defer mu.Unlock()

	hash := typ & cacheIndexMask

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
