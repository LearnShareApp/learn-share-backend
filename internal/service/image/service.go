package image

import (
	"context"

	"github.com/LearnShareApp/learn-share-backend/pkg/storage/object"
)

type ObjectStorage interface {
	IsFileExists(ctx context.Context, fileName string) (bool, error)
	GetFile(ctx context.Context, fileName string) (*object.File, error)
}

type ImageService struct {
	objectStorage ObjectStorage
}

func NewService(objectStorage ObjectStorage) *ImageService {
	return &ImageService{
		objectStorage: objectStorage,
	}
}
