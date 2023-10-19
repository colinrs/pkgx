package structx

import (
	"sync"
)

// SafeSet ...
type SafeSet struct {
	sync.RWMutex
	M map[string]bool
}

// NewSafeSet ...
func NewSafeSet() *SafeSet {
	return &SafeSet{
		M: make(map[string]bool),
	}
}

// Add ...
func (safeSet *SafeSet) Add(key string) {
	safeSet.Lock()
	safeSet.M[key] = true
	safeSet.Unlock()
}

// Remove ...
func (safeSet *SafeSet) Remove(key string) {
	safeSet.Lock()
	delete(safeSet.M, key)
	safeSet.Unlock()
}

// Clear ...
func (safeSet *SafeSet) Clear() {
	safeSet.Lock()
	safeSet.M = make(map[string]bool)
	safeSet.Unlock()
}

// Contains ...
func (safeSet *SafeSet) Contains(key string) bool {
	safeSet.RLock()
	_, exists := safeSet.M[key]
	safeSet.RUnlock()
	return exists
}

// Size ...
func (safeSet *SafeSet) Size() int {
	safeSet.RLock()
	len := len(safeSet.M)
	safeSet.RUnlock()
	return len
}

// ToSlice ...
func (safeSet *SafeSet) ToSlice() []string {
	safeSet.RLock()
	defer safeSet.RUnlock()

	count := len(safeSet.M)
	if count == 0 {
		return []string{}
	}

	r := []string{}
	for key := range safeSet.M {
		r = append(r, key)
	}

	return r
}
