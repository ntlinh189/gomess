package middleware

import (
	"gomess/internal/context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Recovery(log *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		requestID, _ := c.Get(context.RequestIDKey)

		log.Error(
			"panic recovered",
			"request_id", requestID,
			"panic", recovered,
			"path", c.Request.URL.Path,
		)

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	})
}
