package review

import (
	"context"
	"errors"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

// CreateReview create new review about teacher
// Firstly check:
// - is such user (student) exists
// - is this teacher exists
// - is user (student) and teacher the different user
// - is category exists
// - is teacher has such skill
// - is user (student) have had finished lesson with this teacher on this category
func (s *ReviewService) CreateReview(ctx context.Context, review *entities.Review) error {
	// is user exists
	err := s.validateUserExists(ctx, review.StudentID)
	if err != nil {
		return err
	}

	// get teacher
	teacher, err := s.repo.GetTeacherByID(ctx, review.TeacherID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorTeacherNotFound
		}

		return fmt.Errorf("failed to get teacher by id: %w", err)
	}

	// is teacher != student
	if teacher.UserID == review.StudentID {
		return serviceErrs.ErrorStudentAndTeacherSame
	}

	// is category exists
	exists, err := s.repo.IsCategoryExistsByID(ctx, review.CategoryID)
	if err != nil {
		return fmt.Errorf("failed to check category existstance by id: %w", err)
	}

	if !exists {
		return serviceErrs.ErrorCategoryNotFound
	}

	// is teacher have such ACTIVE skill
	skill, err := s.repo.GetSkillByTeacherIDAndCategoryID(ctx, teacher.ID, review.CategoryID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorSkillUnregistered
		}
		return fmt.Errorf("failed to get skill by teacher and category: %w", err)
	}

	review.SkillID = skill.ID

	// is student has finished lesson with this teacher and this category
	exists, err = s.repo.IsLessonExistsByArgs(ctx, review.TeacherID, review.StudentID, review.CategoryID, entities.FinishedStatusName)
	if err != nil {
		return fmt.Errorf("failed to check finished lesson existence by teacher id, student id and category id: %w", err)
	}

	if !exists {
		return serviceErrs.ErrorFinishedLessonNotFound
	}

	// create review
	if err = s.repo.CreateReview(ctx, review); err != nil {
		if errors.Is(err, serviceErrs.ErrorNonUniqueData) {
			return serviceErrs.ErrorReviewExists
		}

		return fmt.Errorf("failed to create review: %w", err)
	}

	return nil
}

func (s *ReviewService) validateUserExists(ctx context.Context, userID int) error {
	exists, err := s.repo.IsUserExistsByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check user existence by id: %w", err)
	}
	if !exists {
		return serviceErrs.ErrorUserNotFound
	}
	return nil
}
