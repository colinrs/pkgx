package kafka

import (
	"context"

	"github.com/colinrs/pkgx/kq/core"
)

func NewInput(ctx context.Context) (core.Input, error) {
	return newConsumer(), nil
}

type consumer struct {
}

func newConsumer() *consumer {
	return &consumer{}
}

func (c *consumer) Consumer(ctx context.Context) (*core.InputMessage, error) {
	return nil, nil
}

func (c *consumer) Close(ctx context.Context) error {
	return nil
}
