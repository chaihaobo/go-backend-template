package controller

import (
	"github.com/chaihaobo/be-template/application"
	"github.com/chaihaobo/be-template/resource"
	"github.com/chaihaobo/be-template/transport/http/controller/health"
	"github.com/chaihaobo/be-template/transport/http/controller/user"
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

func New(res resource.Resource, app application.Application) Controller {
	return &controllers{
		healthController: health.NewController(res, app),
		userController:   user.NewController(res, app),
	}
}
