package get_lesson_shortdata

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

	teacher, err := s.repo.GetTeacherById(ctx, lesson.TeacherId)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracher with teacher's id %d: %w", lesson.TeacherId, err)
	}

	lesson.TeacherUserData = &entities.User{
		Id: teacher.UserId,
	}

	return lesson, nil
}
