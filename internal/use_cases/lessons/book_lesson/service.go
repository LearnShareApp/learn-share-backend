package book_lesson

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

func (s *Service) Do(ctx context.Context, lesson *entities.Lesson) error {
	// is student exists
	exists, err := s.repo.IsUserExistsById(ctx, lesson.StudentId)
	if err != nil {
		return fmt.Errorf("failed to check user existstance by id: %w", err)
	}
	if !exists {
		return serviceErrs.ErrorUserNotFound
	}

	// is student != teacher
	teacher, err := s.repo.GetTeacherById(ctx, lesson.TeacherId)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorTeacherNotFound
		}
		return fmt.Errorf("failed to get teacher by id: %w", err)
	}
	if teacher.UserId == lesson.StudentId {
		return serviceErrs.ErrorStudentAndTeacherSame
	}

	// is category exists
	exists, err = s.repo.IsCategoryExistsById(ctx, lesson.CategoryId)
	if err != nil {
		return fmt.Errorf("failed to check category existstance by id: %w", err)
	}
	if !exists {
		return serviceErrs.ErrorCategoryNotFound
	}

	// is teacher have such skill
	exists, err = s.repo.IsSkillExistsByTeacherIdAndCategoryId(ctx, lesson.TeacherId, lesson.CategoryId)
	if err != nil {
		return fmt.Errorf("failed to check skill existstance by id: %w", err)
	}
	if !exists {
		return serviceErrs.ErrorSkillUnregistered
	}

	// is this time still available and owner is this teacher
	scheduleTime, err := s.repo.GetScheduleTimeById(ctx, lesson.ScheduleTimeId)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorScheduleTimeNotFound
		}
		return fmt.Errorf("failed to get schedule time by id: %w", err)
	}
	if scheduleTime.TeacherId != lesson.TeacherId {
		return serviceErrs.ErrorScheduleTimeForAnotherTeacher
	}
	if !scheduleTime.IsAvailable {
		return serviceErrs.ErrorScheduleTimeUnavailable
	}

	// create unconfirmed lesson
	if err = s.repo.CreateUnconfirmedLesson(ctx, lesson); err != nil {
		if errors.Is(err, serviceErrs.ErrorNonUniqueData) {
			return serviceErrs.ErrorLessonTimeBooked
		}
		return fmt.Errorf("failed to create unconfirmed lesson: %w", err)
	}

	return nil
}
