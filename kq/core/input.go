package core

import (
	"context"
	"sync"

	"github.com/IBM/sarama"
)

type InputMessage struct {
	Raw  *sarama.ConsumerMessage
	done sync.Mutex
}

// Ack acknowledges that the message has been processed.
// It will release the internal lock on the message, allowing the message to be marked in kafka.
func (m *InputMessage) Ack() {
	m.done.Unlock()
}

type Input interface {
	Consumer(ctx context.Context) (*InputMessage, error)
	Close(ctx context.Context) error
}
