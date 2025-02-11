package get_teacher_lessons

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

func (s *Service) Do(ctx context.Context, userId int) ([]*entities.Lesson, error) {
	// is user exists
	exists, err := s.repo.IsUserExistsById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existstance by id: %w", err)
	}
	if !exists {
		return nil, serviceErrs.ErrorUserNotFound
	}

	// get teacher by userId
	teacher, err := s.repo.GetTeacherByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return nil, serviceErrs.ErrorUserIsNotTeacher
		}
		return nil, fmt.Errorf("failed to get teacher by userId: %w", err)
	}

	//get lessons
	lessons, err := s.repo.GetTeacherLessonsByTeacherId(ctx, teacher.Id)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get teacher lessons: %w", err)
	}

	return lessons, nil
}
