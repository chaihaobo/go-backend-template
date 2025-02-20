package health

import (
	"context"

	"github.com/chaihaobo/be-template/constant"
	"github.com/chaihaobo/be-template/infrastructure"
	"github.com/chaihaobo/be-template/resource"
)

type (
	Service interface {
		HealthCheck(ctx context.Context) error
	}

	service struct {
		res   resource.Resource
		infra infrastructure.Infrastructure
	}
)

func (s *service) HealthCheck(ctx context.Context) error {
	var healthChecks = []func(context.Context) error{
		s.infra.Store().Client().Ping,
		s.infra.Cache().Ping,
	}
	for _, check := range healthChecks {
		if err := check(ctx); err != nil {
			s.res.Logger().Error(ctx, "health check failed", err)
			return constant.ErrHealthCheckFailed
		}
	}
	return nil
}

func NewService(res resource.Resource, infra infrastructure.Infrastructure) Service {
	return &service{
		res:   res,
		infra: infra,
	}
}
