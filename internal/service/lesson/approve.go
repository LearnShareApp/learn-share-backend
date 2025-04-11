package lesson

import (
	"context"
	"errors"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

// ApproveLesson approved lesson if all alright.
func (s *LessonService) ApproveLesson(ctx context.Context, userID int, lessonID int) error {
	// is user exists
	exists, err := s.repo.IsUserExistsByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check user existstance by id: %w", err)
	}

	if !exists {
		return serviceErrs.ErrorUserNotFound
	}

	// get teacher
	teacher, err := s.repo.GetTeacherByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorUserIsNotTeacher
		}

		return fmt.Errorf("failed to get teacher by userID: %w", err)
	}

	// get lesson
	lesson, err := s.repo.GetLessonByID(ctx, lessonID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorLessonNotFound
		}

		return fmt.Errorf("failed to get lesson by id: %w", err)
	}

	if lesson.TeacherID != teacher.ID {
		return serviceErrs.ErrorNotRelatedTeacherToLesson
	}

	// get verification statusId
	verificationStatusId, err := s.repo.GetStatusIDByStatusName(ctx, entities.VerificationStatusName)
	if err != nil {
		return fmt.Errorf("failed to get status by status name: %w", err)
	}

	// may approve lesson if only was verification status
	if lesson.StatusID != verificationStatusId {
		return serviceErrs.ErrorStatusNonVerification
	}

	// get waiting statusID
	waitingStatusId, err := s.repo.GetStatusIDByStatusName(ctx, entities.WaitingStatusName)
	if err != nil {
		return fmt.Errorf("failed to get status by status name: %w", err)
	}
	// change lesson status
	if err = s.repo.ChangeLessonStatus(ctx, lessonID, waitingStatusId); err != nil {
		return fmt.Errorf("failed to change lesson status: %w", err)
	}

	return nil
}
