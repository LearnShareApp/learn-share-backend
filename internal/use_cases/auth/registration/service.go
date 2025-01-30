package registration

import (
	"context"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/pkg/hasher"
	"github.com/google/uuid"
	"io"
)

type JwtService interface {
	GenerateJWTToken(int) (string, error)
}

type ObjectStorageService interface {
	UploadFile(ctx context.Context, fileName string, file io.Reader, fileSize int64) error
}

type Service struct {
	repo          repo
	jwtService    JwtService
	objectStorage ObjectStorageService
}

func NewService(repo repo, service JwtService, storageService ObjectStorageService) *Service {
	return &Service{
		repo:          repo,
		jwtService:    service,
		objectStorage: storageService,
	}
}

func (s *Service) Do(ctx context.Context, user *entities.User, avatarReader *io.Reader, avatarSize int64) (string, error) {
	exists, err := s.repo.IsUserExistsByEmail(ctx, user.Email)
	if err != nil {
		return "", fmt.Errorf("failed to find user: %w", err)
	}

	if exists {
		return "", errors.ErrorUserExists
	}

	if len(user.Password) < 4 {
		return "", errors.ErrorPasswordTooShort
	}

	hashedPassword, err := hasher.HashPassword(user.Password)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	var avatarName string
	if avatarReader != nil {
		avatarName = fmt.Sprintf("%s.png", uuid.New().String())
		if err = s.objectStorage.UploadFile(ctx, avatarName, *avatarReader, avatarSize); err != nil {
			return "", fmt.Errorf("failed to upload avatar: %w", err)
		}
	}
	user.Avatar = avatarName

	userId, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to save user: %w", err)
	}

	token, err := s.jwtService.GenerateJWTToken(userId)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
