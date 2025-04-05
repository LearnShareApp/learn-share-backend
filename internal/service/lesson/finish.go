package lesson

import (
	"context"
	"errors"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

func (s *LessonService) FinishLesson(ctx context.Context, userID int, lessonID int) error {
	// is user exists
	exists, err := s.repo.IsUserExistsByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check user existstance by id: %w", err)
	}

	if !exists {
		return serviceErrs.ErrorUserNotFound
	}

	// get lesson
	lesson, err := s.repo.GetLessonByID(ctx, lessonID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorLessonNotFound
		}

		return fmt.Errorf("failed to get lesson by id: %w", err)
	}

	// get teacher
	teacher, err := s.repo.GetTeacherByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorUserIsNotTeacher
		}

		return fmt.Errorf("failed to get teacher by userId: %w", err)
	}

	if lesson.TeacherID != teacher.Id {
		return serviceErrs.ErrorNotRelatedTeacherToLesson
	}

	// get ongoing statusId
	ongoingStatusId, err := s.repo.GetStatusIDByStatusName(ctx, entities.OngoingStatusName)
	if err != nil {
		return fmt.Errorf("failed to get status by status name: %w", err)
	}

	// may finish lesson if only was ongoing status
	if lesson.StatusId != ongoingStatusId {
		return serviceErrs.ErrorStatusNonOngoing
	}

	// get finished statusId
	finishedStatusId, err := s.repo.GetStatusIDByStatusName(ctx, entities.FinishedStatusName)
	if err != nil {
		return fmt.Errorf("failed to get status by status name: %w", err)
	}
	// change lesson status and remove lesson token
	if err = s.repo.ChangeLessonStatus(ctx, lessonID, finishedStatusId); err != nil {
		return fmt.Errorf("failed to change lesson status: %w", err)
	}

	return nil
}
