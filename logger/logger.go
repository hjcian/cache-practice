package logger

import (
	"go.uber.org/zap"
)

func Init() (*zap.Logger, error) {
	// TODO: use env to switch?
	logger, err := zap.NewDevelopment()
	// logger, err := zap.NewProduction()
	return logger, err
}
