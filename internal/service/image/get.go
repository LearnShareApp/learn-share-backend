package image

import (
	"context"
	"fmt"
	"io"
	"slices"

	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

var SupportedExtension = []string{"png", "jpg", "jpeg"}

// GetImage returns image reader by image name.
func (s *ImageService) GetImage(ctx context.Context, filename string) (io.Reader, error) {
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
