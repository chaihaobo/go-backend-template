package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *controller) Health(ctx *gin.Context) {
	_ = c.app.Health().HealthCheck(ctx.Request.Context())
	ctx.JSON(http.StatusOK, "successful")
}
