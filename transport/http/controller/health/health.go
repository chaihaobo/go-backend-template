package health

import (
	"github.com/gin-gonic/gin"
)

func (c *controller) Health(ctx *gin.Context) (any, error) {
	err := c.app.Health().HealthCheck(ctx.Request.Context())
	return nil, err
}
