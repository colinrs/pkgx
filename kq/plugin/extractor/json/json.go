package json

import (
	"context"

	"github.com/colinrs/pkgx/kq/core"
)

type extractor struct {
}

func NewExtractor(ctx context.Context) core.Extractor {
	return nil
}

func (e *extractor) Unmarshal(ctx context.Context, inputMessage *core.InputMessage) (interface{}, error) {
	return nil, nil
}

func (e *extractor) OnDone(cxt context.Context, inputMessage *core.InputMessage) {
	return
}

func (e *extractor) OnError(cxt context.Context, inputMessage *core.InputMessage) (interface{}, error) {
	return nil, nil
}
