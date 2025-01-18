package get_teacher

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

func (s *Service) Do(ctx context.Context, id int) (*entities.User, *entities.Teacher, error) {

	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, internalErrs.ErrorSelectEmpty) {
			return nil, nil, internalErrs.ErrorUserNotFound
		}
		return nil, nil, fmt.Errorf("failed to get user: %w", err)
	}

	teacher, err := s.repo.GetTeacherByUserId(ctx, id)
	if err != nil {
		if errors.Is(err, internalErrs.ErrorSelectEmpty) {
			return nil, nil, internalErrs.ErrorTeacherNotFound
		}
		return nil, nil, fmt.Errorf("failed to get teacher: %w", err)
	}

	teacher.Skills, err = s.repo.GetSkillsByTeacherId(ctx, id)
	if err != nil {
		if errors.Is(err, internalErrs.ErrorSelectEmpty) {
			return user, teacher, nil
		}
		return nil, nil, fmt.Errorf("failed to get skills: %w", err)
	}

	return user, teacher, nil
}
