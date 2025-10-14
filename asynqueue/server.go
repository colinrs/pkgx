package asynqueue

import (
	"fmt"
	"github.com/colinrs/pkgx/utils"
	"log"
	"net/http"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/zeromicro/go-zero/core/logx"
)

type Server interface {
	Start(processWrappers ...[]ProcessWrapper) error
	RegisterHandle(pattern string, handler asynq.Handler)
	Shutdown()
}

type server struct {
	asynqServer *asynq.Server
	inspector   *asynq.Inspector
	serveMux    *asynq.ServeMux
	config      *AsynqConfig
}

var (
	defaultPriority = 10
)

// NewAsynqServer https://github.com/hibiken/asynq
func NewAsynqServer(config AsynqConfig) Server {
	logx.Must(config.Validate())
	serveMux := asynq.NewServeMux()
	serveMux.Use(metricMiddleware())
	serveMux.Use(traceMiddleware())
	return &server{
		asynqServer: initAsynqServer(config.Addr, config.Concurrency, config.Queue),
		inspector:   asynq.NewInspector(asynq.RedisClientOpt{Addr: config.Addr}),
		serveMux:    serveMux,
		config:      &config,
	}
}

func initAsynqServer(addr string, concurrency int, queue string) *asynq.Server {
	return asynq.NewServer(asynq.RedisClientOpt{Addr: addr}, asynq.Config{
		// Specify how many concurrent workers to use
		Concurrency:     concurrency,
		ShutdownTimeout: 25,
		// Optionally specify multiple queues with different priority.
		Queues: map[string]int{
			queue: defaultPriority,
		},
		// See the godoc for other configuration options
	})
}

func asynqmonHttpServer(config *AsynqConfig) {
	if config.Monitoring == nil {
		return
	}
	if !config.Monitoring.Enable {
		logx.Infof("asynqmon is disabled")
		return
	}
	logx.Infof("Starting asynqmon server at :8089")
	h := asynqmon.New(asynqmon.Options{
		RootPath:     config.Monitoring.Path, // RootPath specifies the root for asynqmon app
		RedisConnOpt: asynq.RedisClientOpt{Addr: config.Addr},
	})
	// Note: We need the tailing slash when using net/http.ServeMux.
	http.Handle(h.RootPath()+"/", utils.BasicAuth(h, config.Monitoring.UserName, config.Monitoring.Password))
	// Go to http://localhost:8080/monitoring to see asynqmon homepage.
	logx.Must(http.ListenAndServe(fmt.Sprintf(":%s", config.Monitoring.Port), nil))
}

func (s *server) Start(processWrappers ...[]ProcessWrapper) error {
	// 启动队列指标收集器
	startQueueMetricsCollector(s.inspector)
	go asynqmonHttpServer(s.config)
	// mux maps a type to a handler
	if len(processWrappers) > 0 {
		for _, processWrapper := range processWrappers[0] {
			s.RegisterHandle(processWrapper.EventName(), processWrapper.Handler())
		}
	}
	fmt.Printf("Starting asynq server")
	if err := s.asynqServer.Start(s.serveMux); err != nil {
		log.Fatalf("could not run asynq server: %v", err)
		return err
	}
	return nil
}

func (s *server) RegisterHandle(pattern string, handler asynq.Handler) {
	s.serveMux.Handle(pattern, handler)
}

func (s *server) Shutdown() {
	s.asynqServer.Shutdown()
	if s.inspector != nil {
		s.inspector.Close()
	}
}
