package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hhk7734/gin-test/internal/config"
	"github.com/hhk7734/gin-test/internal/user_interface/gin/middleware"
	"go.uber.org/zap"
)

func main() {
	config.Init()
	c := config.Config()

	initLogger()

	r := gin.New()
	r.Use(middleware.LoggerWithZap([]string{}))
	r.Use(middleware.RecoveryWithZap())

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Port),
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	if err := s.ListenAndServe(); err != nil {
		zap.L().Panic("Server error", zap.Error(err))
	}
}

func initLogger() {
	var l *zap.Logger
	c := config.Config()
	if c.Debug {
		cfg := zap.NewDevelopmentConfig()
		cfg.DisableStacktrace = true
		l, _ = cfg.Build()
	} else {
		cfg := zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "time"
		cfg.DisableStacktrace = true
		l, _ = cfg.Build()
	}
	defer l.Sync()
	zap.ReplaceGlobals(l)
}
