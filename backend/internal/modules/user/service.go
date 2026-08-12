package user

import (
	"gomess/internal/modules/user/dto"
	"gomess/utils"
)

type ServiceInterface interface {
	GetMe(userID int64) (*dto.GetMeResponse, error)
	Search(userID int64, req *dto.SearchRequest) ([]dto.SearchResponse, error)
	DeleteMe(userID int64) error
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

	if user == nil {
		return nil, ErrUserNotFound
	}

	return &dto.GetMeResponse{
		ID:      user.ID,
		Account: user.Account,
		Name:    user.Name,
		Avatar:  user.Avatar,
	}, nil
}

func (s *Service) Search(userID int64, req *dto.SearchRequest) ([]dto.SearchResponse, error) {
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
		userID,
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
			ID:            user.ID,
			Provider:      user.Provider,
			Account:       user.Account,
			Name:          user.Name,
			Avatar:        user.Avatar,
			RequestStatus: user.RequestStatus,
		})
	}

	return resp, nil
}

func (s *Service) DeleteMe(userID int64) error {
	return s.repo.DeleteUser(userID)
}
