package review

import (
	"context"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

type Repository interface {
	IsUserExistsByID(ctx context.Context, userId int) (bool, error)
	GetTeacherByID(ctx context.Context, teacherId int) (*entities.Teacher, error)
	IsCategoryExistsByID(ctx context.Context, categoryID int) (bool, error)
	GetSkillByTeacherIDAndCategoryID(ctx context.Context, teacherID int, categoryID int) (*entities.Skill, error)
	IsLessonExistsByArgs(ctx context.Context, teacherID int, studentID int, categoryID int, stateName entities.StateName) (bool, error)
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
