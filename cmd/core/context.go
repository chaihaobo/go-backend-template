package core

import (
	"gitlab.seakoi.net/engineer/backend/be-template/infrastructure"
	"gitlab.seakoi.net/engineer/backend/be-template/resource"
	"gitlab.seakoi.net/engineer/backend/be-template/transport"
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
