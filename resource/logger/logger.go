package logger

import (
	"github.com/chaihaobo/gocommon/logger"

	"github.com/chaihaobo/be-template/resource/config"
)

type Logger logger.Logger

// New create new instant for the Logger.
func New(config *config.Configuration) (Logger, error) {
	logConfig := config.Logger
	l, _, err := logger.New(logger.Config{
		FileName: logConfig.FileName,
		MaxSize:  logConfig.MaxSize,
		MaxAge:   logConfig.MaxAge,
	})
	return l, err
}
