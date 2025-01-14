package become_teacher

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

func (s *Service) Do(ctx context.Context, teacher *entities.Teacher) error {

	exists, err := s.repo.IsUserExistsById(ctx, teacher.UserId)
	if err != nil {
		return fmt.Errorf("failed to find user by id: %w", err)
	}

	if !exists {
		return errors.ErrorUserNotFound
	}

	exists, err = s.repo.IsTeacherExistsByUserId(ctx, teacher.UserId)
	if err != nil {
		return fmt.Errorf("failed to find teacher by id: %w", err)
	}
	if exists {
		return errors.ErrorTeacherExists
	}
	
	if err = s.repo.CreateTeacher(ctx, teacher); err != nil {
		return fmt.Errorf("failed to create teacher: %w", err)
	}

	return nil
}
