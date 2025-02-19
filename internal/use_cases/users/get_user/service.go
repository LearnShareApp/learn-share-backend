package get_user

import (
	"context"
	"errors"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	internalErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

type Service struct {
	repo repo
}

func NewService(repo repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Do(ctx context.Context, id int) (*entities.User, error) {

	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, internalErrs.ErrorSelectEmpty) {
			return nil, internalErrs.ErrorUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	stat, err := s.repo.GetUserStatByUserId(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user statistic: %w", err)
	}

	if stat != nil {
		user.Stat = *stat
	}

	user.IsTeacher, err = s.repo.IsTeacherExistsByUserId(ctx, user.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to check whether the user is a teacher: %w", err)
	}

	return user, nil
}
