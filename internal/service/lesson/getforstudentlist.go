package lesson

import (
	"context"
	"errors"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

// GetStudentLessonList returns all lessons for this student.
func (s *LessonService) GetStudentLessonList(ctx context.Context, userID int) ([]*entities.Lesson, error) {
	// is user exists
	exists, err := s.repo.IsUserExistsByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existstance by id: %w", err)
	}

	if !exists {
		return nil, serviceErrs.ErrorUserNotFound
	}

	// get lessons
	lessons, err := s.repo.GetStudentLessonsByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get student lessons: %w", err)
	}

	return lessons, nil
}
