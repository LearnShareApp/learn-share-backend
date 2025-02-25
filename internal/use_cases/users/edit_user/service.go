package edit_user

import (
	"context"
	"errors"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	internalErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/pkg/hasher"
	"github.com/LearnShareApp/learn-share-backend/pkg/storage/object"
	"github.com/google/uuid"
	"io"
	"time"
)

type ObjectStorageService interface {
	UploadFile(ctx context.Context, file *object.File) error
}

type Service struct {
	repo          repo
	objectStorage ObjectStorageService
}

func NewService(repo repo, storageService ObjectStorageService) *Service {
	return &Service{
		repo:          repo,
		objectStorage: storageService,
	}
}

func (s *Service) Do(ctx context.Context, userId int, user *entities.User, avatarReader io.Reader, avatarSize int64) error {
	oldUserData, err := s.repo.GetUserById(ctx, userId)
	if err != nil {
		if errors.Is(err, internalErrs.ErrorSelectEmpty) {
			return internalErrs.ErrorUserNotFound
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.Password != "" {
		if len(user.Password) < 4 {
			return internalErrs.ErrorPasswordTooShort
		}

		hashedPassword, err := hasher.HashPassword(user.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		oldUserData.Password = hashedPassword
	}

	if user.Name != "" && oldUserData.Name != user.Name {
		oldUserData.Name = user.Name
	}
	if user.Surname != "" && oldUserData.Surname != user.Surname {
		oldUserData.Surname = user.Surname
	}
	if user.Birthdate != oldUserData.Birthdate &&
		user.Birthdate.After(time.Date(1900, 01, 01, 0, 0, 0, 0, time.UTC)) {
		oldUserData.Birthdate = user.Birthdate
	}

	var avatarName string
	if avatarReader != nil {
		avatarName = fmt.Sprintf("%s.png", uuid.New().String())
		file := object.File{
			Name:       avatarName,
			Size:       avatarSize,
			FileReader: avatarReader,
		}

		if err = s.objectStorage.UploadFile(ctx, &file); err != nil {
			return fmt.Errorf("failed to upload avatar: %w", err)
		}
	}

	if avatarName != "" {
		oldUserData.Avatar = avatarName
	}

	if err = s.repo.UpdateUser(ctx, userId, oldUserData); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
