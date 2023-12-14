package logger

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	LOG_LEVEL_KEY = "log_level"
)

type LogConfig struct {
	Level string
}

func LogPFlags() *pflag.FlagSet {
	f := pflag.NewFlagSet("log", pflag.ContinueOnError)
	f.String(LOG_LEVEL_KEY, "info", "log level")
	return f
}

func LogConfigFromViper() LogConfig {
	return LogConfig{
		Level: viper.GetString(LOG_LEVEL_KEY)}
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

	zap.L().Info("logger config", zap.Dict("config",
		zap.String(LOG_LEVEL_KEY, cfg.Level),
	))
}
