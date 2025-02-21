//go:build wireinject
// +build wireinject

package cmd

import (
	"github.com/google/wire"

	"github.com/chaihaobo/be-template/application"
	"github.com/chaihaobo/be-template/cmd/core"
	"github.com/chaihaobo/be-template/infrastructure"
	"github.com/chaihaobo/be-template/resource"
	"github.com/chaihaobo/be-template/transport"
)

func initContext(configPath string) (*core.Context, error) {
	panic(
		wire.Build(
			resource.ProviderSet,
			infrastructure.ProviderSet,
			application.ProviderSet,
			transport.ProviderSet,
			core.NewContext,
		),
	)
}
