package core

import (
	"context"
)

type Extractor interface {
	Unmarshal(cxt context.Context, message *InputMessage) (Message, error)
	OnDone(cxt context.Context, message *InputMessage)
	OnError(cxt context.Context, message *InputMessage, err error)
}
