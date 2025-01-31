package minio

import (
	"context"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/pkg/object_storage"
	"github.com/minio/minio-go/v7"
	"io"
)

type Service struct {
	client *minio.Client
	bucket string
}

func NewService(client *minio.Client, bucket string) *Service {
	return &Service{
		client: client,
		bucket: bucket,
	}
}

func (s *Service) UploadFile(ctx context.Context, file *object_storage.File) error {
	_, err := s.client.PutObject(ctx, s.bucket, file.Name, file.FileReader, file.Size, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

func (s *Service) GetFile(ctx context.Context, fileName string) (io.Reader, error) {
	object, err := s.client.GetObject(ctx, s.bucket, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}
	return object, nil
}
