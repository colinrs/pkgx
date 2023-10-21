package core

import (
	"context"
	"sync"

	"github.com/IBM/sarama"
)

var (
	inputMessagePool = sync.Pool{
		New: func() interface{} {
			return &InputMessage{}
		},
	}
)

func NewInputMessage() *InputMessage {
	return inputMessagePool.Get().(*InputMessage)
}

func (m *InputMessage) Release() {
	inputMessagePool.Put(m)
}

type InputMessage struct {
	Raw  *sarama.ConsumerMessage
	done sync.Mutex
}

// Ack acknowledges that the message has been processed.
// It will release the internal lock on the message, allowing the message to be marked in kafka.
func (m *InputMessage) Ack() {
	m.done.Unlock()
}

func (m *InputMessage) Lock() {
	m.done.Lock()
}

type Input interface {
	Consumer(ctx context.Context) (*InputMessage, error)
	CommitMessage(ctx context.Context, inputMessage *InputMessage) error
	Close(ctx context.Context) error
}
