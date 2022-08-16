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
