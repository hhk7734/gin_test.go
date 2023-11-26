package main

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hhk7734/gin-test/internal/pkg/logger"
	"github.com/hhk7734/gin-test/internal/userinterface/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	viper.SetConfigName(".env")
	viper.SetConfigType("dotenv")

	// TODO: find .env file from parent directory recursively
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")

	viper.AutomaticEnv()

	pflag.String("log_level", "info", "log level")
	viper.BindPFlags(pflag.CommandLine)

	pflag.Parse()

	logger.SetGlobalZapLogger(logger.LogConfig{Level: viper.GetString("log_level")})

	server := gin.NewGinRestAPI()

	listenErr := make(chan error, 1)
	go func() {
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			listenErr <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-listenErr:
		zap.L().Error("failed to listen and serve", zap.Error(err))
	case <-shutdown:
	}

	zap.L().Info("shutting down server...")

	wg := &sync.WaitGroup{}

	go func() {
		defer wg.Done()
		// blocked until all connections are closed or timeout
		if err := server.Shutdown(); err != nil {
			zap.L().Error("failed to shutdown server", zap.Error(err))
		}
	}()
	wg.Add(1)

	wg.Wait()
}
