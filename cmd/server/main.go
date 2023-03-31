package main

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hhk7734/gin-test/internal/user_interface/restapi/middleware"
	"go.uber.org/zap"
)

func main() {
	workDir, _ := os.Getwd()
	for {
		if _, err := os.Stat(workDir + "/.env"); err == nil {
			os.Chdir(workDir)
			break
		}
		if workDir == "/" {
			break
		}
		workDir = filepath.Dir(workDir)
	}

	r := gin.New()
	r.Use(middleware.LoggerWithZap([]string{}))
	r.Use(middleware.RecoveryWithZap())

	r.StaticFile("/openapi.yaml", "web/static/openapi.yaml")

	s := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	if err := s.ListenAndServe(); err != nil {
		zap.L().Panic("Server error", zap.Error(err))
	}
}
