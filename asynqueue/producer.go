package asynqueue

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type Producer interface {
	Close() error
	NewTaskCtx(ctx context.Context, taskTypeName string, data Marshaller) error
	NewTaskDelayCtx(ctx context.Context, taskTypeName string, data Marshaller, processIn time.Duration) error
}

type Marshaller interface {
	Marshal() []byte
}

type producer struct {
	queue    string
	asynqCli *asynq.Client
	logx.Logger
}

func NewAsynqProducer(config AsynqConfig) Producer {
	logx.Must(config.Validate())
	return &producer{
		asynqCli: asynq.NewClient(asynq.RedisClientOpt{Addr: config.Addr, DB: config.DB}),
		queue:    config.Queue,
		Logger:   logx.WithContext(context.TODO()),
	}
}

func (c *producer) Close() error {
	return c.asynqCli.Close()
}

func (c *producer) NewTaskDelayCtx(ctx context.Context, taskTypeName string, data Marshaller, processIn time.Duration) error {
	taskId := c.getTaskId(taskTypeName, trace.TraceIDFromContext(ctx))
	var opts []asynq.Option
	opts = append(opts, asynq.Queue(c.queue))
	opts = append(opts, asynq.TaskID(taskId))
	opts = append(opts, asynq.ProcessIn(processIn))
	return c.newTaskCtx(ctx, taskTypeName, data, taskId, opts...)
}

func (c *producer) NewTaskCtx(ctx context.Context, taskTypeName string, data Marshaller) error {
	taskId := c.getTaskId(taskTypeName, trace.TraceIDFromContext(ctx))
	var opts []asynq.Option
	opts = append(opts, asynq.Queue(c.queue))
	opts = append(opts, asynq.TaskID(taskId))
	return c.newTaskCtx(ctx, taskTypeName, data, taskId, opts...)
}

func (c *producer) newTaskCtx(ctx context.Context, taskTypeName string, data Marshaller, taskId string, opts ...asynq.Option) error {
	logger := c.Logger.WithContext(ctx)
	tracer := trace.TracerFromContext(ctx)
	spanCtx, span := tracer.Start(
		ctx,
		taskTypeName,
		oteltrace.WithSpanKind(oteltrace.SpanKindProducer),
	)
	defer span.End()
	start := time.Now()
	info, err := c.asynqCli.EnqueueContext(spanCtx, asynq.NewTask(taskTypeName, data.Marshal()), opts...)
	span.SetAttributes(attribute.Int64("duration", time.Since(start).Milliseconds()))
	span.SetAttributes(attribute.String("asynq_event", taskTypeName))
	span.SetAttributes(attribute.String("event_trace_id", trace.TraceIDFromContext(ctx)))
	status := Success
	if err != nil {
		span.SetAttributes(attribute.String("err", err.Error()))
		span.SetStatus(codes.Error, "task push err")
		status = Failed
	} else {
		span.SetStatus(codes.Ok, "task push completed")
	}
	ProducerTotal.Inc(taskTypeName, status)
	if err != nil {
		logger.Errorf("asynq: data:%v, task: id=%s  could not enqueue task: %v", data, taskId, err)
		return err
	}
	logger.Infof("asynq: data:%v, enqueued task: id=%s queue=%s", data, taskId, info.Queue)
	return nil
}

func (c *producer) getTaskId(eventName, traceID string) string {
	return getTaskId(eventName, traceID)
}
