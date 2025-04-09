package health

import (
	"github.com/gin-gonic/gin"

	"github.com/chaihaobo/be-template/application"
	"github.com/chaihaobo/be-template/resource"
)

type (
	Controller interface {
		Health(ctx *gin.Context) (any, error)
	}

	controller struct {
		app application.Application
		res resource.Resource
	}
)

func NewController(res resource.Resource, app application.Application) Controller {
	return &controller{
		app: app,
		res: res,
	}
}
