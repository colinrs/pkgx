package kq

import (
	"context"
	"github.com/colinrs/pkgx/kq/core"

	extractorJson "github.com/colinrs/pkgx/kq/plugin/extractor/json"
	kafkaInput "github.com/colinrs/pkgx/kq/plugin/input/kafka"
	kafkaOutput "github.com/colinrs/pkgx/kq/plugin/output/kafka"
)

func example() {
	ctx := context.Background()
	k := NewKQ(ctx)
	input, _ := kafkaInput.NewInput(ctx)
	extrac := extractorJson.NewExtractor(ctx)
	output, _ := kafkaOutput.NewOutput(ctx)
	k.SetInput(ctx, input).
		SetExtractor(ctx, extrac).
		SetTransformer(ctx, &T{}).
		SetOutput(ctx, output)

}

type T struct {
}

func (t *T) Process(m core.Message) (core.Message, error) {

	return nil, nil
}

func (t *T) OnDone(message core.Message, resp interface{}) {
	return
}

func (t *T) OnError(message core.Message, err error) {
	return
}
