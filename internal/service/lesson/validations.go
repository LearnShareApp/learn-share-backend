package lesson

import (
	"context"
	"errors"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

func (s *LessonService) validateUserExists(ctx context.Context, userID int) error {
	exists, err := s.repo.IsUserExistsByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check user existence by id: %w", err)
	}
	if !exists {
		return serviceErrs.ErrorUserNotFound
	}
	return nil
}

func (s *LessonService) validateUserIsLessonParticipant(ctx context.Context, userID int, lesson *entities.Lesson) error {
	if lesson.StudentID == userID {
		return nil
	}

	// if no student, check is it a teacher
	teacher, err := s.repo.GetTeacherByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorNotRelatedUserToLesson
		}
		return fmt.Errorf("failed to get teacher by userID: %w", err)
	}

	if lesson.TeacherID != teacher.ID {
		return serviceErrs.ErrorNotRelatedUserToLesson
	}

	return nil
}
