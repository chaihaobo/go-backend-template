package application

import (
	"github.com/google/wire"

	"github.com/chaihaobo/be-template/application/health"
	"github.com/chaihaobo/be-template/application/user"
)

var ProviderSet = wire.NewSet(
	health.NewService,
	user.NewService,
	New,
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

func New(
	healthService health.Service,
	userService user.Service,
) Application {
	return &application{
		health: healthService,
		user:   userService,
	}
}
