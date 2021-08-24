package ratelimit

import (
	"errors"
)

var (
	ErrExceededLimit = errors.New("Too many requests, exceeded the limit. ")
)

// RateLimiter ...
type RateLimiter interface {
	Take() error
}
