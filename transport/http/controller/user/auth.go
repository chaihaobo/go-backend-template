package user

import (
	"github.com/gin-gonic/gin"

	"gitlab.seakoi.net/engineer/backend/be-template/constant"
	"gitlab.seakoi.net/engineer/backend/be-template/model/dto/user"
	"gitlab.seakoi.net/engineer/backend/be-template/tools"
)

func (c *controller) Login(gctx *gin.Context) {
	ctx := gctx.Request.Context()
	request := new(user.LoginRequest)

	if err := gctx.ShouldBindJSON(request); err != nil {
		c.res.Logger().Error(ctx, "failed to bind login request", err)
		tools.HTTPWriteErr(gctx.Writer, constant.ErrSystemMalfunction)
		return
	}

	response, err := c.app.User().Login(ctx, request)
	if err != nil {
		tools.HTTPWriteErr(gctx.Writer, err)
		return
	}
	tools.HTTPWrite(gctx.Writer, response)

}
