package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func Init() {
	Logger, _ = zap.NewProduction()
	defer Logger.Sync()
}