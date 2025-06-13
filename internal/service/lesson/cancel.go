package lesson

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

func (s *LessonService) CancelLesson(ctx context.Context, userID, lessonID int) error {
	if err := s.validateUserExists(ctx, userID); err != nil {
		return err
	}

	lesson, err := s.getLessonByID(ctx, lessonID)
	if err != nil {
		return err
	}

	currentState, err := s.getLessonState(ctx, lesson)
	if err != nil {
		return err
	}

	switch currentState.Name {
	case entities.Ongoing:
		//only teacher
		return s.changeLessonStateAsTeacher(ctx, userID, lessonID, entities.Cancel)

	case entities.Planned:
		// teacher or student
		return s.changeLessonStateAsAny(ctx, userID, lessonID, entities.Cancel)

	default:
		return serviceErrs.ErrorUnavailableStateTransition
	}
}
