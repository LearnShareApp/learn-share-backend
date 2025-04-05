package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

func (s *UserService) GetUser(ctx context.Context, id int) (*entities.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return nil, serviceErrs.ErrorUserNotFound
		}

		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	stat, err := s.repo.GetUserStatByUserID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user statistic: %w", err)
	}

	if stat != nil {
		user.Stat = *stat
	}

	user.IsTeacher, err = s.repo.IsTeacherExistsByUserID(ctx, user.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to check whether the user is a teacher: %w", err)
	}

	return user, nil
}
