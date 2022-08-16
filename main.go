package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hhk7734/git-test/internal/config"
	"go.uber.org/zap"
)

// @title       Gin test
// @version     1.0
// @description Gin test

// @schemes  http
// @host     localhost:8080
// @BasePath /
func main() {
	config.Init()
	c := config.Config()

	initLogger()

	r := gin.New()
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
