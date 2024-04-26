package cmd

import (
	"gitlab.seakoi.net/engineer/backend/be-template/application"
	"gitlab.seakoi.net/engineer/backend/be-template/cmd/cmder"
	"gitlab.seakoi.net/engineer/backend/be-template/cmd/core"
	"gitlab.seakoi.net/engineer/backend/be-template/infrastructure"
	"gitlab.seakoi.net/engineer/backend/be-template/resource"
	"gitlab.seakoi.net/engineer/backend/be-template/transport"
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
	tsp := transport.New(res, app)
	ctx := core.NewContext(res, infra, tsp)
	return ctx, nil
}
