package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/maps"
)

type filedsMapKey struct{}

var zConfig zap.Config

func init() {
	var l *zap.Logger
	zConfig = zap.NewProductionConfig()
	zConfig.EncoderConfig.TimeKey = "time"
	zConfig.DisableStacktrace = true
	l, _ = zConfig.Build()
	defer l.Sync()
	zap.ReplaceGlobals(l)
}

func SetLevel(level zapcore.Level) {
	zConfig.Level.SetLevel(level)
}

func Logger(ctx context.Context) *zap.Logger {
	fs := Fields(ctx)
	if fs == nil {
		return zap.L()
	}

	return zap.L().With(fs...)
}

func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	v := ctx.Value(filedsMapKey{})
	if v == nil {
		fm := make(map[string]zap.Field, len(fields))
		for _, f := range fields {
			fm[f.Key] = f
		}
		return context.WithValue(ctx, filedsMapKey{}, fm)
	}

	fm := v.(map[string]zap.Field)
	for _, f := range fields {
		fm[f.Key] = f
	}

	return ctx
}

func Fields(ctx context.Context) []zap.Field {
	v := ctx.Value(filedsMapKey{})
	if v == nil {
		return nil
	}

	return maps.Values(v.(map[string]zap.Field))
}