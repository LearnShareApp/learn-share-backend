package minio

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/LearnShareApp/learn-share-backend/pkg/storage/object"

	"github.com/minio/minio-go/v7"
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

func (s *Service) UploadFile(ctx context.Context, file *object.File) error {
	_, err := s.client.PutObject(ctx, s.bucket, file.Name, file.FileReader, file.Size, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

func (s *Service) GetFile(ctx context.Context, fileName string) (*object.File, error) {
	objectReader, err := s.client.GetObject(ctx, s.bucket, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	stat, err := objectReader.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file stat: %w", err)
	}

	return &object.File{
		Name:       fileName,
		Extension:  strings.TrimPrefix(filepath.Ext(fileName), "."),
		FileReader: objectReader,
		Size:       stat.Size,
	}, nil
}

func (s *Service) IsFileExists(ctx context.Context, fileName string) (bool, error) {
	_, err := s.client.StatObject(ctx, s.bucket, fileName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil // Not found
		}

		return false, fmt.Errorf("failed to check file existanse: %w", err)
	}

	return true, nil
}
