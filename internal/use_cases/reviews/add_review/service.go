package add_review

import (
	"context"
	"errors"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

type Service struct {
	repo repo
}

func NewService(repo repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Do(ctx context.Context, review *entities.Review) error {
	// is user exists
	exists, err := s.repo.IsUserExistsById(ctx, review.StudentId)
	if err != nil {
		return fmt.Errorf("failed to check user existstance by id: %w", err)
	}
	if !exists {
		return serviceErrs.ErrorUserNotFound
	}

	// is teacher exists
	exists, err = s.repo.IsTeacherExistsById(ctx, review.TeacherId)
	if err != nil {
		return fmt.Errorf("failed to check teacher existstance by user id: %w", err)
	}
	if !exists {
		return serviceErrs.ErrorTeacherNotFound
	}
	// get teacher
	teacher, err := s.repo.GetTeacherById(ctx, review.TeacherId)
	if err != nil {
		return fmt.Errorf("failed to get teacher by id: %w", err)
	}
	// is teacher == student
	if teacher.UserId == review.StudentId {
		return serviceErrs.ErrorStudentAndTeacherSame
	}

	// is category exists
	exists, err = s.repo.IsCategoryExistsById(ctx, review.CategoryId)
	if err != nil {
		return fmt.Errorf("failed to check category existstance by id: %w", err)
	}
	if !exists {
		return serviceErrs.ErrorCategoryNotFound
	}

	// is skill exists by teacher id and category id
	exists, err = s.repo.IsSkillExistsByTeacherIdAndCategoryId(ctx, review.TeacherId, review.CategoryId)
	if err != nil {
		return fmt.Errorf("failed to check skill existstance by teacher id and category id: %w", err)
	}
	if !exists {
		return serviceErrs.ErrorSkillUnregistered
	}

	// is student has finished lesson with this teacher and this category
	exists, err = s.repo.IsFinishedLessonExistsByTeacherIdAndStudentIdAndCategoryId(ctx, review.TeacherId, review.StudentId, review.CategoryId)
	if err != nil {
		return fmt.Errorf("failed to check finished lesson existence by teacher id, student id and category id: %w", err)
	}
	if !exists {
		return serviceErrs.ErrorFinishedLessonNotFound
	}

	// get skill id
	skillId, err := s.repo.GetSkillIdByTeacherIdAndCategoryId(ctx, review.TeacherId, review.CategoryId)
	if err != nil {
		return fmt.Errorf("failed to get skill id by teacher id and category id: %w", err)
	}
	review.SkillId = skillId

	// create review
	if err = s.repo.CreateReview(ctx, review); err != nil {
		if errors.Is(err, serviceErrs.ErrorNonUniqueData) {
			return serviceErrs.ErrorReviewExists
		}
		return fmt.Errorf("failed to create review: %w", err)
	}

	return nil
}
