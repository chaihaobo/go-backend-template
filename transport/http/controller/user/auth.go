package user

import (
	"github.com/gin-gonic/gin"

	"github.com/chaihaobo/be-template/constant"
	"github.com/chaihaobo/be-template/model/dto/user"
)

func (c *controller) Login(gctx *gin.Context) (*user.LoginResponse, error) {
	ctx := gctx.Request.Context()
	request := new(user.LoginRequest)

	if err := gctx.ShouldBindJSON(request); err != nil {
		c.res.Logger().Error(ctx, "failed to bind login request", err)
		return nil, constant.ErrSystemMalfunction
	}
	return c.app.User().Login(ctx, request)
}
