package cancel_lesson

import (
	"context"
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

func (s *Service) Do(ctx context.Context, userId int, lessonId int) error {
	// is user exists
	exists, err := s.repo.IsUserExistsById(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to check user existstance by id: %w", err)
	}
	if !exists {
		return serviceErrs.ErrorUserNotFound
	}

	// is lesson exists
	exists, err = s.repo.IsLessonExistsById(ctx, lessonId)
	if err != nil {
		return fmt.Errorf("failed to check lesson existstance by id: %w", err)
	}
	if !exists {
		return serviceErrs.ErrorLessonNotFound
	}

	// get lesson
	lesson, err := s.repo.GetLessonById(ctx, lessonId)
	if err != nil {
		return fmt.Errorf("failed to get lesson by id: %w", err)
	}

	// if user is not student => maybe he is a teacher (check it)
	if lesson.StudentId != userId {
		// is teacher exists by userId
		exists, err = s.repo.IsTeacherExistsByUserId(ctx, userId)
		if err != nil {
			return fmt.Errorf("failed to check teacher existstance by userId: %w", err)
		}

		// if it's not the teacher for this lesson either => error
		if !exists {
			return serviceErrs.ErrorNotRelatedUserToLesson
		}

		teacher, err := s.repo.GetTeacherByUserId(ctx, userId)
		if err != nil {
			return fmt.Errorf("failed to get teacher by userId: %w", err)
		}
		if lesson.TeacherId != teacher.Id {
			return serviceErrs.ErrorNotRelatedUserToLesson
		}
	}

	// if we here => related user to lesson
	// get cancel statusId
	statusId, err := s.repo.GetStatusIdByStatusName(ctx, entities.CancelStatusName)
	if err != nil {
		return fmt.Errorf("failed to get status by status name: %w", err)
	}

	// change lesson status
	if err = s.repo.ChangeLessonStatus(ctx, lessonId, statusId); err != nil {
		return fmt.Errorf("failed to change lesson status: %w", err)
	}

	return nil
}
