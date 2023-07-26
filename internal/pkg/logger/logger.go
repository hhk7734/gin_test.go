package logger

import (
	"os"

	"github.com/hhk7734/zapx.go"
	"go.uber.org/zap"
)

func SetGlobalZapLogger() {
	level, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		level = "info"
	}

	var l *zap.Logger
	zapCfg := zap.NewProductionConfig()
	err := zapCfg.Level.UnmarshalText([]byte(level))
	if err != nil {
		panic(err)
	}
	zapCfg.EncoderConfig.TimeKey = "time"
	l, _ = zapCfg.Build()
	defer l.Sync()
	zap.ReplaceGlobals(l)

	zap.L().Info("logger config", zapx.Dict("config", zap.String("level", level)))
}
