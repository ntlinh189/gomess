package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTInterface interface {
	GenerateAccessToken(userID int64) (string, error)
}

type JWT struct {
	secret string
}

func NewJWT(secret string) *JWT {
	return &JWT{secret: secret}
}

func (j *JWT) GenerateAccessToken(userID int64) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(5 * time.Minute).Unix(),
		})
	return token.SignedString([]byte(j.secret))
}
