package jwt

import (
	"errors"
	"gomess/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTInterface interface {
	GenerateAccessToken(userID int64) (string, error)
	ParseAccessToken(tokenString string) (int64, error)
	GenerateRefreshToken() (string, error)
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

func (j *JWT) ParseAccessToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (any, error) {
			return []byte(j.secret), nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, ErrExpiredToken
		}
		return 0, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, ErrInvalidToken
	}

	rawUserID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, ErrInvalidToken
	}

	return int64(rawUserID), nil
}

func (j *JWT) GenerateRefreshToken() (string, error) {
	return utils.GenerateToken(32)
}
