package logger

import (
	"github.com/chaihaobo/gocommon/logger"
)

type Logger logger.Logger

// New create new instant for the Logger.
func New(config Config) (Logger, func() error, error) {
	return logger.New(logger.Config(config))
}
