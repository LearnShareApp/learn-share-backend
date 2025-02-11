package approve_lesson

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

func (s *Service) Do(ctx context.Context, userId int, lessonId int) error {
	// is user exists
	exists, err := s.repo.IsUserExistsById(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to check user existstance by id: %w", err)
	}
	if !exists {
		return serviceErrs.ErrorUserNotFound
	}

	// get teacher
	teacher, err := s.repo.GetTeacherByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorUserIsNotTeacher
		}
		return fmt.Errorf("failed to get teacher by userId: %w", err)
	}

	// get lesson
	lesson, err := s.repo.GetLessonById(ctx, lessonId)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorLessonNotFound
		}
		return fmt.Errorf("failed to get lesson by id: %w", err)
	}

	if lesson.TeacherId != teacher.Id {
		return serviceErrs.ErrorNotRelatedTeacherToLesson
	}

	// get verification statusId
	verificationStatusId, err := s.repo.GetStatusIdByStatusName(ctx, entities.VerificationStatusName)
	if err != nil {
		return fmt.Errorf("failed to get status by status name: %w", err)
	}

	// may approve lesson if only was verification status
	if lesson.StatusId != verificationStatusId {
		return serviceErrs.ErrorStatusNonVerification
	}

	// get waiting statusId
	waitingStatusId, err := s.repo.GetStatusIdByStatusName(ctx, entities.WaitingStatusName)
	if err != nil {
		return fmt.Errorf("failed to get status by status name: %w", err)
	}
	// change lesson status
	if err = s.repo.ChangeLessonStatus(ctx, lessonId, waitingStatusId); err != nil {
		return fmt.Errorf("failed to change lesson status: %w", err)
	}

	return nil
}
