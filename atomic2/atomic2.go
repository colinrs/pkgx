package atomic2

import (
	"sync"
	"sync/atomic"
)

// AtomicInt32 ...
type AtomicInt32 struct {
	int32
}

// NewAtomicInt32 ...
func NewAtomicInt32(n int32) AtomicInt32 {
	return AtomicInt32{n}
}

// Add ...
func (i *AtomicInt32) Add(n int32) int32 {
	return atomic.AddInt32(&i.int32, n)
}

// Set ...
func (i *AtomicInt32) Set(n int32) {
	atomic.StoreInt32(&i.int32, n)
}

// Get ...
func (i *AtomicInt32) Get() int32 {
	return atomic.LoadInt32(&i.int32)
}

// CompareAndSwap ...
func (i *AtomicInt32) CompareAndSwap(oldval, newval int32) (swapped bool) {
	return atomic.CompareAndSwapInt32(&i.int32, oldval, newval)
}

// AtomicInt64 ...
type AtomicInt64 struct {
	int64
}

// NewAtomicInt64 ...
func NewAtomicInt64(n int64) AtomicInt64 {
	return AtomicInt64{n}
}

// Add ...
func (i *AtomicInt64) Add(n int64) int64 {
	return atomic.AddInt64(&i.int64, n)
}

// Set ...
func (i *AtomicInt64) Set(n int64) {
	atomic.StoreInt64(&i.int64, n)
}

// Get ...
func (i *AtomicInt64) Get() int64 {
	return atomic.LoadInt64(&i.int64)
}

// CompareAndSwap ...
func (i *AtomicInt64) CompareAndSwap(oldval, newval int64) (swapped bool) {
	return atomic.CompareAndSwapInt64(&i.int64, oldval, newval)
}

// AtomicBool ...
type AtomicBool struct {
	int32
}

// NewAtomicBool ...
func NewAtomicBool(n bool) AtomicBool {
	if n {
		return AtomicBool{1}
	}
	return AtomicBool{0}
}

// Set ...
func (i *AtomicBool) Set(n bool) {
	if n {
		atomic.StoreInt32(&i.int32, 1)
	} else {
		atomic.StoreInt32(&i.int32, 0)
	}
}

// Get ...
func (i *AtomicBool) Get() bool {
	return atomic.LoadInt32(&i.int32) != 0
}

// CompareAndSwap ...
func (i *AtomicBool) CompareAndSwap(o, n bool) bool {
	var old, new int32
	if o {
		old = 1
	}
	if n {
		new = 1
	}
	return atomic.CompareAndSwapInt32(&i.int32, old, new)
}

// AtomicString ...
type AtomicString struct {
	mu  sync.Mutex
	str string
}

// Set ...
func (s *AtomicString) Set(str string) {
	s.mu.Lock()
	s.str = str
	s.mu.Unlock()
}

// Get ...
func (s *AtomicString) Get() string {
	s.mu.Lock()
	str := s.str
	s.mu.Unlock()
	return str
}

// CompareAndSwap ...
func (s *AtomicString) CompareAndSwap(oldval, newval string) (swqpped bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.str == oldval {
		s.str = newval
		return true
	}
	return false
}
