package resource

import (
	"github.com/google/wire"

	"github.com/chaihaobo/be-template/resource/config"
	"github.com/chaihaobo/be-template/resource/logger"
	"github.com/chaihaobo/be-template/resource/metric"
	"github.com/chaihaobo/be-template/resource/tracer"
	"github.com/chaihaobo/be-template/resource/validator"
)

var ProviderSet = wire.NewSet(
	config.NewConfiguration,
	logger.New,
	validator.NewValidator,
	metric.NewPrometheusMetric,
	tracer.NewTracer,
	NewResource,
)
