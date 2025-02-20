package core

import (
	"github.com/chaihaobo/be-template/infrastructure"
	"github.com/chaihaobo/be-template/resource"
	"github.com/chaihaobo/be-template/transport"
)

type Context struct {
	Resource       resource.Resource
	Infrastructure infrastructure.Infrastructure
	Transport      transport.Transport
}

func NewContext(res resource.Resource, infra infrastructure.Infrastructure, tsp transport.Transport) *Context {
	return &Context{
		Resource:       res,
		Infrastructure: infra,
		Transport:      tsp,
	}
}
