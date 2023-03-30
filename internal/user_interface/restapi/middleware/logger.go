package middleware

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hhk7734/gin-test/internal/pkg/logger"
	"go.uber.org/zap"
)

const UserIDKey = "user_id"

func LoggerWithZap(skipPaths []string) gin.HandlerFunc {
	skip := make(map[string]bool, len(skipPaths))
	for _, path := range skipPaths {
		skip[path] = true
	}

	return func(c *gin.Context) {
		start := time.Now()
		// 일부 미들웨어는 경로를 중간에 바꾸는 경우가 있음
		path := c.Request.URL.Path
		c.Next()

		if _, ok := skip[path]; !ok {
			end := time.Now()
			latency := end.Sub(start)
			userID := func() int64 {
				if userID, ok := c.Get(UserIDKey); ok {
					return int64(userID.(float64))
				} else {
					return -1
				}
			}()
			xRequestID := c.Request.Header.Get("X-Request-Id")

			l := logger.Logger(c.Request.Context()).With(
				zap.String("method", c.Request.Method),
				zap.String("url", path),
				zap.Int("status", c.Writer.Status()),
				zap.Int64("user_id", userID),
				zap.String("request_id", xRequestID),
				zap.String("remote_address", c.ClientIP()),
				zap.String("user_agent", c.Request.UserAgent()),
				zap.Duration("latency", latency),
			)

			if len(c.Errors) > 0 {
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				c.Error(fmt.Errorf("%s", httpRequest))

				// 컨텍스트에 저장된 에러
				for i, e := range c.Errors {
					l.Error(strconv.Itoa(i), zap.Error(e))
				}
			} else {
				// 정상 처리
				l.Info(path)
			}
		}
	}
}

func RecoveryWithZap() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 연결이 끊겼는지 확인
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				if brokenPipe {
					// 연결이 끊겼다면 Status를 설정할 수 없음
					c.Error(err.(error))
					c.Abort()
					return
				}

				c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("%v\n%s", err, string(debug.Stack())))
			}
		}()
		c.Next()
	}
}
