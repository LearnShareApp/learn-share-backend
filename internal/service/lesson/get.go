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

		return nil, fmt.Errorf("failed to get lesson with %d id: %w", lessonID, err)
	}

	// get state machine
	lesson.StateMachineItem, err = s.repo.GetStateMachineItemByID(ctx, lesson.StateMachineItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get statemachine item with %d id: %w", lesson.StateMachineItemID, err)
	}

	lesson.StudentUserData, err = s.repo.GetUserByID(ctx, lesson.StudentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user with %d id: %w", lesson.StudentID, err)
	}

	teacher, err := s.repo.GetTeacherByID(ctx, lesson.TeacherID)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher with %d id: %w", lesson.TeacherID, err)
	}

	lesson.TeacherUserData, err = s.repo.GetUserByID(ctx, teacher.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user with %d id: %w", lesson.TeacherID, err)
	}

	return lesson, nil
}
