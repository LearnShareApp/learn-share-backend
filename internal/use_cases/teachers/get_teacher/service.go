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

	exists, err := s.repo.IsUserExistsById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to check is user exists: %w", err)
	}
	if !exists {
		return nil, internalErrs.ErrorUserNotFound
	}

	user, err := s.repo.GetUserById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	exists, err = s.repo.IsTeacherExistsByUserId(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to check is teacher exists by user id: %w", err)
	}
	if !exists {
		return nil, internalErrs.ErrorTeacherNotFound
	}

	teacher, err := s.repo.GetTeacherByUserId(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher: %w", err)
	}

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
