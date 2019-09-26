package contains

import (
	"sync"
)

// SafeSet ...
type SafeSet struct {
	set  []interface{}
	m    map[interface{}]struct{}
	lock *sync.Locker
}
