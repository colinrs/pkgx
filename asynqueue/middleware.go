package asynqueue

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func traceMiddleware() asynq.MiddlewareFunc {
	return func(next asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
			tracer := otel.Tracer("asynqueue")

			spanName := t.Type()
			taskId := t.ResultWriter().TaskID()
			spanCtx, span := tracer.Start(
				ctx,
				spanName,
				trace.WithSpanKind(trace.SpanKindServer),
			)
			defer span.End()

			start := time.Now()
			// do something before
			err := next.ProcessTask(spanCtx, t)
			// do something after
			span.SetAttributes(attribute.Int64("duration", time.Since(start).Milliseconds()))
			span.SetAttributes(attribute.String("asynq_event", t.Type()))
			span.SetAttributes(attribute.String("event_trace_id", getTraceIDFromTaskId(taskId)))
			if err != nil {
				span.SetAttributes(attribute.String("err", err.Error()))
				span.SetStatus(codes.Error, "task completed")
				return err
			}
			span.SetStatus(codes.Ok, "task completed")
			return nil
		})
	}
}

func metricMiddleware() asynq.MiddlewareFunc {
	return func(next asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
			logx.WithContext(ctx).Debugf("start process task %s,payload:%s", t.ResultWriter().TaskID(), t.Payload())
			now := time.Now()
			err := next.ProcessTask(ctx, t)
			status := Success
			if err != nil {
				status = Failed
			}
			logx.WithContext(ctx).Debugf("end process task %s,payload:%s", t.ResultWriter().TaskID(), t.Payload())
			EventTotal.Inc(t.Type(), status)
			EventDurationHistogram.ObserveFloat(float64(time.Since(now).Milliseconds()), t.Type())
			return err
		})
	}
}
