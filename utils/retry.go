package utils

import (
	"time"
)

// Retry ...
func Retry(fc func() error, maxRetries int, interval time.Duration) error {
	var err error
	for i := 1; i <= maxRetries; i++ {
		if err = fc(); err != nil {
			time.Sleep(interval)
			continue
		}
		return nil
	}
	return err
}
