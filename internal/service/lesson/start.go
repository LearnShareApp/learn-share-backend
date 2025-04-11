package lesson

import (
	"context"
	"errors"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

func (s *LessonService) StartLesson(ctx context.Context, userID, lessonID int) (string, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return "", serviceErrs.ErrorUserNotFound
		}

		return "", fmt.Errorf("failed to check user existstance by id: %w", err)
	}

	// get teacher
	teacher, err := s.repo.GetTeacherByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return "", serviceErrs.ErrorUserIsNotTeacher
		}

		return "", fmt.Errorf("failed to get teacher by userID: %w", err)
	}

	// get lesson
	lesson, err := s.repo.GetLessonByID(ctx, lessonID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return "", serviceErrs.ErrorLessonNotFound
		}

		return "", fmt.Errorf("failed to get lesson by id: %w", err)
	}

	if lesson.TeacherID != teacher.ID {
		return "", serviceErrs.ErrorNotRelatedTeacherToLesson
	}

	// get waiting statusId
	waitingStatusId, err := s.repo.GetStatusIDByStatusName(ctx, entities.WaitingStatusName)
	if err != nil {
		return "", fmt.Errorf("failed to get status by status name: %w", err)
	}

	// may start lesson if only was waiting status
	if lesson.StatusID != waitingStatusId {
		return "", serviceErrs.ErrorStatusNonWaiting
	}

	// get ongoing statusId
	ongoingStatusId, err := s.repo.GetStatusIDByStatusName(ctx, entities.OngoingStatusName)
	if err != nil {
		return "", fmt.Errorf("failed to get status by status name: %w", err)
	}

	token, err := s.meetCreator.GenerateMeetingToken(s.meetCreator.NameRoomByLessonID(lessonID),
		s.meetCreator.GetUserIdentityString(user.Name, user.Surname, user.ID))
	if err != nil {
		return "", fmt.Errorf("failed to generate meeting token: %w", err)
	}

	// save token and edit status
	if err = s.repo.ChangeLessonStatus(ctx, lessonID, ongoingStatusId); err != nil {
		return "", fmt.Errorf("failed to edit lesson status: %w", err)
	}

	return token, nil
}
