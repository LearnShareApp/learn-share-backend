package get_reviews

import (
	"context"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	internalErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

type Service struct {
	repo repo
}

func NewService(repo repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Do(ctx context.Context, teacherId int) ([]*entities.Review, error) {
	exists, err := s.repo.IsTeacherExistsById(ctx, teacherId)
	if err != nil {
		return nil, fmt.Errorf("failed to check teacher existance: %w", err)
	}
	if !exists {
		return nil, internalErrs.ErrorTeacherNotFound
	}

	reviews, err := s.repo.GetReviewsByTeacherId(ctx, teacherId)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews by teacher id: %w", err)
	}

	return reviews, nil
}
