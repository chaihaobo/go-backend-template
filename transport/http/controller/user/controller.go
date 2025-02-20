package user

import (
	"github.com/gin-gonic/gin"

	"github.com/chaihaobo/be-template/application"
	"github.com/chaihaobo/be-template/model/dto/user"
	"github.com/chaihaobo/be-template/resource"
)

type (
	Controller interface {
		Login(ctx *gin.Context) (*user.LoginResponse, error)
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
