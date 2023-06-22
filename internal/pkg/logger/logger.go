package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/hhk7734/gin-test/internal/pkg/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/maps"
)

func init() {
	if err := SetConfig(NewConfigFromENV()); err != nil {
		panic(err)
	}
}

type filedsMapKey struct{}

var zConfig = zap.NewProductionConfig()

func NewConfigFromENV() Config {
	if err := env.Load(".env"); err != nil {
		fmt.Println(`{"level":"warn","msg":"` + err.Error() + `"}`)
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	return cfg
}

type Config struct {
	Format string `env:"LOG_FORMAT" envDefault:"json"`
	Level  string `env:"LOG_LEVEL" envDefault:"info"`
}

func SetConfig(c Config) error {
	if err := SetLevel(c.Level); err != nil {
		return err
	}
	if err := SetFormat(c.Format); err != nil {
		return err
	}
	return nil
}

func SetFormat(format string) error {
	zConfig.EncoderConfig.TimeKey = "time"
	zConfig.DisableStacktrace = true
	switch format {
	case "text":
		zConfig.Encoding = "console"
		zConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.999Z0700"))
		})
	case "json":
		zConfig.Encoding = "json"
		zConfig.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		zConfig.EncoderConfig.EncodeTime = zapcore.EpochTimeEncoder
	default:
		return fmt.Errorf("`%s` is not supported log format, use `text` or `json`", format)
	}

	var l *zap.Logger
	l, _ = zConfig.Build()
	defer l.Sync()
	zap.ReplaceGlobals(l)
	return nil
}

func SetLevel(level string) error {
	l, err := zapcore.ParseLevel(level)
	if err != nil {
		return err
	}
	zConfig.Level.SetLevel(l)
	return nil
}

func Logger(ctx context.Context) *zap.Logger {
	fs := Fields(ctx)
	if fs == nil {
		return zap.L()
	}

	return zap.L().With(fs...)
}

func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	fieldMap := make(map[string]zap.Field, len(fields))

	if fm := ctx.Value(filedsMapKey{}); fm != nil {
		for k, v := range fm.(map[string]zap.Field) {
			fieldMap[k] = v
		}
	}

	for _, f := range fields {
		fieldMap[f.Key] = f
	}

	return context.WithValue(ctx, filedsMapKey{}, fieldMap)
}

func Fields(ctx context.Context) []zap.Field {
	v := ctx.Value(filedsMapKey{})
	if v == nil {
		return nil
	}

	return maps.Values(v.(map[string]zap.Field))
}
