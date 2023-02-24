package common

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"time"
)

// 创建链路追踪实例

func NewTracer(serviceName string, addr string) (opentracing.Tracer, io.Closer, error) {
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			QueueSize:                  0,
			BufferFlushInterval:        1 * time.Second,
			LogSpans:                   true,
			LocalAgentHostPort:         addr,
			DisableAttemptReconnecting: false,
			AttemptReconnectInterval:   0,
			CollectorEndpoint:          "",
			User:                       "",
			Password:                   "",
			HTTPHeaders:                nil,
		},
	}

	return cfg.NewTracer()
}
