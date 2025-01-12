package login

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

func (s *Service) Do(ctx context.Context, reqUser *entities.User) (string, error) {
	exists, err := s.repo.ExistsByEmail(ctx, reqUser.Email)
	if err != nil {
		return "", fmt.Errorf("failed to find reqUser: %w", err)
	}

	if !exists {
		return "", errors.ErrorUserNotFound
	}

	realUser, err := s.repo.GetUserByEmail(ctx, reqUser.Email)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	if !hasher.ComparePassword(reqUser.Password, realUser.Password) {
		return "", errors.ErrorPasswordIncorrect
	}

	token, err := s.jwtService.GenerateJWTToken(realUser.Id)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
