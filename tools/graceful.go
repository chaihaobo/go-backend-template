package tools

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/samber/lo"
)

func GracefulShutdown(actions ...func() error) {
	graceful := make(chan os.Signal, 1)
	signal.Notify(graceful, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-graceful
	lo.ForEach(actions, func(item func() error, _ int) {
		if err := item(); err != nil {
			panic(err)
		}
	})
	slog.Error("graceful shutdown successful")

}
