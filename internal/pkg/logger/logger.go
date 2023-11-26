package logger

import (
	"github.com/hhk7734/zapx.go"
	"go.uber.org/zap"
)

const (
	LOG_LEVEL_KEY = "log_level"
)

type LogConfig struct {
	Level string
}

func SetGlobalZapLogger(cfg LogConfig) {
	var l *zap.Logger
	zapCfg := zap.NewProductionConfig()
	err := zapCfg.Level.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		panic(err)
	}
	zapCfg.EncoderConfig.TimeKey = "time"
	l, _ = zapCfg.Build()
	defer l.Sync()
	zap.ReplaceGlobals(l)

	zap.L().Info("logger config", zapx.Dict("config", zap.String("level", cfg.Level)))
}
