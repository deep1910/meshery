// Package store provides methods for interacting
// with a thread safe global store.
package store

import "strings"

// globalStore is intentionally a global variable so
// that the life time of the variable is through
// out the program lifecycle
var globalStore *threadSafeStore

// Initialize initializes the global store
//
// It will skip initialization if the
// store has already been initialized once
func Initialize() {
	if globalStore == nil {
		globalStore = newThreadSafeStore()
	}
}

// Set will set the key value pairs in the global store
//
// If the key already exists then it will OVERWRITE
// those values
func Set(key string, value interface{}) {
	globalStore.Lock()
	defer globalStore.Unlock()

	if globalStore.store[key] == nil {
		globalStore.store[key] = []interface{}{}
	}

	globalStore.store[key] = append(globalStore.store[key], value)
}

// Get will get the value corresponding to the given key
func Get(key string) (interface{}, bool) {
	globalStore.RLock()
	defer globalStore.RUnlock()

	val, ok := globalStore.store[key]
	if ok && len(val) > 0 {
		return val[0], true
	}

	return nil, false
}

// GetAll returns all the values stored against the key
func GetAll(key string) []interface{} {
	globalStore.RLock()
	defer globalStore.RUnlock()

	val, ok := globalStore.store[key]
	if !ok {
		return []interface{}{}
	}

	return val
}

// PrefixMatch will return all the values which matches the given key
func PrefixMatch(key string) (res []interface{}) {
	globalStore.RLock()
	defer globalStore.RUnlock()

	for k, v := range globalStore.store {
		if strings.HasPrefix(k, key) {
			res = append(res, v)
		}
	}

	return
}
