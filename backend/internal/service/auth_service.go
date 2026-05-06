package service

import (
	"errors"
	"gomess/internal/model"
	"gomess/internal/repository"
	"gomess/pkg/authprovider"
	"gomess/pkg/jwt"
)

type AuthService struct {
	repo      repository.UserRepository
	providers map[string]authprovider.Provider
	jwt       *jwt.JWT
}

func NewAuthService(
	repo repository.UserRepository,
	providers map[string]authprovider.Provider,
	jwt *jwt.JWT,
) *AuthService {
	return &AuthService{
		repo:      repo,
		providers: providers,
		jwt:       jwt,
	}
}

func (s *AuthService) Login(providerName, token string) (string, error) {
	provider, ok := s.providers[providerName]

	if !ok {
		return "", errors.New("Unsupported provider")
	}

	info, err := provider.Verify(token)

	if err != nil {
		return "", err
	}

	user, _ := s.repo.FindByProvider(providerName, info.ID)

	if user == nil {
		user = &model.User{
			Provider:   providerName,
			ProviderID: info.ID,
			Email:      info.Email,
			Name:       info.Name,
			Avatar:     info.Avatar,
		}
		s.repo.Create(user)
	}

	return s.jwt.Generate(user.ID)
}
