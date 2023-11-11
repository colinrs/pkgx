package kafka

import (
	"context"
	"fmt"
)

func NewConsumerExample() {
	ctx := context.Background()
	c, err := NewConsumer(ctx, []string{"127.0.0.1:9092"}, "", "", "test", []string{"test"})
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		msg, ok := c.ConsumeMessage()
		if !ok {
			break
		}
		fmt.Println(msg)
		msg.Session().MarkMessage(msg.ConsumerMessage(), "")
	}
}
