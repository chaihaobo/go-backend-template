package tools

import (
	"log"
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
	log.Println("graceful shutdown successful")

}
