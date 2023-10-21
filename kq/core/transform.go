package core

import "time"

type Message interface {
	ID() string
	Timestamp() time.Time
}

type Transformer interface {
	Process(Message) (Message, error)
	OnDone(message Message, resp interface{})
	OnError(message Message, err error)
}
