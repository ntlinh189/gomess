package middleware

import (
	"context"
	appcontext "gomess/internal/context"
	"gomess/internal/logger"
	"gomess/internal/redis"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimit(
	rdb redis.RedisInterface,
	prefix string,
	limit int,
	window time.Duration,
	keyFunc func(c *gin.Context) string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "ratelimit:" + prefix + ":" + keyFunc(c)

		count, err := rdb.Incr(context.Background(), key)
		if err != nil {
			// Redis error/cannot connect -> skip
			logger.FromGin(c).Error("rate limit error", "error", err)
			c.Next()
			return
		}

		if count == 1 {
			if err := rdb.Expire(context.Background(), key, window); err != nil {
				logger.FromGin(c).Error("rate limit expire error", "error", err)
			}
		}

		if count > int64(limit) {
			c.Header("Retry-After", window.String())
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RateLimitByIP(c *gin.Context) string {
	return c.ClientIP()
}

func RateLimitByUser(c *gin.Context) string {
	userID := c.GetInt64(appcontext.UserIDKey)
	return strconv.FormatInt(userID, 10)
}