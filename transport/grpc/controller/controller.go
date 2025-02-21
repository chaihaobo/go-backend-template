package controller

import (
	"github.com/google/wire"

	"github.com/chaihaobo/be-template/application"
	"github.com/chaihaobo/be-template/resource"
	"github.com/chaihaobo/be-template/transport/grpc/controller/hello"
)

var ProviderSet = wire.NewSet(
	hello.NewController,
	NewController,
)

type (
	Controller interface {
		Hello() hello.Controller
	}
	controller struct {
		hello hello.Controller
	}
)

func (c controller) Hello() hello.Controller {
	return c.hello
}

func NewController(res resource.Resource, app application.Application) Controller {
	return &controller{
		hello: hello.NewController(res, app),
	}
}
