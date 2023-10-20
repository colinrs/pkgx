package os

import (
	"fmt"
	"os"
	"runtime"
)

var gomaxprocs = runtime.GOMAXPROCS
var numCPU = runtime.NumCPU

// Home returns the os-specific home path as specified in the environment.
func Home() string {
	return os.Getenv("HOME")
}

// SetHome sets the os-specific home path in the environment.
func SetHome(s string) error {
	return os.Setenv("HOME", s)
}

// UseMultipleCPUs sets GOMAXPROCS to the number of CPU cores unless it has
// already been overridden by the GOMAXPROCS environment variable.
func UseMultipleCPUs() {
	if envGOMAXPROCS := os.Getenv("GOMAXPROCS"); envGOMAXPROCS != "" {
		n := gomaxprocs(0)
		fmt.Printf("GOMAXPROCS already set in environment to %q, %d internally\n",
			envGOMAXPROCS, n)
		return
	}
	n := numCPU()
	fmt.Printf("setting GOMAXPROCS to %d\n", n)
	gomaxprocs(n)
}
