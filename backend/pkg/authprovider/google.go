package authprovider

import (
	"context"
	"errors"

	"google.golang.org/api/idtoken"
)

type GoogleProvider struct {
	clientID string
}

func NewGoogleProvider(clientID string) *GoogleProvider {
	return &GoogleProvider{clientID: clientID}
}

func (g *GoogleProvider) Verify(token string) (*UserInfo, error) {
	payload, err := idtoken.Validate(context.Background(), token, g.clientID)
	if err != nil {
		return nil, err
	}

	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)

	if email == "" {
		return nil, errors.New("Invalid token")
	}

	return &UserInfo{
		ID:     payload.Subject,
		Email:  email,
		Name:   name,
		Avatar: picture,
	}, nil
}
