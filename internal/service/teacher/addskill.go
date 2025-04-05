package teacher

import (
	"context"
	"errors"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

func (s *TeacherService) AddSkill(ctx context.Context, userID, categoryID int, videoCardLink string, about string) error {
	// is user exists
	exists, err := s.repo.IsUserExistsByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user by id: %w", err)
	}

	if !exists {
		return serviceErrs.ErrorUserNotFound
	}

	// is teacher exists by user id
	teacherID, err := s.repo.CreateTeacherIfNotExists(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to create teacher: %w", err)
	}

	// is category exists
	exists, err = s.repo.IsCategoryExistsByID(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("failed to find category by id: %w", err)
	}

	if !exists {
		return serviceErrs.ErrorCategoryNotFound
	}

	// create skill
	skill := &entities.Skill{
		TeacherId:     teacherID,
		CategoryId:    categoryID,
		VideoCardLink: videoCardLink,
		About:         about,
	}

	if err = s.repo.CreateSkill(ctx, skill); err != nil {
		if errors.Is(err, serviceErrs.ErrorNonUniqueData) {
			return serviceErrs.ErrorSkillRegistered
		}

		return fmt.Errorf("failed to create skill: %w", err)
	}

	return nil
}
