package middleware

import (
	"errors"
	"gomess/internal/context"
	"gomess/pkg/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(j jwt.JWTInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")

		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header is required",
				"code":  "TOKEN_MISSING",
			})
			c.Abort()
			return
		}

		if !strings.HasPrefix(auth, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
				"code":  "TOKEN_INVALID",
			})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")

		userID, err := j.ParseAccessToken(token)
		if err != nil {
			if errors.Is(err, jwt.ErrExpiredToken) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "access token expired",
					"code":  "TOKEN_EXPIRED",
				})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "invalid token",
					"code":  "TOKEN_INVALID",
				})
			}
			c.Abort()
			return
		}

		c.Set(context.UserIDKey, userID)

		c.Next()
	}
}