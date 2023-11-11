package kafka

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

type ConsumeFunc func(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error

// Consumer is to consume message from kafka and process it.
type Consumer interface {
	ConsumeMessage() (*Message, bool)
	Stop(ctx context.Context)
}

// NewConsumer will return a kafka consumer and you can use Process to consume messages.
func NewConsumer(ctx context.Context, brokers []string,
	username, password, groupID string, topics []string, options ...Option) (Consumer, error) {
	opts := &Options{}
	for _, option := range options {
		option(opts)
	}
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V1_0_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.MaxProcessingTime = 2 * time.Second
	if username != "" {
		config.Version = sarama.V2_3_0_0
		config.Net.SASL.Enable = true
		config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		config.Net.SASL.User = username
		config.Net.SASL.Password = password
		config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &XDGSCRAMClient{HashGeneratorFcn: SHA512}
		}
		if opts.useSASLPlainText {
			config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
			config.Net.SASL.SCRAMClientGeneratorFunc = nil
		}
	}
	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Printf("Failed to create consumer group client: err=[%v]", err)
		return nil, err
	}

	//listen client's errors.
	go func() {
		for err := range client.Errors() {
			log.Printf("Consumer got errors. err=[%v]", err)
		}
	}()
	c := &consumer{
		client:  client,
		groupID: groupID,
		topics:  topics,
		message: make(chan *Message),
	}
	go c.consumeMessage(ctx)
	return c, nil
}

type consumer struct {
	client  sarama.ConsumerGroup
	groupID string
	message chan *Message
	topics  []string
}

func (c *consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			c.message <- &Message{
				consumerMessage: message,
				session:         session,
			}
		case <-session.Context().Done():
			log.Printf("topics:%+v, session done", c.topics)
			return nil
		}
	}
}

// consumeMessage begins to consume messages from kafka in topics.
func (c *consumer) consumeMessage(parentCtx context.Context) {
	ctx, cancel := context.WithCancel(parentCtx)
	wg := &sync.WaitGroup{}
	topics := c.topics
	go func() {
		wg.Add(1)
		defer wg.Done()
		for {
			select {
			case <-parentCtx.Done():
				log.Printf("Consume parentCtx was cancel |topics=%v\n", topics)
				return
			default:
				if err := c.client.Consume(ctx, topics, c); err != nil {
					log.Printf("error occur from consumer topics %s: err=[%v]", topics, err)
				}
				// check if context was cancelled, signaling that the consumer should stop
				if ctx.Err() != nil {
					log.Printf("consumer exit. topics:%+v, err=[%v]", topics, ctx.Err())
					return
				}
				log.Printf("reconnect kafka topics:%+v", topics)
			}
		}
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("consumer terminating: context cancelled")
	case <-sigterm:
		log.Println("consumer terminating: via signal")
	}
	cancel()
	wg.Wait()
	c.Stop(ctx)
}

func (c *consumer) Stop(_ context.Context) {
	close(c.message)
	if err := c.client.Close(); err != nil {
		log.Printf("failed to close consumer client: err=[%v]", err)
	}
}

func (c *consumer) ConsumeMessage() (*Message, bool) {
	msg, end := <-c.message
	return msg, end
}

type Message struct {
	consumerMessage *sarama.ConsumerMessage
	session         sarama.ConsumerGroupSession
}

func (m *Message) ConsumerMessage() *sarama.ConsumerMessage {
	return m.consumerMessage
}

func (m *Message) Session() sarama.ConsumerGroupSession {
	return m.session
}
