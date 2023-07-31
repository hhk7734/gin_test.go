package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hhk7734/gin-test/internal/pkg/env"
	"github.com/hhk7734/gin-test/internal/pkg/logger"
	"github.com/hhk7734/gin-test/internal/userinterface/gin/controller"
	"github.com/hhk7734/gin-test/internal/userinterface/gin/middleware"
	"github.com/hhk7734/zapx.go"
	"go.uber.org/zap"
)

func main() {
	env.Load(".env")
	logger.SetGlobalZapLogger()

	ctx := context.Background()

	lm := &middleware.GinLoggerMiddleware{}
	ratelimit := middleware.NewGinRateLimitMiddleware()

	engin := gin.New()
	engin.RemoteIPHeaders = append([]string{"X-Envoy-External-Address"}, engin.RemoteIPHeaders...)
	engin.Use(lm.Logger([]string{"/healthz"}))
	engin.Use(lm.Recovery)
	engin.Use(middleware.GinRequestIDMiddleware(true))

	engin.GET("/healthz",
		ratelimit.IPRateLimit("healthz", 20, 10*time.Second),
		controller.GinHealthzController)

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
		zapx.Ctx(ctx).Error("failed to listen and serve", zap.Error(err))
	case <-shutdown:
	}

	zapx.Ctx(ctx).Info("shutting down server...")

	wg := &sync.WaitGroup{}

	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()

		// blocked until all connections are closed or timeout
		if err := server.Shutdown(ctx); err != nil {
			zapx.Ctx(ctx).Error("failed to shutdown server", zap.Error(err))
		}
	}()
	wg.Add(1)

	wg.Wait()
}
