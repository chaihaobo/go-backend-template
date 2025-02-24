package logger

import (
	"github.com/chaihaobo/gocommon/logger"
)

type Logger logger.Logger

// New create new instant for the Logger.
func New(config Config) (Logger, func() error, error) {
	log, flusher, err := logger.New(logger.Config{
		FileName:   config.FileName,
		MaxSize:    config.MaxSize,
		MaxAge:     config.MaxAge,
		WithCaller: config.WithCaller,
		CallerSkip: config.CallerSkip,
	})
	if err != nil {
		return nil, nil, err
	}
	logger.SetSlogDefaultFromLogger(log)
	return log, flusher, nil
}
