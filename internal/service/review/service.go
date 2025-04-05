package review

import (
	"context"
	"errors"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

type Repository interface {
	IsUserExistsByID(ctx context.Context, userId int) (bool, error)
	GetTeacherByID(ctx context.Context, teacherId int) (*entities.Teacher, error)
	IsCategoryExistsByID(ctx context.Context, categoryID int) (bool, error)
	GetSkillIdByTeacherIdAndCategoryId(ctx context.Context, teacherId, categoryId int) (int, error)
	IsFinishedLessonExistsByTeacherIdAndStudentIdAndCategoryId(ctx context.Context, teacherId, studentId, categoryId int) (bool, error)
	CreateReview(ctx context.Context, review *entities.Review) error
	IsTeacherExistsById(ctx context.Context, teacherID int) (bool, error)
	GetReviewsByTeacherId(ctx context.Context, teacherID int) ([]*entities.Review, error)
}

type ReviewService struct {
	repo Repository
}

func NewService(repo Repository) *ReviewService {
	return &ReviewService{
		repo: repo,
	}
}

// AddReview create new review about teacher
// Firstly check:
// - is such user (student) exists
// - is this teacher exists
// - is user (student) and teacher the different user
// - is category exists
// - is teacher has such skill
// - is user (student) have had finished lesson with this teacher on this category
func (s *ReviewService) AddReview(ctx context.Context, review *entities.Review) error {
	// is user exists
	exists, err := s.repo.IsUserExistsByID(ctx, review.StudentId)
	if err != nil {
		return fmt.Errorf("failed to check user existstance by id: %w", err)
	}

	if !exists {
		return serviceErrs.ErrorUserNotFound
	}

	// get teacher
	teacher, err := s.repo.GetTeacherByID(ctx, review.TeacherId)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorTeacherNotFound
		}

		return fmt.Errorf("failed to get teacher by id: %w", err)
	}

	// is teacher != student
	if teacher.UserID == review.StudentId {
		return serviceErrs.ErrorStudentAndTeacherSame
	}

	// is category exists
	exists, err = s.repo.IsCategoryExistsByID(ctx, review.CategoryId)
	if err != nil {
		return fmt.Errorf("failed to check category existstance by id: %w", err)
	}

	if !exists {
		return serviceErrs.ErrorCategoryNotFound
	}

	// get skill id
	skillId, err := s.repo.GetSkillIdByTeacherIdAndCategoryId(ctx, review.TeacherId, review.CategoryId)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorSkillUnregistered
		}

		return fmt.Errorf("failed to get skill id by teacher id and category id: %w", err)
	}

	review.SkillId = skillId

	// is student has finished lesson with this teacher and this category
	exists, err = s.repo.IsFinishedLessonExistsByTeacherIdAndStudentIdAndCategoryId(ctx, review.TeacherId, review.StudentId, review.CategoryId)
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

// GetReviews returns all reviews about the teacher.
func (s *ReviewService) GetReviews(ctx context.Context, teacherId int) ([]*entities.Review, error) {
	exists, err := s.repo.IsTeacherExistsById(ctx, teacherId)
	if err != nil {
		return nil, fmt.Errorf("failed to check teacher existence: %w", err)
	}

	if !exists {
		return nil, serviceErrs.ErrorTeacherNotFound
	}

	reviews, err := s.repo.GetReviewsByTeacherId(ctx, teacherId)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews by teacher id: %w", err)
	}

	return reviews, nil
}
