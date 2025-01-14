package get_profile

import (
	"context"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"github.com/LearnShareApp/learn-share-backend/internal/errors"
)

type Service struct {
	repo repo
}

func NewService(repo repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Do(ctx context.Context, id int64) (*entities.User, error) {

	exists, err := s.repo.IsUserExistsById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	if !exists {
		return nil, errors.ErrorUserNotFound
	}

	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
