package resource

import (
	"context"

	"github.com/chaihaobo/gocommon/constant"
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"

	svcconstant "github.com/chaihaobo/be-template/constant"
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
		loggerFlusher func() error
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
		r.loggerFlusher,
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

func New(configPath string) (Resource, error) {
	configuration, err := config.NewConfiguration(configPath)
	if err != nil {
		return nil, err
	}
	constant.MergeServiceErrorCode2HTTPStatus(svcconstant.ServiceErrorCode2HTTPStatus)
	logConfig := configuration.Logger
	logger, f, err := logger.New(logger.Config{
		FileName:   logConfig.FileName,
		MaxSize:    logConfig.MaxSize,
		MaxAge:     logConfig.MaxAge,
		WithCaller: true,
		CallerSkip: 1,
	})
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	configuration.OnConfigChange(func(e fsnotify.Event) {
		logger.Info(ctx, "config changed", zap.Any("event", e))
		if err := configuration.Unmarshal(configuration); err != nil {
			logger.Error(ctx, "failed to unmarshal changed config", err)
		}
	})

	validator, err := validator.NewValidator()
	if err != nil {
		return nil, err

	}

	prometheusMetric, err := metric.NewPrometheusMetric(configuration)
	if err != nil {
		return nil, err
	}

	tracer, err := tracer.NewTracer(configuration)
	if err != nil {
		return nil, err
	}

	return &resource{
		configuration: configuration,
		logger:        logger,
		validator:     validator,
		loggerFlusher: f,
		metric:        prometheusMetric,
		tracer:        tracer,
	}, nil
}
