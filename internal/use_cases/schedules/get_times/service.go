package get_times

import (
	"context"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"github.com/LearnShareApp/learn-share-backend/internal/errors"
)

type Service struct {
	repo repo
}

func NewService(repo repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Do(ctx context.Context, userId int) ([]*entities.ScheduleTime, error) {
	// is user exists
	exists, err := s.repo.IsUserExistsById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return nil, errors.ErrorUserNotFound
	}

	// is teacher exists
	exists, err = s.repo.IsTeacherExistsByUserId(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to check teacher existstance by user id: %w", err)
	}

	if !exists {
		return nil, errors.ErrorUserIsNotTeacher
	}

	// get teacher
	teacher, err := s.repo.GetTeacherByUserId(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher by user id: %w", err)
	}

	times, err := s.repo.GetScheduleTimesByTeacherId(ctx, teacher.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get available schedule times by teacher id: %w", err)
	}

	return times, nil
}
