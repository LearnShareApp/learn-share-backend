package registration

import (
	"context"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/pkg/hasher"
)

type JwtService interface {
	GenerateJWTToken(int64) (string, error)
}

type Service struct {
	repo       repo
	jwtService JwtService
}

func NewService(repo repo, service JwtService) *Service {
	return &Service{
		repo:       repo,
		jwtService: service,
	}
}

func (s *Service) Do(ctx context.Context, user *entities.User) (string, error) {
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
