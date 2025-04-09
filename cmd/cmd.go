package cmd

import (
	"github.com/chaihaobo/be-template/application"
	"github.com/chaihaobo/be-template/cmd/cmder"
	"github.com/chaihaobo/be-template/cmd/core"
	"github.com/chaihaobo/be-template/infrastructure"
	"github.com/chaihaobo/be-template/resource"
	"github.com/chaihaobo/be-template/transport"
)

func Execute() error {
	ctx, err := initialContext()
	if err != nil {
		return err
	}
	return cmder.NewRoot().Command(ctx).Execute()
}

func initialContext() (*core.Context, error) {
	res, err := resource.New("./configuration.yaml")
	if err != nil {
		return nil, err
	}

	infra, err := infrastructure.New(res)
	if err != nil {
		return nil, err
	}
	app := application.New(res, infra)
	tsp := transport.New(res, infra, app)
	ctx := core.NewContext(res, infra, tsp)
	return ctx, nil
}
