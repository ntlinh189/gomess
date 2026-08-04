package message

import (
	"context"
	"database/sql"
	"gomess/internal/database"
	"gomess/internal/models"
	"gomess/internal/modules/message/dto"
	"gomess/pkg/storage"
	"strings"
	"time"
)

type FriendRepositoryInterface interface {
	IsFriend(user1ID, user2ID int64) (bool, error)
}

type Broadcaster interface {
	SendToUser(userID int64, event string, data any)
}

type ServiceInterface interface {
	SendMessage(senderID, receiverID int64, content string, attachments []dto.AttachmentInput) (*dto.MessageResponse, error)
	GetHistory(userID, friendID int64, beforeID int64, limit int) ([]dto.MessageResponse, error)
	DeleteForMe(userID, messageID int64) error
	RevokeMessage(userID, messageID int64) error
}

type Service struct {
	repo        RepositoryInterface
	friendRepo  FriendRepositoryInterface
	storage     storage.StorageInterface
	db          database.DatabaseInterface
	broadcaster Broadcaster
}

func NewService(
	repo RepositoryInterface,
	friendRepo FriendRepositoryInterface,
	storage storage.StorageInterface,
	db database.DatabaseInterface,
	broadcaster Broadcaster,
) *Service {
	return &Service{
		repo:        repo,
		friendRepo:  friendRepo,
		storage:     storage,
		db:          db,
		broadcaster: broadcaster,
	}
}

func (s *Service) SendMessage(senderID, receiverID int64, content string, attachments []dto.AttachmentInput) (*dto.MessageResponse, error) {
	if senderID == receiverID {
		return nil, ErrCannotMessageYourself
	}

	if content == "" && len(attachments) == 0 {
		return nil, ErrEmptyMessage
	}

	if len(attachments) > maxAttachmentsPerMessage {
		return nil, ErrTooManyAttachments
	}

	isFriend, err := s.friendRepo.IsFriend(senderID, receiverID)
	if err != nil {
		return nil, err
	}
	if !isFriend {
		return nil, ErrNotFriend
	}

	messageAttachments := make([]models.MessageAttachment, 0, len(attachments))
	for _, a := range attachments {
		info, err := s.storage.StatObject(context.Background(), a.ObjectKey)
		if err != nil {
			return nil, err
		}
		if !info.Exists {
			return nil, ErrAttachmentNotFound
		}

		messageAttachments = append(messageAttachments, models.MessageAttachment{
			Type:      detectAttachmentType(info.ContentType),
			ObjectKey: a.ObjectKey,
			FileName:  a.FileName,
			MimeType:  info.ContentType,
			SizeBytes: info.SizeBytes,
		})
	}

	message := &models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
	}

	err = s.db.WithTransaction(func(tx *sql.Tx) error {
		if err := s.repo.CreateTx(tx, message); err != nil {
			return err
		}
		return s.repo.CreateAttachmentsTx(tx, message.ID, messageAttachments)
	})
	if err != nil {
		return nil, err
	}

	message.Attachments = messageAttachments

	resp, err := s.toMessageResponse(message)
	if err != nil {
		return nil, err
	}

	s.broadcaster.SendToUser(receiverID, eventNewMessage, resp)
	s.broadcaster.SendToUser(senderID, eventNewMessage, resp)

	return resp, nil
}

func (s *Service) GetHistory(userID, friendID int64, beforeID int64, limit int) ([]dto.MessageResponse, error) {
	isFriend, err := s.friendRepo.IsFriend(userID, friendID)
	if err != nil {
		return nil, err
	}
	if !isFriend {
		return nil, ErrNotFriend
	}

	if limit <= 0 {
		limit = defaultHistoryLimit
	}
	if limit > maxHistoryLimit {
		limit = maxHistoryLimit
	}

	messages, err := s.repo.GetHistory(userID, friendID, beforeID, limit)
	if err != nil {
		return nil, err
	}

	messageIDs := make([]int64, len(messages))
	for i, m := range messages {
		messageIDs[i] = m.ID
	}

	attachmentsByMessage, err := s.repo.GetAttachmentsByMessageIDs(messageIDs)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.MessageResponse, 0, len(messages))
	for _, m := range messages {
		m.Attachments = attachmentsByMessage[m.ID]

		r, err := s.toMessageResponse(&m)
		if err != nil {
			return nil, err
		}
		resp = append(resp, *r)
	}

	return resp, nil
}

func (s *Service) DeleteForMe(userID, messageID int64) error {
	message, err := s.repo.GetByID(messageID)
	if err != nil {
		return err
	}
	if message == nil {
		return ErrMessageNotFound
	}
	if message.SenderID != userID && message.ReceiverID != userID {
		return ErrForbidden
	}

	if err := s.repo.CreateDeletion(messageID, userID); err != nil {
		return err
	}

	s.broadcaster.SendToUser(userID, eventMessageDeleted, map[string]int64{"id": messageID})

	return nil
}

func (s *Service) RevokeMessage(userID, messageID int64) error {
	message, err := s.repo.GetByID(messageID)
	if err != nil {
		return err
	}
	if message == nil {
		return ErrMessageNotFound
	}
	if message.SenderID != userID {
		return ErrForbidden
	}
	if message.RevokedAt != nil {
		return ErrAlreadyRevoked
	}
	if time.Since(message.CreatedAt) > revokeWindow {
		return ErrRevokeWindowExpired
	}

	if err := s.repo.RevokeMessage(messageID); err != nil {
		return err
	}

	payload := map[string]int64{
		"id": messageID, "sender_id": message.SenderID, "receiver_id": message.ReceiverID,
	}
	s.broadcaster.SendToUser(message.SenderID, eventMessageRevoked, payload)
	s.broadcaster.SendToUser(message.ReceiverID, eventMessageRevoked, payload)

	return nil
}

func (s *Service) toMessageResponse(m *models.Message) (*dto.MessageResponse, error) {
	if m.RevokedAt != nil {
		return &dto.MessageResponse{
			ID:         m.ID,
			SenderID:   m.SenderID,
			ReceiverID: m.ReceiverID,
			Revoked:    true,
			CreatedAt:  m.CreatedAt,
		}, nil
	}

	attachments := make([]dto.AttachmentResponse, 0, len(m.Attachments))
	for _, a := range m.Attachments {
		url, err := s.storage.PresignedGetURL(context.Background(), a.ObjectKey, presignGetExpiry)
		if err != nil {
			return nil, err
		}

		attachments = append(attachments, dto.AttachmentResponse{
			ID:        a.ID,
			Type:      a.Type,
			URL:       url,
			FileName:  a.FileName,
			MimeType:  a.MimeType,
			SizeBytes: a.SizeBytes,
		})
	}

	return &dto.MessageResponse{
		ID:          m.ID,
		SenderID:    m.SenderID,
		ReceiverID:  m.ReceiverID,
		Content:     m.Content,
		Attachments: attachments,
		Revoked:     false,
		CreatedAt:   m.CreatedAt,
	}, nil
}

func detectAttachmentType(mimeType string) string {
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return "image"
	case strings.HasPrefix(mimeType, "video/"):
		return "video"
	case strings.HasPrefix(mimeType, "audio/"):
		return "audio"
	default:
		return "file"
	}
}
