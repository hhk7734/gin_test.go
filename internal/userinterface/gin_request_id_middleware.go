package userinterface

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hhk7734/zapx.go"
	"go.uber.org/zap"
)

type GinRequestIDMiddleware struct{}

// RequestID는 요청에 X-Request-Id 헤더가 있는지 확인하는 미들웨어입니다. generateIfNotExist가
// true고 X-Request-Id 헤더가 없으면 UUID를 생성하여 헤더에 추가합니다. generateIfNotExist가 false고
// X-Request-Id 헤더가 없으면 요청을 거부합니다.
func (r *GinRequestIDMiddleware) RequestID(generateIfNotExist bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-Id")
		if _, err := uuid.Parse(requestID); err != nil {
			if generateIfNotExist {
				requestID = uuid.New().String()
				c.Request.Header.Set("X-Request-Id", requestID)
			} else {
				c.AbortWithStatusJSON(http.StatusBadRequest, RequestIDRequiredResponse)
				return
			}
		}
		// Response
		c.Header("X-Request-Id", requestID)

		ctx := c.Request.Context()
		ctx = zapx.WithFields(ctx, zap.String("request_id", requestID))

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
