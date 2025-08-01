package lesson

import (
	"context"
	"errors"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

func (s *LessonService) getLessonByID(ctx context.Context, lessonID int) (*entities.Lesson, error) {
	lesson, err := s.repo.GetLessonByID(ctx, lessonID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return nil, serviceErrs.ErrorLessonNotFound
		}
		return nil, fmt.Errorf("failed to get lesson by id: %w", err)
	}
	return lesson, nil
}

func (s *LessonService) getLessonState(ctx context.Context, lesson *entities.Lesson) (*entities.State, error) {
	stateMachineItem, err := s.repo.GetStateMachineItemByID(ctx, lesson.StateMachineItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get statemachine item by id: %w", err)
	}

	currentState, err := s.repo.GetStateByID(ctx, stateMachineItem.StateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current state: %w", err)
	}

	return currentState, nil
}

// changeLessonStateAsStudent set lesson in new state.
// only for local-package usage
func (s *LessonService) changeLessonStateAsStudent(ctx context.Context, userID int, lessonID int, state entities.StateName) error {
	if err := s.validateUserExists(ctx, userID); err != nil {
		return err
	}

	lesson, err := s.getLessonByID(ctx, lessonID)
	if err != nil {
		return err
	}

	if lesson.StudentID != userID {
		return serviceErrs.ErrorNotRelatedUserToLesson
	}

	return s.changeLessonState(ctx, lesson, state)
}

// changeLessonStateAsTeacher set lesson in new state.
// only for local-package usage
func (s *LessonService) changeLessonStateAsTeacher(ctx context.Context, teacherUserID int, lessonID int, state entities.StateName) error {
	if err := s.validateUserExists(ctx, teacherUserID); err != nil {
		return err
	}

	teacher, err := s.repo.GetTeacherByUserID(ctx, teacherUserID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorUserIsNotTeacher
		}
		return fmt.Errorf("failed to get teacher by userID: %w", err)
	}

	lesson, err := s.getLessonByID(ctx, lessonID)
	if err != nil {
		return err
	}

	if lesson.TeacherID != teacher.ID {
		return serviceErrs.ErrorNotRelatedTeacherToLesson
	}

	return s.changeLessonState(ctx, lesson, state)
}

// changeLessonStateAsAny set lesson in new state.
// only for local-package usage
func (s *LessonService) changeLessonStateAsAny(ctx context.Context, userID int, lessonID int, state entities.StateName) error {
	if err := s.validateUserExists(ctx, userID); err != nil {
		return err
	}

	lesson, err := s.getLessonByID(ctx, lessonID)
	if err != nil {
		return err
	}

	if err := s.validateUserIsLessonParticipant(ctx, userID, lesson); err != nil {
		return err
	}

	return s.changeLessonState(ctx, lesson, state)
}

// changeLessonState set lesson in new state.
// only for local-package usage
func (s *LessonService) changeLessonState(ctx context.Context, lesson *entities.Lesson, state entities.StateName) error {
	nextStateID, err := s.repo.GetStateIDByName(ctx, state)
	if err != nil {
		return fmt.Errorf("failed to get next stateID by name: %w", err)
	}

	stateMachineItem, err := s.repo.GetStateMachineItemByID(ctx, lesson.StateMachineItemID)
	if err != nil {
		return fmt.Errorf("failed to get statemachine item by id: %w", err)
	}

	available, err := s.repo.CheckIsTransitionAvailable(
		ctx,
		stateMachineItem.StateMachineID,
		stateMachineItem.StateID,
		nextStateID)
	if err != nil {
		return fmt.Errorf("failed to check if transition is available: %w", err)
	}

	if !available {
		return serviceErrs.ErrorUnavailableStateTransition
	}

	err = s.repo.UpdateStateMachineItemState(ctx, stateMachineItem.ID, nextStateID)
	if err != nil {
		return fmt.Errorf("failed to change statemachine item state: %w", err)
	}

	return nil
}
