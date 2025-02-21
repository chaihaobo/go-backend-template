package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/google/wire"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	commonGrpc "github.com/chaihaobo/gocommon/middleware/grpc"

	"github.com/chaihaobo/be-template/application"
	"github.com/chaihaobo/be-template/constant"
	"github.com/chaihaobo/be-template/proto"
	"github.com/chaihaobo/be-template/resource"
	"github.com/chaihaobo/be-template/transport/grpc/controller"
)

var ProviderSet = wire.NewSet(
	NewTransport,
	controller.ProviderSet,
)

type (
	Transport interface {
		Serve() error
		GracefulStop()
	}
	transport struct {
		resource   resource.Resource
		app        application.Application
		controller controller.Controller
		server     *grpc.Server
	}
)

// Serve the Grpc Server
func (t transport) Serve() error {
	proto.RegisterHelloServiceServer(t.server, t.controller.Hello())
	reflection.Register(t.server)
	grpcPort := t.resource.Configuration().Service.GrpcPort
	t.resource.Logger().Info(context.Background(), "grpc server started.", zap.String("addr", grpcPort))
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return fmt.Errorf("failed to start grpc server: %w", err)
	}
	return t.server.Serve(lis)
}

// GracefulStop Grpc Server
func (t transport) GracefulStop() {
	t.server.GracefulStop()
}

func NewTransport(res resource.Resource, app application.Application, controller controller.Controller) Transport {
	serviceConfig := res.Configuration().Service
	commonIntercept := commonGrpc.WithDefault(constant.ServiceErrorCodeToGRPCErrorCode)

	chainedUnaryInterceptor := grpc.ChainUnaryInterceptor(
		commonGrpc.TelemetryUnaryServerInterceptor(serviceConfig.Name, serviceConfig.Env, res.Logger()),
	)
	interceptors := append(commonIntercept, chainedUnaryInterceptor)
	server := grpc.NewServer(interceptors...)
	return &transport{
		resource:   res,
		app:        app,
		controller: controller,
		server:     server,
	}

}
