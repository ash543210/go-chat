package middleware

import (
	"github.com/ash543210/go-chat/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LoggerMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate a new trace ID for this request
		traceID := uuid.NewString()

		// Add trace ID to context (your logger's context helper)
		ctx := logger.WithTraceID(c.Request.Context(), traceID)
		c.Request = c.Request.WithContext(ctx)

		// Optionally, store logger in Gin context for other handlers
		c.Set("logger", log)

		// You can also expose traceID for debugging/headers
		c.Writer.Header().Set(string(logger.TraceKey), traceID)

		c.Next()
	}
}
