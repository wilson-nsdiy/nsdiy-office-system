package middleware

import (
	"strings"

	"oa-nsdiy/backend/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const requestIDHeader = "X-Request-ID"

// RequestLogger injects a request-scoped logger with request_id into the context.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request == nil {
			c.Next()
			return
		}

		requestID := strings.TrimSpace(c.GetHeader(requestIDHeader))
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Header(requestIDHeader, requestID)

		reqLogger := logger.With(
			zap.String("component", "http"),
			zap.String("request_id", requestID),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)

		ctx := logger.IntoContext(c.Request.Context(), reqLogger)
		c.Request = c.Request.WithContext(ctx)
		c.Set("request_id", requestID)

		c.Next()
	}
}
