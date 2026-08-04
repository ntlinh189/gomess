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
	Login(providerName, token string) (*dto.LoginResponse, string, error)
	Refresh(refreshToken string) (*dto.RefreshResponse, string, error)
	Logout(refreshToken string) error
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
			"facebook": authprovider.NewFacebookProvider(cfg.GetFacebookAppID(), cfg.GetFacebookAppSecret()),
		},
		jwt:   jwt,
		redis: redis,
	}
}

func (s *Service) Login(providerName, token string) (*dto.LoginResponse, string, error) {
	provider, ok := s.providers[providerName]
	if !ok {
		return nil, "", ErrUnsupportedProvider
	}

	info, err := provider.Verify(token)
	if err != nil {
		return nil, "", err
	}

	u, err := s.repo.FindByProvider(providerName, info.ID)
	if err != nil {
		return nil, "", err
	}

	if u == nil {
		u = &models.User{
			Provider:   providerName,
			ProviderID: info.ID,
			Account:    info.Account,
			Name:       info.Name,
			Avatar:     info.Avatar,
		}
		if err := s.repo.Create(u); err != nil {
			return nil, "", err
		}
	}

	accessToken, err := s.jwt.GenerateAccessToken(u.ID)
	if err != nil {
		return nil, "", err
	}

	refreshToken, err := s.issueRefreshToken(u.ID)
	if err != nil {
		return nil, "", err
	}

	return &dto.LoginResponse{
		AccessToken: accessToken,
	}, refreshToken, nil
}

func (s *Service) Refresh(refreshToken string) (*dto.RefreshResponse, string, error){
	userID, err := s.redis.GetDel(
		context.Background(),
		"refresh:"+utils.SHA256(refreshToken),
	)
	if err != nil {
		if errors.Is(err, redis.ErrNotFound) {
			return nil, "", ErrInvalidRefreshToken
		}
		return nil, "", err
	}

	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, "", ErrInvalidRefreshToken
	}

	accessToken, err := s.jwt.GenerateAccessToken(id)
	if err != nil {
		return nil, "", err
	}

	newRefreshToken, err := s.issueRefreshToken(id)
	if err != nil {
		return nil, "", err
	}

	return &dto.RefreshResponse{
		AccessToken: accessToken,
	}, newRefreshToken, nil
}

func (s *Service) Logout(refreshToken string) error {
	return s.redis.Delete(
		context.Background(),
		"refresh:"+utils.SHA256(refreshToken),
	)
}

func (s *Service) issueRefreshToken(userID int64) (string, error) {
	refreshToken, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return "", err
	}

	err = s.redis.Set(
		context.Background(),
		"refresh:"+utils.SHA256(refreshToken),
		strconv.FormatInt(userID, 10),
		30 * 24 * time.Hour,
	)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}