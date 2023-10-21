package core

import (
	"context"
	"sync"
)

type OutputMessage struct {
	data interface{}
	done sync.Mutex
}

type Output interface {
	SendOutput(ctx context.Context, msg *OutputMessage) error
	Close(ctx context.Context) error
}
