package get_student_lessons

import (
	"context"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

type Service struct {
	repo repo
}

func NewService(repo repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Do(ctx context.Context, userId int) ([]*entities.Lesson, error) {
	// is user exists
	exists, err := s.repo.IsUserExistsById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existstance by id: %w", err)
	}
	if !exists {
		return nil, serviceErrs.ErrorUserNotFound
	}

	//get lessons
	lessons, err := s.repo.GetStudentLessonsByUserId(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get student lessons: %w", err)
	}

	return lessons, nil
}
