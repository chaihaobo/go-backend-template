package http

import (
	"context"
	"net/http"

	ginmiddewate "github.com/chaihaobo/gocommon/middleware/http/gin"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"gitlab.seakoi.net/engineer/backend/be-template/application"
	"gitlab.seakoi.net/engineer/backend/be-template/resource"
	"gitlab.seakoi.net/engineer/backend/be-template/transport/http/controller"
	"gitlab.seakoi.net/engineer/backend/be-template/transport/http/middleware"
)

type (
	Transport interface {
		Serve() error
		Shutdown() error
	}

	transport struct {
		resource   resource.Resource
		engine     *gin.Engine
		controller controller.Controller
		server     *http.Server
	}
)

func (t *transport) Serve() error {
	var (
		addr = t.resource.Configuration().Service.HTTPPort
		name = t.resource.Configuration().Service.Name
	)
	t.resource.Logger().Info(context.TODO(), "http server started",
		zap.String("name", name),
		zap.String("addr", addr))
	return t.server.ListenAndServe()
}

func (t *transport) Shutdown() error {
	return t.server.Shutdown(context.TODO())
}

func (t *transport) applyRoutes() {
	router := t.engine
	healthController := t.controller.Health()
	userController := t.controller.User()
	router.GET("/health", healthController.Health)

	userGroup := router.Group("/user")
	{
		userGroup.POST("/login", userController.Login)
	}

}

// middware
func NewTransport(res resource.Resource, app application.Application) Transport {
	svcConfig := res.Configuration().Service
	gin.SetMode(lo.If(svcConfig.Debug, gin.DebugMode).
		Else(gin.ReleaseMode))
	engine := gin.New()
	engine.Use(
		ginmiddewate.TelemetryMiddleware(svcConfig.Name, svcConfig.Env, res.Logger()),
		gin.Recovery(),
		middleware.AuthMiddleware(res, app),
	)

	tsp := &transport{
		resource:   res,
		engine:     engine,
		controller: controller.New(res, app),
		server:     &http.Server{Addr: svcConfig.HTTPPort, Handler: engine},
	}
	tsp.applyRoutes()
	return tsp
}
