package utils

import (
	"os"
)

// Home returns the os-specific home path as specified in the environment.
func Home() string {
	return os.Getenv("HOME")
}

// SetHome sets the os-specific home path in the environment.
func SetHome(s string) error {
	return os.Setenv("HOME", s)
}
