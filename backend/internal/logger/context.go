package logger

import (
	"gomess/internal/context"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func Inject(c *gin.Context, log *slog.Logger) {
	c.Set(context.LoggerKey, log)
}

func FromGin(c *gin.Context) *slog.Logger {
	if v, ok := c.Get(context.LoggerKey); ok {
		if l, ok := v.(*slog.Logger); ok {
			return l
		}
	}
	return slog.Default()
}