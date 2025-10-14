package asynqueue

import "context"

// entry ...
type entry struct {
	event      string
	newPayload func() any
	handle     func(ctx context.Context, v any) error
	decoder    Decoder
}

// WrapperProcessWithDecoder ---- 包级“泛型函数”做适配
func wrapperProcessWithDecoder[T any](p Process[T], dec Decoder) *entry {
	return &entry{
		event: p.EventName(),
		newPayload: func() any {
			return new(PayLoad[T]) // *T
		},
		handle: func(ctx context.Context, v any) error {
			return p.Handle(ctx, *v.(*PayLoad[T]))
		},
		decoder: dec,
	}
}

// wrapperProcess ---- 默认 JSON 解码器的适配 ----
func wrapperProcess[T any](p Process[T]) *entry {
	return wrapperProcessWithDecoder(p, JSONDecoder)
}
