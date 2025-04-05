package user

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/pkg/storage/object"

	"github.com/google/uuid"
)

func (s *UserService) EditUser(ctx context.Context, userID int, user *entities.User, avatarReader io.Reader, avatarSize int64) error {
	oldUserData, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorUserNotFound
		}

		return fmt.Errorf("failed to get user: %w", err)
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
		avatarName = uuid.New().String() + ".png"
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

	if err = s.repo.UpdateUser(ctx, userID, oldUserData); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
