package kafka

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"github.com/colinrs/pkgx/logger"
	"log"
)

func newConcurrentConsume(ctx context.Context, client sarama.ConsumerGroup) *ConcurrentConsume {
	return &ConcurrentConsume{
		Client: client,
		ready:  make(chan bool),
	}
}

type ConcurrentConsume struct {
	Client sarama.ConsumerGroup
	ready  chan bool
}

func (c *ConcurrentConsume) StartConcurrentConsume(ctx context.Context, topics []string, consumerAPI Consumer) error {
	go func() {
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := c.Client.Consume(ctx, topics, c); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			c.ready = make(chan bool)
		}
	}()
	return nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *ConcurrentConsume) Setup(sarama.ConsumerGroupSession) error {
	// Mark the c as ready
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *ConcurrentConsume) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (c *ConcurrentConsume) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			logger.Info("Message claimed: value = %s, timestamp = %v, topic = %s",
				string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")
		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
