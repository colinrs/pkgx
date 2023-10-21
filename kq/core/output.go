package core

import (
	"context"
)

type OutputMessage struct {
	Ctx  context.Context
	Data interface{}
}

type Output interface {
	SendOutput(ctx context.Context, message *OutputMessage) error
	OnDone(cxt context.Context, message *OutputMessage)
	OnError(cxt context.Context, message *OutputMessage, err error)
	Close(ctx context.Context) error
}
