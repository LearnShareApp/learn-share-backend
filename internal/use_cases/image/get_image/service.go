package get_image

import (
	"context"
	"fmt"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/pkg/object_storage"
	"io"
	"slices"
)

var SupportedExtension = []string{"png", "jpg", "jpeg"}

type ObjectStorageService interface {
	IsFileExists(ctx context.Context, fileName string) (bool, error)
	GetFile(ctx context.Context, fileName string) (*object_storage.File, error)
}

type Service struct {
	objectStorage ObjectStorageService
}

func NewService(storageService ObjectStorageService) *Service {
	return &Service{
		objectStorage: storageService,
	}
}

func (s *Service) Do(ctx context.Context, filename string) (io.Reader, error) {
	exists, err := s.objectStorage.IsFileExists(ctx, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to check if file exists: %w", err)
	}
	if !exists {
		return nil, serviceErrs.ErrorImageNotFound
	}
	file, err := s.objectStorage.GetFile(ctx, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	if !slices.Contains(SupportedExtension, file.Extension) {
		return nil, serviceErrs.ErrorIncorrectFileFormat
	}

	return file.FileReader, nil
}
