package lesson

import (
	"context"
	"errors"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

func (s *LessonService) GetLesson(ctx context.Context, lessonID int) (*entities.Lesson, error) {
	// get lesson
	lesson, err := s.repo.GetLessonByID(ctx, lessonID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return nil, serviceErrs.ErrorLessonNotFound
		}

		return nil, fmt.Errorf("failed to get lesson with id %d: %w", lessonID, err)
	}

	lesson.StudentUserData, err = s.repo.GetUserByID(ctx, lesson.StudentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user with id %d: %w", lesson.StudentID, err)
	}

	lesson.TeacherUserData, err = s.repo.GetUserByID(ctx, lesson.TeacherID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user with id %d: %w", lesson.TeacherID, err)
	}

	return lesson, nil
}
