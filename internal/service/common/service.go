package common

import (
	"context"
	"errors"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

type Repository interface {
	GetTeacherByUserID(ctx context.Context, userID int) (*entities.Teacher, error)
	GetTeacherByID(ctx context.Context, teacherID int) (*entities.Teacher, error)
}

type CommonService struct {
	repo Repository
}

func NewService(repo Repository) *CommonService {
	return &CommonService{
		repo: repo,
	}
}

// GetTeacherByID returns teacher by his id.
func (s *CommonService) GetTeacherByID(ctx context.Context, id int) (*entities.Teacher, error) {
	// get teacher
	teacher, err := s.repo.GetTeacherByID(ctx, id)
	if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
		return nil, serviceErrs.ErrorTeacherNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get teacher by id: %w", err)
	}

	return teacher, nil
}

// GetTeacherByUserID returns teacher by his user id.
func (s *CommonService) GetTeacherByUserID(ctx context.Context, userID int) (*entities.Teacher, error) {
	teacher, err := s.repo.GetTeacherByUserID(ctx, userID)
	if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
		return nil, serviceErrs.ErrorTeacherNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get teacher by user id: %w", err)
	}

	return teacher, nil
}
