package application

import (
	"github.com/chaihaobo/be-template/application/health"
	"github.com/chaihaobo/be-template/application/user"
	"github.com/chaihaobo/be-template/infrastructure"
	"github.com/chaihaobo/be-template/resource"
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
