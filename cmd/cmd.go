package cmd

import (
	"github.com/chaihaobo/be-template/cmd/cmder"
)

func Execute() error {
	ctx, err := initContext("configuration.yaml")
	if err != nil {
		return err
	}
	return cmder.NewRoot().Command(ctx).Execute()
}
