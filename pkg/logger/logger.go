package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey struct{}

var defaultLogger *zap.Logger

type Config struct {
	Level       string
	Development bool
}

func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return logger
	}
	return defaultLogger
}

func ToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Sugar().Infof(format, args...)
}

func Info(ctx context.Context, msg string) {
	FromContext(ctx).Sugar().Infow(msg)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Sugar().Errorf(format, args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Sugar().Fatalf(format, args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Sugar().Debugf(format, args...)
}

func Init(config Config) {
	var err error
	var zapConfig zap.Config

	if config.Development {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}

	switch config.Level {
	case "debug":
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "error":
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	defaultLogger, err = zapConfig.Build()
	if err != nil {
		panic(err)
	}

	defer defaultLogger.Sync()
}

func GetLogger() *zap.Logger {
	return defaultLogger
}
