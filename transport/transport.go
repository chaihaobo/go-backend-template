package transport

import (
	"gitlab.seakoi.net/engineer/backend/be-template/application"
	"gitlab.seakoi.net/engineer/backend/be-template/resource"
	"gitlab.seakoi.net/engineer/backend/be-template/transport/http"
)

type (
	Transport interface {
		ServeAll() error
		ShutdownAll() error
		HTTP() http.Transport
	}

	transport struct {
		http http.Transport
	}
)

func (t *transport) ShutdownAll() error {
	return t.http.Shutdown()
}

func (t *transport) ServeAll() error {
	return t.HTTP().Serve()
}

func (t *transport) HTTP() http.Transport {
	return t.http
}

func New(res resource.Resource, app application.Application) Transport {
	httpTransport := http.NewTransport(res, app)
	return &transport{
		http: httpTransport,
	}
}
