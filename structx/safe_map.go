package structx

import (
	"sync"
)

// SafeMap ...
type SafeMap struct {
	sync.RWMutex
	M map[string]interface{}
}

// NewSafeMap ...
func NewSafeMap() *SafeMap {
	return &SafeMap{
		M: make(map[string]interface{}),
	}
}

// Put ...
func (safeMap *SafeMap) Put(key string, val interface{}) {
	safeMap.Lock()
	safeMap.M[key] = val
	safeMap.Unlock()
}

// Get ...
func (safeMap *SafeMap) Get(key string) (interface{}, bool) {
	safeMap.RLock()
	val, exists := safeMap.M[key]
	safeMap.RUnlock()
	return val, exists
}

// Remove ...
func (safeMap *SafeMap) Remove(key string) {
	safeMap.Lock()
	delete(safeMap.M, key)
	safeMap.Unlock()
}

// GetAndRemove ...
func (safeMap *SafeMap) GetAndRemove(key string) (interface{}, bool) {
	safeMap.Lock()
	val, exists := safeMap.M[key]
	if exists {
		delete(safeMap.M, key)
	}
	safeMap.Unlock()
	return val, exists
}

// Clear ...
func (safeMap *SafeMap) Clear() {
	safeMap.Lock()
	safeMap.M = make(map[string]interface{})
	safeMap.Unlock()
}

// Keys ...
func (safeMap *SafeMap) Keys() []string {
	safeMap.RLock()
	defer safeMap.RUnlock()

	keys := make([]string, 0)
	for key := range safeMap.M {
		keys = append(keys, key)
	}
	return keys
}

// Slice ...
func (safeMap *SafeMap) Slice() []interface{} {
	safeMap.RLock()
	defer safeMap.RUnlock()

	vals := make([]interface{}, 0)
	for _, val := range safeMap.M {
		vals = append(vals, val)
	}
	return vals
}

// ContainsKey ...
func (safeMap *SafeMap) ContainsKey(key string) bool {
	safeMap.RLock()
	_, exists := safeMap.M[key]
	safeMap.RUnlock()
	return exists
}

// Size ...
func (safeMap *SafeMap) Size() int {
	safeMap.RLock()
	len := len(safeMap.M)
	safeMap.RUnlock()
	return len
}

// IsEmpty ...
func (safeMap *SafeMap) IsEmpty() bool {
	safeMap.RLock()
	empty := (len(safeMap.M) == 0)
	safeMap.RUnlock()
	return empty
}
