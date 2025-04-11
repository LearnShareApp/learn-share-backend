package lesson

import (
	"context"
	"errors"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

func (s *LessonService) CancelLesson(ctx context.Context, userID int, lessonID int) error {
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

	// if user is not student => maybe he is a teacher (check it)
	if lesson.StudentID != userID {
		teacher, err := s.repo.GetTeacherByUserID(ctx, userID)
		if err != nil {
			// if it's not the teacher for this lesson either => error
			if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
				return serviceErrs.ErrorNotRelatedUserToLesson
			}

			return fmt.Errorf("failed to get teacher by userID: %w", err)
		}

		if lesson.TeacherID != teacher.ID {
			return serviceErrs.ErrorNotRelatedUserToLesson
		}
	}

	// if we here => related user to lesson

	//can cancel only if lesson isn't already canceled or finished
	finishStatusID, err := s.repo.GetStatusIDByStatusName(ctx, entities.FinishedStatusName)
	if err != nil {
		return fmt.Errorf("failed to get status by status name: %w", err)
	}

	if lesson.StatusID == finishStatusID {
		return serviceErrs.ErrorFinishedLessonCanNotBeCancel
	}

	// get cancel cancelStatusID
	cancelStatusID, err := s.repo.GetStatusIDByStatusName(ctx, entities.CancelStatusName)
	if err != nil {
		return fmt.Errorf("failed to get status by status name: %w", err)
	}

	if lesson.StatusID == cancelStatusID {
		return serviceErrs.ErrorLessonAlreadyCanceled
	}

	// change lesson status
	if err = s.repo.ChangeLessonStatus(ctx, lessonID, cancelStatusID); err != nil {
		return fmt.Errorf("failed to change lesson status: %w", err)
	}

	return nil
}
