package fx

import (
	"sync"

	"github.com/colinrs/pkgx/concurrent"
)

// A RoutineGroup is used to group goroutines together and all wait all goroutines to be done.
type RoutineGroup struct {
	waitGroup sync.WaitGroup
	Limit     *concurrent.Limit
}

// NewRoutineGroup returns a RoutineGroup.
func NewRoutineGroup() *RoutineGroup {
	return &RoutineGroup{
		waitGroup: sync.WaitGroup{},
		Limit:     concurrent.NewLimit(defaultGoConcurrentLimit),
	}
}

// Run runs the given fn in RoutineGroup.
// Don't reference the variables from outside,
// because outside variables can be changed by other goroutines
func (g *RoutineGroup) Run(fn func()) {
	g.waitGroup.Add(1)
	g.Limit.Acquire()

	go func() {
		defer g.waitGroup.Done()
		defer g.Limit.Release()
		fn()
	}()
}

// runSafe runs the given fn in RoutineGroup, and avoid panics.
// Don't reference the variables from outside,
// because outside variables can be changed by other goroutines
func (g *RoutineGroup) RunGoSafe(fn func()) {
	g.waitGroup.Add(1)
	g.Limit.Acquire()
	GoSafe(func() {
		defer g.waitGroup.Done()
		defer g.Limit.Release()
		fn()
	})
}

// Wait waits all running functions to be done.
func (g *RoutineGroup) Wait() {
	g.waitGroup.Wait()
}
