package application

import (
	"gitlab.seakoi.net/engineer/backend/be-template/application/health"
	"gitlab.seakoi.net/engineer/backend/be-template/application/user"
	"gitlab.seakoi.net/engineer/backend/be-template/infrastructure"
	"gitlab.seakoi.net/engineer/backend/be-template/resource"
)

type (
	Application interface {
		Health() health.Service
		User() user.Service
	}

	application struct {
		health health.Service
		user   user.Service
	}
)

func (a *application) User() user.Service {
	return a.user
}

func (a *application) Health() health.Service {
	return a.health
}

func New(res resource.Resource, infra infrastructure.Infrastructure) Application {
	return &application{
		health: health.NewService(res, infra),
		user:   user.NewService(res, infra),
	}
}
