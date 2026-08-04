package upload

import (
	"context"
	"gomess/internal/modules/upload/dto"
	"gomess/pkg/storage"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
)

type ServiceInterface interface {
	PresignUpload(userID int64, fileName string) (*dto.PresignResponse, error)
}

type Service struct {
	storage storage.StorageInterface
}

func NewService(storage storage.StorageInterface) *Service {
	return &Service{storage: storage}
}

func (s *Service) PresignUpload(userID int64, fileName string) (*dto.PresignResponse, error) {
	ext := filepath.Ext(fileName)
	if len(ext) > 10 {
		ext = ""
	}

	objectKey := "attachments/" + strconv.FormatInt(userID, 10) + "/" + uuid.NewString() + ext

	uploadURL, err := s.storage.PresignedPutURL(context.Background(), objectKey, presignPutExpiry)
	if err != nil {
		return nil, err
	}

	return &dto.PresignResponse{
		UploadURL: uploadURL,
		ObjectKey: objectKey,
	}, nil
}
