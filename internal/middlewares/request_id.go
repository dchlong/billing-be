package middlewares

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dchlong/billing-be/pkg/logger"
)

func SetRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			id, _ := uuid.NewUUID()
			requestID = id.String()
		}

		baseCtx := c.Request.Context()
		baseCtx = context.WithValue(baseCtx, logger.RequestIDContextKey{}, requestID)
		c.Request = c.Request.WithContext(baseCtx)
		c.Writer.Header().Set("X-Request-ID", requestID) // Set X-Request-ID header
	}
}
