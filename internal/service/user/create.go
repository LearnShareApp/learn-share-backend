package user

import (
	"context"
	"fmt"
	"io"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/pkg/hasher"
	"github.com/LearnShareApp/learn-share-backend/pkg/storage/object"

	"github.com/google/uuid"
)

// CreateUser save new user and returns his id from db.
func (s *UserService) CreateUser(ctx context.Context, user *entities.User, avatarReader io.Reader, avatarSize int64) (int, error) {
	exists, err := s.repo.IsUserExistsByEmail(ctx, user.Email)
	if err != nil {
		return 0, fmt.Errorf("failed to find user: %w", err)
	}

	if exists {
		return 0, serviceErrs.ErrorUserExists
	}

	if len(user.Password) < 4 {
		return 0, serviceErrs.ErrorPasswordTooShort
	}

	hashedPassword, err := hasher.HashPassword(user.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword

	var avatarName string
	if avatarReader != nil {
		avatarName = uuid.New().String() + ".png"
		file := object.File{
			Name:       avatarName,
			Size:       avatarSize,
			FileReader: avatarReader,
		}

		if err = s.objectStorage.UploadFile(ctx, &file); err != nil {
			return 0, fmt.Errorf("failed to upload avatar: %w", err)
		}
	}

	user.Avatar = avatarName

	userID, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("failed to save user: %w", err)
	}

	return userID, nil
}
