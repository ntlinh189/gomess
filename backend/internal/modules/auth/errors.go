package auth

import "errors"

var (
	ErrUnsupportedProvider = errors.New("unsupported provider")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
)