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

func (s *Service) Do(ctx context.Context, userId int) (*entities.User, error) {
	// get user data
	user, err := s.repo.GetUserById(ctx, userId)
	if err != nil {
		if errors.Is(err, internalErrs.ErrorSelectEmpty) {
			return nil, internalErrs.ErrorUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	teacher, err := s.repo.GetTeacherByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, internalErrs.ErrorSelectEmpty) {
			return nil, internalErrs.ErrorTeacherNotFound
		}
		return nil, fmt.Errorf("failed to get teacher: %w", err)
	}

	stat, err := s.repo.GetShortStatTeacherById(ctx, teacher.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get stat teacher: %w", err)
	}

	teacher.TeacherStat = *stat

	user.TeacherData = teacher

	teacher.Skills, err = s.repo.GetSkillsByTeacherId(ctx, teacher.Id)
	if err != nil {
		if errors.Is(err, internalErrs.ErrorSelectEmpty) {
			return user, nil
		}
		return nil, fmt.Errorf("failed to get skills: %w", err)
	}

	return user, nil
}
