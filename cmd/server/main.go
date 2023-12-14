package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/hhk7734/gin-test/internal/pkg/logger"
	"github.com/hhk7734/gin-test/internal/userinterface/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// .env file
	viper.SetConfigName(".env")
	viper.SetConfigType("dotenv")

	workDir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(workDir, ".env")); err == nil {
			viper.AddConfigPath(workDir)
			break
		}
		if workDir == "/" {
			break
		}
		workDir = filepath.Dir(workDir)
	}

	if err := viper.ReadInConfig(); err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}

	// env
	viper.AutomaticEnv()

	// flag
	pflag.CommandLine.AddFlagSet(logger.LogPFlags())
	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)

	logger.SetGlobalZapLogger(logger.LogConfigFromViper())

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
