package get_lesson

import (
	"context"
	"errors"
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

func (s *Service) Do(ctx context.Context, lessonId int) (*entities.Lesson, error) {
	// get lesson
	lesson, err := s.repo.GetLessonById(ctx, lessonId)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return nil, serviceErrs.ErrorLessonNotFound
		}
		return nil, fmt.Errorf("failed to get lesson with id %d: %w", lessonId, err)
	}

	lesson.StudentUserData, err = s.repo.GetUserById(ctx, lesson.StudentId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user with id %d: %w", lesson.StudentId, err)
	}

	lesson.TeacherUserData, err = s.repo.GetUserById(ctx, lesson.TeacherId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user with id %d: %w", lesson.TeacherId, err)
	}

	return lesson, nil
}
