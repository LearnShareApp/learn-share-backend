package get_lesson_shortdata

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

func (s *Service) Do(ctx context.Context, lessonId int) (*entities.Lesson, error) {
	// is lesson exists
	exists, err := s.repo.IsLessonExistsById(ctx, lessonId)
	if err != nil {
		return nil, fmt.Errorf("failed to check existence of lesson with id %d: %w", lessonId, err)
	}

	if !exists {
		return nil, serviceErrs.ErrorLessonNotFound
	}

	lesson, err := s.repo.GetLessonById(ctx, lessonId)
	if err != nil {
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
