package middleware

import (
	"gomess/pkg/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(jwt jwt.JWTInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "missing authorization header"},
			)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := jwt.ParseAccessToken(token)

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "invalid token"},
			)
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
