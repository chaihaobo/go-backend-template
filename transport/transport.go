package transport

import (
	"context"

	"github.com/google/wire"

	"github.com/chaihaobo/be-template/application"
	"github.com/chaihaobo/be-template/resource"
	"github.com/chaihaobo/be-template/transport/grpc"
	"github.com/chaihaobo/be-template/transport/http"
)

var ProviderSet = wire.NewSet(
	New,
	http.ProviderSet,
	grpc.ProviderSet,
)

type (
	Transport interface {
		ServeAll() error
		ShutdownAll() error
		HTTP() http.Transport
		Grpc() grpc.Transport
	}

	transport struct {
		res  resource.Resource
		http http.Transport
		grpc grpc.Transport
	}
)

func (t *transport) Grpc() grpc.Transport {
	return t.grpc
}

func (t *transport) ShutdownAll() error {
	return t.http.Shutdown()
}

func (t *transport) ServeAll() error {

	go func() {
		if err := t.Grpc().Serve(); err != nil {
			t.res.Logger().Error(context.Background(), "failed to serve grpc server", err)
		}
	}()

	return t.HTTP().Serve()
}

func (t *transport) HTTP() http.Transport {
	return t.http
}

func New(res resource.Resource, app application.Application, httpTransport http.Transport, grpcTransport grpc.Transport) Transport {
	return &transport{
		res:  res,
		http: httpTransport,
		grpc: grpcTransport,
	}
}
