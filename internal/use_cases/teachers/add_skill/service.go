package add_skill

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

func (s *Service) Do(ctx context.Context, userId int, categoryId int, videoCardLink string, about string) error {
	// TODO: check whether user is teacher, check is he not already teach this and add skill to db

	// is user exists
	exists, err := s.repo.IsUserExistsById(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to find user by id: %w", err)
	}
	if !exists {
		return errors.ErrorUserNotFound
	}

	// is teacher exists by user id
	teacherId, err := s.repo.CreateTeacherIfNotExists(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to create teacher: %w", err)
	}

	// is category exists
	exists, err = s.repo.IsCategoryExistsById(ctx, categoryId)
	if err != nil {
		return fmt.Errorf("failed to find category by id: %w", err)
	}
	if !exists {
		return errors.ErrorCategoryNotFound
	}

	// create skill
	skill := &entities.Skill{
		TeacherId:     teacherId,
		CategoryId:    categoryId,
		VideoCardLink: videoCardLink,
		About:         about,
	}

	if err = s.repo.CreateSkill(ctx, skill); err != nil {
		if err == errors.ErrorNonUniqueData {
			return errors.ErrorSkillRegistered
		}
		return fmt.Errorf("failed to create skill: %w", err)
	}

	return nil
}
