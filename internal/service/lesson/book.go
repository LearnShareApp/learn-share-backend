package lesson

import (
	"context"
	"errors"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

func (s *LessonService) BookLesson(ctx context.Context, lesson *entities.Lesson) error {
	// is student exists
	exists, err := s.repo.IsUserExistsByID(ctx, lesson.StudentID)
	if err != nil {
		return fmt.Errorf("failed to check user existstance by id: %w", err)
	}

	if !exists {
		return serviceErrs.ErrorUserNotFound
	}

	// is student != teacher
	teacher, err := s.repo.GetTeacherByID(ctx, lesson.TeacherID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorTeacherNotFound
		}

		return fmt.Errorf("failed to get teacher by id: %w", err)
	}

	if teacher.UserID == lesson.StudentID {
		return serviceErrs.ErrorStudentAndTeacherSame
	}

	// is categories exists
	exists, err = s.repo.IsCategoryExistsByID(ctx, lesson.CategoryID)
	if err != nil {
		return fmt.Errorf("failed to check categories existstance by id: %w", err)
	}

	if !exists {
		return serviceErrs.ErrorCategoryNotFound
	}

	// is teacher have such skill
	exists, err = s.repo.IsSkillExistsByTeacherIDAndCategoryID(ctx, lesson.TeacherID, lesson.CategoryID)
	if err != nil {
		return fmt.Errorf("failed to check skill existstance by id: %w", err)
	}

	if !exists {
		return serviceErrs.ErrorSkillUnregistered
	}

	// is this time still available and owner is this teacher
	scheduleTime, err := s.repo.GetScheduleTimeByID(ctx, lesson.ScheduleTimeID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorScheduleTimeNotFound
		}

		return fmt.Errorf("failed to get schedules time by id: %w", err)
	}

	if scheduleTime.TeacherId != lesson.TeacherID {
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
