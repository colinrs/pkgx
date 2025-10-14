package asynqueue

import (
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynq/x/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/metric"
)

// StartQueueMetricsCollector starts the official asynq queue metrics collector and registers it with Prometheus
func startQueueMetricsCollector(inspector *asynq.Inspector) {
	// Use the inspector from ServiceContext
	collector := metrics.NewQueueMetricsCollector(inspector)

	// Register with Prometheus default registry
	// Note: go-zero's /metrics endpoint uses the default prometheus.DefaultRegisterer
	// So using prometheus.Register here is correct
	if err := prometheus.Register(collector); err != nil {
		logx.Errorf("Failed to register asynq metrics collector: %v", err)
		return
	}
	logx.Info("Official asynq queue metrics collector registered with Prometheus")
}

const (
	serverNamespace   = "go_asynq_server"
	producerNamespace = "go_asynq_producer"
)

// EventTotal instance
var (
	EventTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "event",
		Name:      "total",
		Help:      "go common asynq server event custom error count.",
		Labels:    []string{"name", "status"},
	})
	EventDurationHistogram = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Subsystem: "event",
		Name:      "duration_ms",
		Help:      "go common asynq event duration(ms).",
		Labels:    []string{"name"},
		Buckets:   []float64{0.25, 0.5, 1, 2, 5, 10, 25, 50, 100, 250, 500, 1000, 2000, 5000, 10000, 15000},
	})

	ProducerTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: producerNamespace,
		Subsystem: "event",
		Name:      "total",
		Help:      "go common asynq server event custom error count.",
		Labels:    []string{"name", "status"},
	})
)
