package fx

import (
	"fmt"

	"github.com/colinrs/pkgx/utils"
)

// GoSafe runs the given fn using another goroutine, recovers if fn panics.
func GoSafe(fn func()) {
	go runSafe(fn)
}

// runSafe runs the given fn, recovers if fn panics.
func runSafe(fn func()) {
	defer recoverFunc()
	fn()
}

// runSafeWithRecover runs the given fn, recovers if fn panics.
func runSafeWithRecover(fn func(), recover ...func()) {
	defer recoverFunc(recover...)
	fn()
}

func recoverFunc(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		fmt.Printf("recove:%+v, stack:%s\n", p, utils.Stack())
	}
}

// GoSafeWithRecover runs the given fn, recovers if fn panics.
func GoSafeWithRecover(fn func(), recover ...func()) {
	go runSafeWithRecover(fn, recover...)
}
