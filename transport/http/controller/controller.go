package controller

import (
	"github.com/google/wire"

	"github.com/chaihaobo/be-template/transport/http/controller/health"
	"github.com/chaihaobo/be-template/transport/http/controller/user"
)

var ProviderSet = wire.NewSet(
	health.NewController,
	user.NewController,
	New,
)

type (
	Controller interface {
		Health() health.Controller
		User() user.Controller
	}

	controllers struct {
		healthController health.Controller
		userController   user.Controller
	}
)

func (c *controllers) User() user.Controller {
	return c.userController
}

func (c *controllers) Health() health.Controller {
	return c.healthController
}

func New(healthController health.Controller, userController user.Controller) Controller {
	return &controllers{
		healthController: healthController,
		userController:   userController,
	}
}
