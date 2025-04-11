package schedule

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

type Repository interface {
	IsUserExistsByID(ctx context.Context, id int) (bool, error)
	IsScheduleTimeExistsByTeacherIDAndDatetime(ctx context.Context, id int, datetime time.Time) (bool, error)
	GetTeacherByUserID(ctx context.Context, userId int) (*entities.Teacher, error)
	GetScheduleTimesByTeacherID(ctx context.Context, id int) ([]*entities.ScheduleTime, error)
	CreateScheduleTime(ctx context.Context, teacherId int, datetime time.Time) error
}

type ScheduleService struct {
	repo Repository
}

func NewService(repo Repository) *ScheduleService {
	return &ScheduleService{
		repo: repo,
	}
}

// AddTime create new time in schedule for teacher
// Firstly check:
// - is such user (teacher) exists
// - is this user a teacher
// - is this schedule time not already exists
func (s *ScheduleService) AddTime(ctx context.Context, userID int, datetime time.Time) error {
	// is user exists
	exists, err := s.repo.IsUserExistsByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check user existstance by id: %w", err)
	}

	if !exists {
		return serviceErrs.ErrorUserNotFound
	}

	// get teacher
	teacher, err := s.repo.GetTeacherByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorUserIsNotTeacher
		}

		return fmt.Errorf("failed to get teacher by user id: %w", err)
	}

	exists, err = s.repo.IsScheduleTimeExistsByTeacherIDAndDatetime(ctx, teacher.ID, datetime)
	if err != nil {
		return fmt.Errorf("failed to check time existstance by user id: %w", err)
	}

	if exists {
		return serviceErrs.ErrorScheduleTimeExists
	}

	if err = s.repo.CreateScheduleTime(ctx, teacher.ID, datetime); err != nil {
		return fmt.Errorf("failed to create teacher time: %w", err)
	}

	return nil
}

// GetTimes returns teacher's schedule times.
func (s *ScheduleService) GetTimes(ctx context.Context, teacher *entities.Teacher) ([]*entities.ScheduleTime, error) {
	times, err := s.repo.GetScheduleTimesByTeacherID(ctx, teacher.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get available schedule times by teacher id: %w", err)
	}

	return times, nil
}
