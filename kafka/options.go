package kafka

import (
	"github.com/IBM/sarama"
)

type Options struct {
	useSASLPlainText     bool
	producerInterceptors []sarama.ProducerInterceptor
	consumerInterceptors []sarama.ConsumerInterceptor
}

type Option func(*Options)

// WithSASLPlainText use SASLTypePlaintext
func WithSASLPlainText() func(*Options) {
	return func(o *Options) {
		o.useSASLPlainText = true
	}
}

// WithProducerInterceptors ...
func WithProducerInterceptors(producerInterceptors []sarama.ProducerInterceptor) func(*Options) {
	return func(o *Options) {
		o.producerInterceptors = append(o.producerInterceptors, producerInterceptors...)
	}
}

// WithConsumerInterceptors ...
func WithConsumerInterceptors(consumerInterceptors []sarama.ConsumerInterceptor) func(*Options) {
	return func(o *Options) {
		o.consumerInterceptors = append(o.consumerInterceptors, consumerInterceptors...)
	}
}
