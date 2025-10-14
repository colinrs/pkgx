package asynqueue

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
)

// Process 业务方实现这个接口：只关心自己的载荷 T
type Process[T any] interface {
	EventName() string
	Handle(ctx context.Context, m PayLoad[T]) error
}

// ProcessWrapper ---- 非泛型的 Wrapper ----
type ProcessWrapper interface {
	EventName() string
	RegisterHandle(server Server)
	Handler() asynq.Handler
	WithHook(h HookFunc)
}

type wrapper struct {
	entry *entry
	hook  HookFunc
}

func NewProcessWrapperWithDecoder[T any](p Process[T], dec Decoder) ProcessWrapper {
	e := wrapperProcessWithDecoder(p, dec)
	return &wrapper{
		entry: e}
}

func NewProcessWrapper[T any](p Process[T]) ProcessWrapper {
	e := wrapperProcess(p)
	return &wrapper{
		entry: e}
}

func (w *wrapper) EventName() string {
	return w.entry.event
}

func (w *wrapper) RegisterHandle(server Server) {
	server.RegisterHandle(w.EventName(), w.Handler())
}

func (w *wrapper) Handler() asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		payload := w.entry.newPayload()
		fields := []logx.LogField{
			logx.Field(EventField, t.Type()),
		}
		ctx = logx.WithFields(
			ctx,
			fields...,
		)
		if err := w.entry.decoder(t.Payload(), payload); err != nil {
			logx.WithContext(ctx).Errorf("asynq: could not decode task: %v, data:%s", err, t.Payload())
			return err
		}
		err := w.entry.handle(ctx, payload)
		if err != nil {
			logx.WithContext(ctx).Errorf("asynq: handle task err: %v, data:%s", err, t.Payload())
			return err
		}
		if w.hook != nil {
			if err := w.hook(ctx, t.Type(), payload); err != nil {
				logx.WithContext(ctx).Errorf("asynq: could not hook task: %v, data:%s", err, t.Payload())
				return err
			}
		}
		return nil
	})
}

func (w *wrapper) WithHook(h HookFunc) {
	w.hook = h
}
