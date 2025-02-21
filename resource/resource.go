package resource

import (
	"context"

	"github.com/chaihaobo/be-template/resource/config"
	"github.com/chaihaobo/be-template/resource/logger"
	"github.com/chaihaobo/be-template/resource/metric"
	"github.com/chaihaobo/be-template/resource/tracer"
	"github.com/chaihaobo/be-template/resource/validator"
)

type (
	Resource interface {
		Configuration() *config.Configuration
		Logger() logger.Logger
		Validator() validator.Validator
		Metric() metric.PrometheusMetric
		Tracer() tracer.Tracer
		Close() error
	}

	resource struct {
		configuration *config.Configuration
		logger        logger.Logger
		validator     validator.Validator
		metric        metric.PrometheusMetric
		tracer        tracer.Tracer
	}
)

func (r *resource) Tracer() tracer.Tracer {
	return r.tracer
}

func (r *resource) Metric() metric.PrometheusMetric {
	return r.metric
}

func (r *resource) Validator() validator.Validator {
	return r.validator
}

func (r *resource) Logger() logger.Logger {
	return r.logger
}

func (r *resource) Close() error {
	ctx := context.Background()
	closeFuncs := []func() error{
		func() error {
			return r.metric.Close(ctx)
		},
		func() error {
			return r.tracer.Close(ctx)
		},
	}
	for _, closeFun := range closeFuncs {
		if err := closeFun(); err != nil {
			return err
		}
	}

	return nil
}

func (r *resource) Configuration() *config.Configuration {
	return r.configuration
}

func NewResource(configuration *config.Configuration, log logger.Logger, validator validator.Validator, prometheusMetric metric.PrometheusMetric, tracer tracer.Tracer) Resource {
	return &resource{
		configuration: configuration,
		logger:        log,
		validator:     validator,
		metric:        prometheusMetric,
		tracer:        tracer,
	}
}
