package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level, env string) *zap.Logger {
	var cfg zap.Config

	if env == "production" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	l, err := zapcore.ParseLevel(level)
	if err == nil {
		cfg.Level.SetLevel(l)
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger
}
