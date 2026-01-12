package logger

import (
	"go.uber.org/zap"
)

func New(level string) *zap.Logger {
	var (
		log *zap.Logger
		err error
	)

	if level == "prod" {
		log, err = zap.NewProduction(zap.AddStacktrace(zap.PanicLevel))
	} else {
		log, err = zap.NewDevelopment(zap.AddStacktrace(zap.PanicLevel))
	}

	if err != nil {
		log = zap.NewExample()
	}

	return log
}
