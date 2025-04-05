package lesson

import (
	"context"
	"errors"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

// GetTeacherLessonList returns all lessons for this teacher.
func (s *LessonService) GetTeacherLessonList(ctx context.Context, userID int) ([]*entities.Lesson, error) {
	// is user exists
	exists, err := s.repo.IsUserExistsByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existstance by id: %w", err)
	}

	if !exists {
		return nil, serviceErrs.ErrorUserNotFound
	}

	// get teacher by userID
	teacher, err := s.repo.GetTeacherByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return nil, serviceErrs.ErrorUserIsNotTeacher
		}

		return nil, fmt.Errorf("failed to get teacher by userID: %w", err)
	}

	// get lessons
	lessons, err := s.repo.GetTeacherLessonsByTeacherID(ctx, teacher.Id)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get teacher's lessons: %w", err)
	}

	return lessons, nil
}
