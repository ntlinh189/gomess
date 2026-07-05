package user

import (
	"gomess/internal/modules/user/dto"
	"gomess/utils"
)

type ServiceInterface interface {
	GetMe(userID int64) (*dto.GetMeResponse, error)
	Search(req *dto.SearchRequest) ([]dto.SearchResponse, error)
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

func (s *Service) Search(req *dto.SearchRequest) ([]dto.SearchResponse, error) {
	if req.Skip < 0 {
		req.Skip = 0
	}

	if req.Limit <= 0 {
		req.Limit = 20
	}

	if req.Limit > 100 {
		req.Limit = 100
	}

	users, err := s.repo.Search(
		req.Provider,
		utils.BuildLikePattern(req.Keyword),
		req.Skip,
		req.Limit,
	)

	if err != nil {
		return nil, err
	}

	resp := make([]dto.SearchResponse, 0, len(users))

	for _, user := range users {

		resp = append(resp, dto.SearchResponse{
			ID:       user.ID,
			Provider: user.Provider,
			Email:    user.Email,
			Name:     user.Name,
			Avatar:   user.Avatar,
		})
	}

	return resp, nil
}
