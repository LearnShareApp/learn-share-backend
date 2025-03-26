package start_lesson

import (
	"context"
	"errors"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	internalErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

type MeetService interface {
	GenerateMeetingToken(roomName string, userName string) (string, error)
	NameRoomByLessonID(lessonId int) string
	GetUserIdentityString(userName, userSurname string, id int) string
}

type Service struct {
	repo        repo
	meetService MeetService
}

func NewService(repo repo, meetService MeetService) *Service {
	return &Service{
		repo:        repo,
		meetService: meetService,
	}
}

func (s *Service) Do(ctx context.Context, userId int, lessonId int) (string, error) {
	user, err := s.repo.GetUserById(ctx, userId)
	if err != nil {
		if errors.Is(err, internalErrs.ErrorSelectEmpty) {
			return "", internalErrs.ErrorUserNotFound
		}

		return "", fmt.Errorf("failed to check user existstance by id: %w", err)
	}

	// get teacher
	teacher, err := s.repo.GetTeacherByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, internalErrs.ErrorSelectEmpty) {
			return "", internalErrs.ErrorUserIsNotTeacher
		}
		return "", fmt.Errorf("failed to get teacher by userId: %w", err)
	}

	// get lesson
	lesson, err := s.repo.GetLessonById(ctx, lessonId)
	if err != nil {
		if errors.Is(err, internalErrs.ErrorSelectEmpty) {
			return "", internalErrs.ErrorLessonNotFound
		}
		return "", fmt.Errorf("failed to get lesson by id: %w", err)
	}

	if lesson.TeacherId != teacher.Id {
		return "", internalErrs.ErrorNotRelatedTeacherToLesson
	}

	// get waiting statusId
	waitingStatusId, err := s.repo.GetStatusIdByStatusName(ctx, entities.WaitingStatusName)
	if err != nil {
		return "", fmt.Errorf("failed to get status by status name: %w", err)
	}

	// may start lesson if only was waiting status
	if lesson.StatusId != waitingStatusId {
		return "", internalErrs.ErrorStatusNonWaiting
	}

	// get ongoing statusId
	ongoingStatusId, err := s.repo.GetStatusIdByStatusName(ctx, entities.OngoingStatusName)
	if err != nil {
		return "", fmt.Errorf("failed to get status by status name: %w", err)
	}

	token, err := s.meetService.GenerateMeetingToken(s.meetService.NameRoomByLessonID(lessonId),
		s.meetService.GetUserIdentityString(user.Name, user.Surname, user.Id))
	if err != nil {
		return "", fmt.Errorf("failed to generate meeting token: %w", err)
	}

	// save token and edit status
	if err = s.repo.ChangeLessonStatus(ctx, lessonId, ongoingStatusId); err != nil {
		return "", fmt.Errorf("failed to edit lesson status: %w", err)
	}

	return token, nil
}
