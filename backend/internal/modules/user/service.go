package user

import (
	"gomess/internal/modules/user/dto"
)

type ServiceInterface interface {
	GetMe(userID int64) (*dto.GetMeResponse, error)
}

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetMe(userID int64) (*dto.GetMeResponse, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return &dto.GetMeResponse{
		ID:     user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Avatar: user.Avatar,
	}, nil
}