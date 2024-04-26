package core

import "github.com/spf13/cobra"

type (
	Cmder interface {
		Command(ctx *Context) *cobra.Command
	}
	CmderFunc func(ctx *Context) *cobra.Command
)

func (c CmderFunc) Command(ctx *Context) *cobra.Command {
	return c(ctx)
}
