package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hhk7734/ratelimit.go"
	"github.com/hhk7734/zapx.go"
	"go.uber.org/zap"
)

func NewGinRateLimitMiddleware() *GinRateLimitMiddleware {
	mr := ratelimit.NewMemoryRateLimit()
	hr := ratelimit.NewHttpRateLimit(mr)
	return &GinRateLimitMiddleware{ratelimit: hr}
}

type GinRateLimitMiddleware struct {
	ratelimit *ratelimit.HttpRateLimit
}

func (r *GinRateLimitMiddleware) IPRateLimit(key string, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ip := c.ClientIP()
		fmt.Println(ip)
		k := fmt.Sprintf("%s:%s", ip, key)
		err := r.ratelimit.SlidingWindowLog(ctx, c.Writer, k, limit, window)
		switch {
		case errors.Is(err, ratelimit.ErrLimitExceeded):
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		case err != nil:
			zapx.Ctx(ctx).Error("ip rate limit", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Next()
	}
}
