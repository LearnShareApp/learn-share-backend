package lesson

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

func (s *LessonService) JoinLesson(ctx context.Context, userID int, lessonID int) (string, error) {
	err := s.validateUserExists(ctx, userID)
	if err != nil {
		return "", err
	}

	lesson, err := s.getLessonByID(ctx, lessonID)
	if err != nil {
		return "", err
	}

	if err = s.validateUserIsLessonParticipant(ctx, userID, lesson); err != nil {
		return "", err
	}

	currentState, err := s.getLessonState(ctx, lesson)
	if err != nil {
		return "", err
	}

	if currentState.Name != entities.Ongoing {
		return "", serviceErrs.ErrorUnavailableOperationState
	}

	return s.generateLessonMeetingToken(ctx, userID, lessonID)
}
