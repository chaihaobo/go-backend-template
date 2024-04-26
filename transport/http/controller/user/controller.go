package user

import (
	"github.com/gin-gonic/gin"

	"gitlab.seakoi.net/engineer/backend/be-template/application"
	"gitlab.seakoi.net/engineer/backend/be-template/resource"
)

type (
	Controller interface {
		Login(ctx *gin.Context)
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
