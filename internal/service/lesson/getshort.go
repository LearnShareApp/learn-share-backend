package lesson

import (
	"context"
	"errors"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

func (s *LessonService) GetLessonShortData(ctx context.Context, lessonID int) (*entities.Lesson, error) {
	// get lesson
	lesson, err := s.repo.GetLessonByID(ctx, lessonID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return nil, serviceErrs.ErrorLessonNotFound
		}

		return nil, fmt.Errorf("failed to get lesson with id %d: %w", lessonID, err)
	}

	teacher, err := s.repo.GetTeacherByID(ctx, lesson.TeacherID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracher with teacher's id %d: %w", lesson.TeacherID, err)
	}

	lesson.TeacherUserData = &entities.User{
		Id: teacher.UserID,
	}

	return lesson, nil
}
