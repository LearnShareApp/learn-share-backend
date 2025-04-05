package teacher

import (
	"context"
	"fmt"

	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

func (s *TeacherService) BecomeTeacher(ctx context.Context, userID int) error {
	exists, err := s.repo.IsUserExistsByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user by id: %w", err)
	}

	if !exists {
		return serviceErrs.ErrorUserNotFound
	}

	exists, err = s.repo.IsTeacherExistsByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find teacher by id: %w", err)
	}

	if exists {
		return serviceErrs.ErrorTeacherExists
	}

	if err = s.repo.CreateTeacher(ctx, userID); err != nil {
		return fmt.Errorf("failed to create teacher: %w", err)
	}

	return nil
}
