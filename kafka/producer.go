package kafka

import (
	"context"
	"crypto/sha512"
	"time"

	"github.com/IBM/sarama"
	"github.com/xdg/scram"
)

var _ Producer = (*producer)(nil)

var (
	sha512Instance scram.HashGeneratorFcn = sha512.New
)

const (
	MaxRequestSize = 10 * 1024 * 1024
)

type Producer interface {
	SyncProduce(ctx context.Context, topic string, message []byte) (partition int32, offset int64, err error)
	SyncProduceMessages(ctx context.Context, topic string, messages [][]byte) error
	SyncProduceWithKey(ctx context.Context, topic string, key string,
		content []byte) (partition int32, offset int64, err error)
	SyncProduceWithHeader(ctx context.Context, topic string, content []byte,
		headers []sarama.RecordHeader) (partition int32, offset int64, err error)
	SyncProduceWithKeyHeader(ctx context.Context, topic string, key string, content []byte,
		headers []sarama.RecordHeader) (partition int32, offset int64, err error)
}

type producer struct {
	producer sarama.SyncProducer
	brokers  []string
}

//NewProducer Build new kafka producer
func NewProducer(ctx context.Context, brokers []string, username, password string, options ...Option) (Producer,
	error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Timeout = 5 * time.Second
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Version = sarama.V2_3_0_0
	opts := &Options{}
	for _, option := range options {
		option(opts)
	}
	if len(opts.producerInterceptors) > 0 {
		kafkaConfig.Producer.Interceptors = opts.producerInterceptors
	}
	if username != "" {
		kafkaConfig.Version = sarama.V2_3_0_0
		kafkaConfig.Net.SASL.Enable = true
		kafkaConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		kafkaConfig.Net.SASL.User = username
		kafkaConfig.Net.SASL.Password = password
		kafkaConfig.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &XDGSCRAMClient{HashGeneratorFcn: sha512Instance}
		}
		if opts.useSASLPlainText {
			kafkaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
			kafkaConfig.Net.SASL.SCRAMClientGeneratorFunc = nil
		}
	}
	//default 100M, should <= server config
	sarama.MaxRequestSize = MaxRequestSize
	p, err := sarama.NewSyncProducer(brokers, kafkaConfig)
	if err != nil {
		return nil, err
	}

	producerInstance := &producer{
		producer: p,
		brokers:  brokers,
	}
	return producerInstance, nil
}

func (p *producer) SyncProduce(ctx context.Context, topic string, message []byte) (int32, int64, error) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}
	return p.producer.SendMessage(msg)
}

func (p *producer) SyncProduceMessages(ctx context.Context, topic string, messages [][]byte) error {
	messagesToMq := make([]*sarama.ProducerMessage, 0, len(messages))
	for _, message := range messages {
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(message),
		}
		messagesToMq = append(messagesToMq, msg)
	}
	return p.producer.SendMessages(messagesToMq)
}

func (p *producer) SyncProduceWithKey(ctx context.Context, topic string, key string, message []byte) (int32, int64,
	error) {
	msg := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(key),
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}
	return p.producer.SendMessage(msg)
}
func (p *producer) SyncProduceWithHeader(ctx context.Context, topic string, message []byte,
	headers []sarama.RecordHeader) (int32, int64, error) {
	msg := &sarama.ProducerMessage{
		Topic:   topic,
		Value:   sarama.ByteEncoder(message),
		Headers: headers,
	}
	return p.producer.SendMessage(msg)
}
func (p *producer) SyncProduceWithKeyHeader(ctx context.Context, topic string, key string, message []byte,
	headers []sarama.RecordHeader) (int32, int64, error) {
	msg := &sarama.ProducerMessage{
		Key:     sarama.StringEncoder(key),
		Topic:   topic,
		Value:   sarama.ByteEncoder(message),
		Headers: headers,
	}
	return p.producer.SendMessage(msg)
}
