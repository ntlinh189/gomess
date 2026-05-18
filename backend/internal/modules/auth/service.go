package auth

import (
	"context"
	"errors"
	"gomess/internal/config"
	"gomess/internal/models"
	"gomess/internal/modules/auth/dto"
	"gomess/internal/redis"
	"gomess/pkg/authprovider"
	"gomess/pkg/jwt"
	"gomess/utils"
	"strconv"
	"time"
)

type ServiceInterface interface {
	Login(providerName, token string) (*dto.LoginResponse, error)
}

type Service struct {
	repo      RepositoryInterface
	providers map[string]authprovider.ProviderInterface
	jwt       jwt.JWTInterface
	redis     redis.RedisInterface
}

func NewService(repo RepositoryInterface, jwt jwt.JWTInterface, cfg config.ConfigInterface, redis redis.RedisInterface) *Service {
	return &Service{
		repo: repo,
		providers: map[string]authprovider.ProviderInterface{
			"google": authprovider.NewGoogleProvider(cfg.GetGoogleClientID()),
		},
		jwt:   jwt,
		redis: redis,
	}
}

func (s *Service) Login(providerName, token string) (*dto.LoginResponse, error) {
	provider, ok := s.providers[providerName]
	if !ok {
		return nil, errors.New("Unsupported provider")
	}

	info, err := provider.Verify(token)
	if err != nil {
		return nil, err
	}

	u, err := s.repo.FindByProvider(providerName, info.ID)
	if err != nil {
		return nil, err
	}

	if u == nil {
		u = &models.User{
			Provider:   providerName,
			ProviderID: info.ID,
			Email:      info.Email,
			Name:       info.Name,
			Avatar:     info.Avatar,
		}
		if err := s.repo.Create(u); err != nil {
			return nil, err
		}
	}

	accessToken, err := s.jwt.GenerateAccessToken(u.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateToken(32)
	if err != nil {
		return nil, err
	}

	err = s.redis.Set(
		context.Background(),
		"refresh:"+utils.SHA256(refreshToken),
		strconv.FormatInt(u.ID, 10),
		30*24*time.Hour,
	)

	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
