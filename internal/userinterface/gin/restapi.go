package gin

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/hhk7734/gin_test.go/internal/pkg/validator"
	"github.com/hhk7734/gin_test.go/internal/userinterface/gin/controller"
	"github.com/hhk7734/gin_test.go/internal/userinterface/gin/middleware"
)

func init() {
	binding.Validator = &validator.GinValidator{}
}

type GinRestAPI struct {
	engin  *gin.Engine
	server *http.Server
}

func NewGinRestAPI() *GinRestAPI {
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

	return &GinRestAPI{
		engin:  engin,
		server: server,
	}
}

func (g *GinRestAPI) Run() error {
	return g.server.ListenAndServe()
}

func (g *GinRestAPI) Shutdown() error {
	return g.server.Shutdown(context.Background())
}
