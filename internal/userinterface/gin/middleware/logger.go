package middleware

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

func (l *GinLoggerMiddleware) Logger(skipPaths []string) gin.HandlerFunc {
	skip := make(map[string]bool, len(skipPaths))
	for _, path := range skipPaths {
		skip[path] = true
	}

	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path

		c.Next()

		if _, ok := skip[path]; !ok {
			status := c.Writer.Status()
			ctx := zapx.WithFields(c.Request.Context(),
				zap.String("method", c.Request.Method),
				zap.String("url", path),
				zap.Int("status", status),
				zap.String("remote_address", c.ClientIP()),
				zap.String("user_agent", c.Request.UserAgent()),
				zap.Duration("latency", time.Since(start)))
			logger := zapx.Ctx(ctx)

			if status < http.StatusInternalServerError {
				logger.Info("request")
			} else {
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				logger.Error("request", zap.String("request", string(httpRequest)))
			}

			for _, e := range c.Errors {
				logger.Error("unknown error", zap.Error(e))
			}
		}

	}
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
				zapx.Ctx(c.Request.Context()).Warn("broken pipe", zap.Error(err.(error)))
				c.Abort()
				return
			}

			zapx.Ctx(c.Request.Context()).Error("internal server error", zap.Error(fmt.Errorf("%v\n%s", err, string(debug.Stack()))))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()

	c.Next()
}
