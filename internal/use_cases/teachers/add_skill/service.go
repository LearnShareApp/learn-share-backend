package add_skill

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

func (s *Service) Do(ctx context.Context, userID, categoryID int, videoCardLink string, about string) error {
	// is user exists
	exists, err := s.repo.IsUserExistsById(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user by id: %w", err)
	}

	if !exists {
		return internalErrs.ErrorUserNotFound
	}

	// is teacher exists by user id
	teacherID, err := s.repo.CreateTeacherIfNotExists(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to create teacher: %w", err)
	}

	// is category exists
	exists, err = s.repo.IsCategoryExistsById(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("failed to find category by id: %w", err)
	}

	if !exists {
		return internalErrs.ErrorCategoryNotFound
	}

	// create skill
	skill := &entities.Skill{
		TeacherId:     teacherID,
		CategoryId:    categoryID,
		VideoCardLink: videoCardLink,
		About:         about,
	}

	if err = s.repo.CreateSkill(ctx, skill); err != nil {
		if errors.Is(err, internalErrs.ErrorNonUniqueData) {
			return internalErrs.ErrorSkillRegistered
		}

		return fmt.Errorf("failed to create skill: %w", err)
	}

	return nil
}
