package get_times

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

func (s *Service) DoByTeacherId(ctx context.Context, teacherId int) ([]*entities.ScheduleTime, error) {
	// get teacher
	teacher, err := s.repo.GetTeacherById(ctx, teacherId)
	if errors.Is(err, internalErrs.ErrorSelectEmpty) {
		return nil, internalErrs.ErrorTeacherNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher by user id: %w", err)
	}

	return s.do(ctx, teacher)
}

func (s *Service) DoByUserId(ctx context.Context, userId int) ([]*entities.ScheduleTime, error) {
	// get teacher
	teacher, err := s.repo.GetTeacherByUserId(ctx, userId)
	if errors.Is(err, internalErrs.ErrorSelectEmpty) {
		return nil, internalErrs.ErrorUserIsNotTeacher
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher by user id: %w", err)
	}

	return s.do(ctx, teacher)
}

func (s *Service) do(ctx context.Context, teacher *entities.Teacher) ([]*entities.ScheduleTime, error) {
	times, err := s.repo.GetScheduleTimesByTeacherId(ctx, teacher.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get available schedule times by teacher id: %w", err)
	}

	return times, nil
}
