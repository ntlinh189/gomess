package storage

import (
	"context"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ObjectInfo struct {
	Exists      bool
	SizeBytes   int64
	ContentType string
}

type ObjectSummary struct {
	Key          string
	LastModified time.Time
}

type StorageInterface interface {
	PresignedPutURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error)
	PresignedGetURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error)
	StatObject(ctx context.Context, objectKey string) (ObjectInfo, error)
	ListObjects(ctx context.Context, prefix string) ([]ObjectSummary, error)
	DeleteObject(ctx context.Context, objectKey string) error
}

type Storage struct {
	client       *minio.Client
	publicClient *minio.Client
	bucket       string
}

func NewStorage(endpoint, publicEndpoint, accessKey, secretKey, bucket string) (*Storage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	if publicEndpoint == "" {
		publicEndpoint = endpoint
	}
	publicClient, err := minio.New(publicEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
		Region: "us-east-1",
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	return &Storage{
		client:       client,
		publicClient: publicClient,
		bucket:       bucket,
	}, nil
}

func (s *Storage) PresignedPutURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	u, err := s.publicClient.PresignedPutObject(ctx, s.bucket, objectKey, expiry)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (s *Storage) PresignedGetURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	u, err := s.publicClient.PresignedGetObject(ctx, s.bucket, objectKey, expiry, url.Values{})
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (s *Storage) StatObject(ctx context.Context, objectKey string) (ObjectInfo, error) {
	info, err := s.client.StatObject(ctx, s.bucket, objectKey, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return ObjectInfo{Exists: false}, nil
		}
		return ObjectInfo{}, err
	}
	return ObjectInfo{
		Exists:      true,
		SizeBytes:   info.Size,
		ContentType: info.ContentType,
	}, nil
}

func (s *Storage) ListObjects(ctx context.Context, prefix string) ([]ObjectSummary, error) {
	var result []ObjectSummary

	for obj := range s.client.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}) {
		if obj.Err != nil {
			return nil, obj.Err
		}
		result = append(result, ObjectSummary{Key: obj.Key, LastModified: obj.LastModified})
	}

	return result, nil
}

func (s *Storage) DeleteObject(ctx context.Context, objectKey string) error {
	return s.client.RemoveObject(ctx, s.bucket, objectKey, minio.RemoveObjectOptions{})
}
