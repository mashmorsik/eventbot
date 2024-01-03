package Logger

import (
	"go.uber.org/zap"
	"log/slog"
)

var Sugar *zap.SugaredLogger

func InitSugarLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	Sugar = logger.Sugar()

	slog.Info("Logger initialized")
}
