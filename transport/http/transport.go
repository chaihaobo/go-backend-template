package http

import (
	"context"
	"fmt"
	"net/http"

	ginmiddewate "github.com/chaihaobo/gocommon/middleware/http/gin"
	"github.com/chaihaobo/gocommon/restkit"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/chaihaobo/be-template/application"
	"github.com/chaihaobo/be-template/infrastructure"
	"github.com/chaihaobo/be-template/model/dto/user"
	"github.com/chaihaobo/be-template/resource"
	"github.com/chaihaobo/be-template/resource/discovery"
	"github.com/chaihaobo/be-template/transport/http/controller"
	"github.com/chaihaobo/be-template/transport/http/middleware"
	"github.com/chaihaobo/be-template/utils"
)

type (
	Transport interface {
		Serve() error
		Shutdown() error
	}

	transport struct {
		resource   resource.Resource
		engine     *gin.Engine
		infra      infrastructure.Infrastructure
		controller controller.Controller
		server     *http.Server
		serviceID  string
	}
)

func (t *transport) Serve() error {
	var (
		port = t.resource.Configuration().Service.HTTPPort
		name = t.resource.Configuration().Service.Name
	)
	ip, err := utils.GetOutboundIP()
	if err != nil {
		return err
	}
	ctx := context.TODO()
	serviceID, err := t.resource.Discovery().RegisterService(ctx, &discovery.Service{
		Name:            fmt.Sprintf("%s-http", name),
		IP:              ip,
		Port:            port,
		Type:            discovery.ServiceTypeHTTP,
		HealthCheckPath: "/health",
	})
	if err != nil {
		return fmt.Errorf("failed to register service to consul: %w", err)
	}
	t.serviceID = serviceID
	t.resource.Logger().Info(ctx, "http server started",
		zap.String("name", name),
		zap.Int("port", port))
	return t.server.ListenAndServe()
}

func (t *transport) Shutdown() error {
	ctx := context.TODO()
	if err := t.resource.Discovery().DeregisterService(ctx, t.serviceID); err != nil {
		t.resource.Logger().Error(ctx, "failed to deregister http service from consul", err)
	}
	return t.server.Shutdown(ctx)
}

func (t *transport) applyRoutes() {
	router := t.engine
	healthController := t.controller.Health()
	userController := t.controller.User()
	router.GET("/health", restkit.AdaptToGinHandler(restkit.HandlerFunc[any](healthController.Health)))

	userGroup := router.Group("/user")
	{
		userGroup.POST("/login", restkit.AdaptToGinHandler(restkit.HandlerFunc[*user.LoginResponse](userController.Login)))
	}

}

func NewTransport(res resource.Resource, infra infrastructure.Infrastructure, app application.Application) Transport {
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
		infra:      infra,
		controller: controller.New(res, app),
		server:     &http.Server{Addr: fmt.Sprintf(":%d", svcConfig.HTTPPort), Handler: engine},
	}
	tsp.applyRoutes()
	return tsp
}
