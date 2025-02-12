package get_teacher

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

func (s *Service) DoByUserId(ctx context.Context, userId int) (*entities.User, error) {

	// get teacher by user id
	teacher, err := s.repo.GetTeacherByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, internalErrs.ErrorSelectEmpty) {
			return nil, internalErrs.ErrorUserIsNotTeacher
		}
		return nil, fmt.Errorf("failed to get teacher: %w", err)
	}

	return s.do(ctx, teacher)
}

func (s *Service) DoByTacherId(ctx context.Context, teacherId int) (*entities.User, error) {

	// get teacher
	teacher, err := s.repo.GetTeacherById(ctx, teacherId)
	if err != nil {
		if errors.Is(err, internalErrs.ErrorSelectEmpty) {
			return nil, internalErrs.ErrorTeacherNotFound
		}
		return nil, fmt.Errorf("failed to get teacher: %w", err)
	}

	return s.do(ctx, teacher)
}

func (s *Service) do(ctx context.Context, teacher *entities.Teacher) (*entities.User, error) {
	// get teacher's user data (common)
	user, err := s.repo.GetUserById(ctx, teacher.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher's user data: %w", err)
	}

	stat, err := s.repo.GetShortStatTeacherById(ctx, teacher.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get stat teacher: %w", err)
	}

	teacher.TeacherStat = *stat

	user.TeacherData = teacher

	teacher.Skills, err = s.repo.GetSkillsByTeacherId(ctx, teacher.Id)
	if err != nil {
		if errors.Is(err, internalErrs.ErrorSelectEmpty) {
			return user, nil
		}
		return nil, fmt.Errorf("failed to get skills: %w", err)
	}

	return user, nil
}
