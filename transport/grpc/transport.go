package grpc

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	commonGrpc "github.com/chaihaobo/gocommon/middleware/grpc"

	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/chaihaobo/be-template/application"
	"github.com/chaihaobo/be-template/constant"
	"github.com/chaihaobo/be-template/infrastructure"
	"github.com/chaihaobo/be-template/infrastructure/discovery"
	"github.com/chaihaobo/be-template/proto"
	"github.com/chaihaobo/be-template/resource"
	"github.com/chaihaobo/be-template/transport/grpc/controller"
	"github.com/chaihaobo/be-template/utils"
)

type (
	Transport interface {
		Serve() error
		GracefulStop()
	}
	transport struct {
		resource   resource.Resource
		infra      infrastructure.Infrastructure
		app        application.Application
		controller controller.Controller
		server     *grpc.Server
		serviceID  string
	}
)

// Serve the Grpc Server
func (t *transport) Serve() error {
	grpcPort := t.resource.Configuration().Service.GrpcPort

	ip, err := utils.GetOutboundIP()
	if err != nil {
		return fmt.Errorf("failed to get outbound ip: %w", err)
	}
	serviceID, err := t.infra.DiscoveryClient().RegisterService(context.Background(), &discovery.Service{
		Name: fmt.Sprintf("%s-grpc", t.resource.Configuration().Service.Name),
		IP:   ip,
		Port: grpcPort,
		Type: discovery.ServiceTypeGRPC,
	})
	if err != nil {
		return fmt.Errorf("failed to register service to consul: %w", err)
	}
	t.serviceID = serviceID
	proto.RegisterHelloServiceServer(t.server, t.controller.Hello())
	healthpb.RegisterHealthServer(t.server, health.NewServer())
	reflection.Register(t.server)

	t.resource.Logger().Info(context.Background(), "grpc server started.", zap.Int("port", grpcPort))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return fmt.Errorf("failed to start grpc server: %w", err)
	}
	return t.server.Serve(lis)
}

// GracefulStop Grpc Server
func (t *transport) GracefulStop() {
	if err := t.infra.DiscoveryClient().DeregisterService(context.Background(), t.serviceID); err != nil {
		t.resource.Logger().Error(context.Background(), "failed to deregister service from consul", err)
	}
	t.server.GracefulStop()

}

func NewTransport(res resource.Resource, infra infrastructure.Infrastructure, app application.Application) Transport {
	serviceConfig := res.Configuration().Service
	commonIntercept := commonGrpc.WithDefault(constant.ServiceErrorCodeToGRPCErrorCode)

	chainedUnaryInterceptor := grpc.ChainUnaryInterceptor(
		commonGrpc.TelemetryUnaryServerInterceptor(serviceConfig.Name, serviceConfig.Env, res.Logger()),
	)
	interceptors := append(commonIntercept, chainedUnaryInterceptor)
	server := grpc.NewServer(interceptors...)
	return &transport{
		resource:   res,
		infra:      infra,
		app:        app,
		controller: controller.NewController(res, app),
		server:     server,
	}

}
