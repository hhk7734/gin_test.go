package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hhk7734/gin-test/internal/pkg/logger"
	"go.uber.org/zap"
)

type RequestIDMiddleware struct{}

func (r *RequestIDMiddleware) RequestID(c *gin.Context) {
	requestID := c.GetHeader("X-Request-Id")
	if _, err := uuid.Parse(requestID); err != nil {
		requestID = uuid.New().String()
		c.Request.Header.Set("X-Request-Id", requestID)
	}
	// Response
	c.Header("X-Request-Id", requestID)

	ctx := c.Request.Context()
	ctx = logger.WithFields(ctx, zap.String("request_id", requestID))

	c.Request = c.Request.WithContext(ctx)

	c.Next()
}
