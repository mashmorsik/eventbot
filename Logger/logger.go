package Logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log/slog"
	"time"
)

var Sugar *zap.SugaredLogger

func InitSugarLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	logger, err = loggerConfig.Build()
	if err != nil {
		panic(err)
	}

	Sugar = logger.Sugar()

	slog.Info("Logger initialized")
}
