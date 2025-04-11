package lesson

import (
	"context"
	"errors"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

func (s *LessonService) JoinLesson(ctx context.Context, userID int, lessonID int) (string, error) {
	// get user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return "", serviceErrs.ErrorUserNotFound
		}

		return "", fmt.Errorf("failed to check user existstance by id: %w", err)
	}

	// get lesson
	lesson, err := s.repo.GetLessonByID(ctx, lessonID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return "", serviceErrs.ErrorLessonNotFound
		}

		return "", fmt.Errorf("failed to get lesson by id: %w", err)
	}

	// if user is not student => maybe he is a teacher (check it)
	if lesson.StudentID != userID {
		// is teacher exists by userID
		// if it's not the teacher for this lesson either => error
		teacher, err := s.repo.GetTeacherByUserID(ctx, userID)
		if err != nil {
			if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
				return "", serviceErrs.ErrorNotRelatedUserToLesson
			}

			return "", fmt.Errorf("failed to get teacher by userID: %w", err)
		}

		if lesson.TeacherID != teacher.ID {
			return "", serviceErrs.ErrorNotRelatedUserToLesson
		}
	}

	// get ongoing statusId
	ongoingStatusId, err := s.repo.GetStatusIDByStatusName(ctx, entities.OngoingStatusName)
	if err != nil {
		return "", fmt.Errorf("failed to get status by status name: %w", err)
	}

	// may join lesson if only was ongoing status
	if lesson.StatusID != ongoingStatusId {
		return "", serviceErrs.ErrorStatusNonOngoing
	}

	token, err := s.meetCreator.GenerateMeetingToken(s.meetCreator.NameRoomByLessonID(lessonID),
		s.meetCreator.GetUserIdentityString(user.Name, user.Surname, user.ID))
	if err != nil {
		return "", fmt.Errorf("failed to generate meeting token: %w", err)
	}

	return token, nil
}
