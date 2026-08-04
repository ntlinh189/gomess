package middleware

import (
	"gomess/internal/context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(context.RequestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set(context.RequestIDKey, requestID)
		c.Header(context.RequestIDHeader, requestID)

		c.Next()
	}
}