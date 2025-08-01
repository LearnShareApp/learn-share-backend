package lesson

import (
	"context"
	"errors"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

func (s *LessonService) BookLesson(ctx context.Context, lesson *entities.Lesson) error {
	if err := s.validateUserExists(ctx, lesson.StudentID); err != nil {
		return err
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
	exists, err := s.repo.IsCategoryExistsByID(ctx, lesson.CategoryID)
	if err != nil {
		return fmt.Errorf("failed to check categories existstance by id: %w", err)
	}
	if !exists {
		return serviceErrs.ErrorCategoryNotFound
	}

	// is teacher have such ACTIVE skill
	skill, err := s.repo.GetSkillByTeacherIDAndCategoryID(ctx, lesson.TeacherID, lesson.CategoryID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorSkillUnregistered
		}
		return fmt.Errorf("failed to get skill by teacher and category: %w", err)
	}
	if !skill.IsActive {
		return serviceErrs.ErrorSkillInactive
	}

	// is this time still available and owner is this teacher
	scheduleTime, err := s.repo.GetScheduleTimeByID(ctx, lesson.ScheduleTimeID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorScheduleTimeNotFound
		}

		return fmt.Errorf("failed to get schedules time by id: %w", err)
	}

	if scheduleTime.TeacherID != lesson.TeacherID {
		return serviceErrs.ErrorScheduleTimeForAnotherTeacher
	}

	if !scheduleTime.IsAvailable {
		return serviceErrs.ErrorScheduleTimeUnavailable
	}

	// book lesson
	if err = s.repo.BookLesson(ctx,
		lesson.ScheduleTimeID,
		lesson.StudentID,
		lesson.TeacherID,
		lesson.CategoryID); err != nil {
		// if some another booked faster between check and upd
		if errors.Is(err, serviceErrs.ErrorNonUniqueData) {
			return serviceErrs.ErrorLessonTimeBooked
		}

		return fmt.Errorf("failed to book lesson: %w", err)
	}

	return nil
}
