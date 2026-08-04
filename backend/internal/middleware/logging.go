package middleware

import (
	"gomess/internal/context"
	"gomess/internal/logger"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func WithLogger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, _ := c.Get(context.RequestIDKey)
		logger.Inject(c, log.With("request_id", requestID))
		c.Next()
	}
}

func Logging(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		if raw := c.Request.URL.RawQuery; raw != "" {
			path += "?" + raw
		}

		c.Next()

		requestID, _ := c.Get(context.RequestIDKey)

		attrs := []any{
			"request_id", requestID,
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs, "errors", c.Errors.String())
		}

		switch {
		case c.Writer.Status() >= 500:
			log.Error("request completed", attrs...)
		case c.Writer.Status() >= 400:
			log.Warn("request completed", attrs...)
		default:
			log.Info("request completed", attrs...)
		}
	}
}
