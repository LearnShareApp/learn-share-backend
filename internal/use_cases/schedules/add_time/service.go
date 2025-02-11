package add_time

import (
	"context"
	"errors"
	"fmt"
	intenalErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"time"
)

type Service struct {
	repo repo
}

func NewService(repo repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Do(ctx context.Context, userId int, datetime time.Time) error {
	// is user exists
	exists, err := s.repo.IsUserExistsById(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to check user existstance by id: %w", err)
	}

	if !exists {
		return intenalErrs.ErrorUserNotFound
	}

	// get teacher
	teacher, err := s.repo.GetTeacherByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, intenalErrs.ErrorSelectEmpty) {
			return intenalErrs.ErrorUserIsNotTeacher
		}
		return fmt.Errorf("failed to get teacher by user id: %w", err)
	}

	exists, err = s.repo.IsScheduleTimeExistsByTeacherIdAndDatetime(ctx, teacher.Id, datetime)
	if err != nil {
		return fmt.Errorf("failed to check time existstance by user id: %w", err)
	}
	if exists {
		return intenalErrs.ErrorScheduleTimeExists
	}

	if err = s.repo.CreateScheduleTime(ctx, teacher.Id, datetime); err != nil {
		return fmt.Errorf("failed to create teacher time: %w", err)
	}

	return nil
}
