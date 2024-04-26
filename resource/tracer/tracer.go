package tracer

import (
	"github.com/chaihaobo/gocommon/trace"

	"gitlab.seakoi.net/engineer/backend/be-template/resource/config"
)

type (
	Tracer trace.CloseableTracer
)

func NewTracer(config *config.Configuration) (Tracer, error) {
	tracer, err := trace.NewZipkinTracer(trace.Config{
		CollectorURL: config.Service.TraceCollectorURL,
		ServiceName:  config.Service.Name,
	})
	if err != nil {
		return nil, err
	}

	return tracer, nil
}
