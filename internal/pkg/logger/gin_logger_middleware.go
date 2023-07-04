package logger

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hhk7734/zapx.go"
	"go.uber.org/zap"
)

type GinLoggerMiddleware struct{}

func (l *GinLoggerMiddleware) Logger(c *gin.Context) {
	start := time.Now()
	// some evil middlewares modify this values
	path := c.Request.URL.Path

	c.Next()

	ctx := zapx.WithFields(c.Request.Context(),
		zap.String("method", c.Request.Method),
		zap.String("url", path),
		zap.Int("status", c.Writer.Status()),
		zap.String("remote_address", c.ClientIP()),
		zap.String("user_agent", c.Request.UserAgent()),
		zap.Duration("latency", time.Since(start)))

	log := zapx.Ctx(ctx)

	if c.Writer.Status() >= http.StatusInternalServerError || len(c.Errors) > 0 {
		httpRequest, _ := httputil.DumpRequest(c.Request, false)
		log.Error("internal server error", zap.String("request", string(httpRequest)))
		for _, e := range c.Errors {
			log.Error("internal server error", zap.Error(e), zap.Any("meta", e.Meta))
		}
		return
	}

	log.Info("request")
}

func (l *GinLoggerMiddleware) Recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			var brokenPipe bool
			if ne, ok := err.(*net.OpError); ok {
				if se, ok := ne.Err.(*os.SyscallError); ok {
					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						brokenPipe = true
					}
				}
			}

			if brokenPipe {
				c.Error(err.(error))
				c.Abort()
				return
			}

			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("%v\n%s", err, string(debug.Stack())))
		}
	}()

	c.Next()
}
