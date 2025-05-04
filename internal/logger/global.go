package logger

import "go.uber.org/zap"

var log *zap.Logger

func Inject(logger *zap.Logger) {
	log = logger
}

func Get() *zap.Logger {
	if log == nil {
		panic("global logger not initialized")
	}
	return log
}
