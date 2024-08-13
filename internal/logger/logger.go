package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey struct{}

var defaultLogger *zap.Logger

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

func Init() {
	var err error
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	defaultLogger, err = config.Build()
	if err != nil {
		panic(err)
	}

	defer defaultLogger.Sync()
}

func GetLogger() *zap.Logger {
	return defaultLogger
}
