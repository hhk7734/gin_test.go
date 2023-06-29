package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hhk7734/gin-test/internal/interface/gin/controller"
	"github.com/hhk7734/gin-test/internal/interface/gin/middleware"
	"github.com/hhk7734/gin-test/internal/pkg/logger"
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

	lm := &logger.GinLoggerMiddleware{}

	engin := gin.New()
	engin.Use(lm.Logger)
	engin.Use(lm.Recovery)
	engin.Use((&middleware.RequestIDMiddleware{}).RequestID)

	engin.GET("/healthz", (&controller.HealthzController{}).Healthz)

	engin.StaticFile("/openapi.yaml", "web/static/openapi.yaml")

	server := &http.Server{
		Addr:         ":8080",
		Handler:      engin,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	listenErr := make(chan error, 1)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			listenErr <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-listenErr:
		logger.Logger(context.Background()).Error("failed to listen and serve", zap.Error(err))
	case <-shutdown:
	}

	logger.Logger(context.Background()).Info("shutting down server...")

	wg := &sync.WaitGroup{}

	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		// blocked until all connections are closed or timeout
		if err := server.Shutdown(ctx); err != nil {
			logger.Logger(context.Background()).Error("failed to shutdown server", zap.Error(err))
		}
	}()
	wg.Add(1)

	wg.Wait()
}
