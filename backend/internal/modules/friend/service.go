package friend

import (
	"database/sql"
	"gomess/internal/database"
	"gomess/internal/models"
)

type UserRepositoryInterface interface {
	ExistsUser(userID int64) (bool, error)
}

type MessageRepositoryInterface interface {
	DeleteConversationTx(tx *sql.Tx, user1ID, user2ID int64) error
}

type ServiceInterface interface {
	SendRequest(senderID, receiverID int64) error
	AcceptRequest(userID, requestID int64) error
	RejectRequest(userID, requestID int64) error
	DeleteFriend(userID, friendID int64) error
	GetFriends(userID int64) ([]models.User, error)
	GetReceivedRequests(userID int64) ([]models.FriendRequest, error)
	GetSentRequests(userID int64) ([]models.FriendRequest, error)
}

type Service struct {
	repository  RepositoryInterface
	userRepo    UserRepositoryInterface
	messageRepo MessageRepositoryInterface
	db          database.DatabaseInterface
}

func NewService(
	repository RepositoryInterface,
	userRepo UserRepositoryInterface,
	messageRepo MessageRepositoryInterface,
	db database.DatabaseInterface,
) *Service {
	return &Service{
		repository:  repository,
		userRepo:    userRepo,
		messageRepo: messageRepo,
		db:          db,
	}
}

func (s *Service) SendRequest(senderID, receiverID int64) error {
	if senderID == receiverID {
		return ErrCannotAddYourself
	}

	exists, err := s.userRepo.ExistsUser(receiverID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrUserNotFound
	}

	isFriend, err := s.repository.IsFriend(senderID, receiverID)
	if err != nil {
		return err
	}
	if isFriend {
		return ErrAlreadyFriend
	}

	existing, err := s.repository.FindPendingRequest(senderID, receiverID)
	if err != nil {
		return err
	}
	if existing != nil {
		if existing.SenderID == receiverID {
			return s.db.WithTransaction(func(tx *sql.Tx) error {
				if err := s.repository.UpdateRequestStatusTx(tx, existing.ID, RequestAccepted); err != nil {
					return err
				}
				return s.repository.CreateFriendTx(tx, existing.SenderID, existing.ReceiverID)
			})
		}
		return ErrRequestExists
	}

	request := &models.FriendRequest{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     RequestPending,
	}

	return s.repository.CreateRequest(request)
}

func (s *Service) AcceptRequest(userID, requestID int64) error {
	request, err := s.repository.GetRequestByID(requestID)
	if err != nil {
		return err
	}
	if request == nil {
		return ErrRequestNotFound
	}
	if request.ReceiverID != userID {
		return ErrForbidden
	}
	if request.Status != RequestPending {
		return ErrRequestNotFound
	}

	return s.db.WithTransaction(func(tx *sql.Tx) error {
		if err := s.repository.UpdateRequestStatusTx(tx, request.ID, RequestAccepted); err != nil {
			return err
		}
		return s.repository.CreateFriendTx(tx, request.SenderID, request.ReceiverID)
	})
}

func (s *Service) RejectRequest(userID, requestID int64) error {
	request, err := s.repository.GetRequestByID(requestID)
	if err != nil {
		return err
	}
	if request == nil {
		return ErrRequestNotFound
	}
	if request.ReceiverID != userID {
		return ErrForbidden
	}
	if request.Status != RequestPending {
		return ErrRequestNotFound
	}

	return s.repository.UpdateRequestStatus(request.ID, RequestRejected)
}

func (s *Service) DeleteFriend(userID, friendID int64) error {
	if userID == friendID {
		return ErrCannotAddYourself
	}

	isFriend, err := s.repository.IsFriend(userID, friendID)
	if err != nil {
		return err
	}
	if !isFriend {
		return ErrNotFriend
	}

	return s.db.WithTransaction(func(tx *sql.Tx) error {
		if err := s.messageRepo.DeleteConversationTx(tx, userID, friendID); err != nil {
			return err
		}
		return s.repository.DeleteFriendTx(tx, userID, friendID)
	})
}

func (s *Service) GetFriends(userID int64) ([]models.User, error) {
	return s.repository.GetFriends(userID)
}

func (s *Service) GetReceivedRequests(userID int64) ([]models.FriendRequest, error) {
	return s.repository.GetReceivedRequests(userID)
}

func (s *Service) GetSentRequests(userID int64) ([]models.FriendRequest, error) {
	return s.repository.GetSentRequests(userID)
}
