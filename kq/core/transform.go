package core

import (
	"context"
	"time"
)

type Message interface {
	ID() string
	Timestamp() time.Time
	Ctx() context.Context
}

type Transformer interface {
	Process(ctx context.Context, message Message) (*OutputMessage, error)
	OnDone(cxt context.Context, message Message)
	OnError(cxt context.Context, message Message, err error)
}
