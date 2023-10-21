package kafka

import (
	"context"
	"fmt"

	"github.com/colinrs/pkgx/kq/core"

	"github.com/IBM/sarama"
)

func NewOutput(ctx context.Context) (core.Output, error) {
	return newProducer(), nil
}

type producer struct {
	client sarama.SyncProducer
}

func newProducer() *producer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		fmt.Println(err)
	}
	return &producer{client: client}
}

func (p *producer) SendOutput(ctx context.Context, msg *core.OutputMessage) error {
	producerMessage := &sarama.ProducerMessage{
		Topic: "shopping",
		Value: sarama.StringEncoder("20220411happy02"),
	}
	partition, offset, err := p.client.SendMessage(producerMessage)
	fmt.Println("partition:")
	fmt.Println(partition)
	fmt.Println("offset:")
	fmt.Println(offset)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (p *producer) Close(ctx context.Context) error {
	return p.client.Close()
}
