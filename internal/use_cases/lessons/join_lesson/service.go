package join_lesson

import (
	"context"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

type MeetService interface {
	GenerateMeetingToken(roomName string) (string, error)
	NameRoomByLessonId(lessonId int) string
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

	// is user exists
	exists, err := s.repo.IsUserExistsById(ctx, userId)
	if err != nil {
		return "", fmt.Errorf("failed to check user existstance by id: %w", err)
	}
	if !exists {
		return "", serviceErrs.ErrorUserNotFound
	}

	// is lesson exists
	exists, err = s.repo.IsLessonExistsById(ctx, lessonId)
	if err != nil {
		return "", fmt.Errorf("failed to check lesson existstance by id: %w", err)
	}
	if !exists {
		return "", serviceErrs.ErrorLessonNotFound
	}

	// get lesson
	lesson, err := s.repo.GetLessonById(ctx, lessonId)
	if err != nil {
		return "", fmt.Errorf("failed to get lesson by id: %w", err)
	}

	// if user is not student => maybe he is a teacher (check it)
	if lesson.StudentId != userId {
		// is teacher exists by userId
		exists, err = s.repo.IsTeacherExistsByUserId(ctx, userId)
		if err != nil {
			return "", fmt.Errorf("failed to check teacher existstance by userId: %w", err)
		}

		// if it's not the teacher for this lesson either => error
		if !exists {
			return "", serviceErrs.ErrorNotRelatedUserToLesson
		}

		teacher, err := s.repo.GetTeacherByUserId(ctx, userId)
		if err != nil {
			return "", fmt.Errorf("failed to get teacher by userId: %w", err)
		}
		if lesson.TeacherId != teacher.Id {
			return "", serviceErrs.ErrorNotRelatedUserToLesson
		}
	}

	// get ongoing statusId
	ongoingStatusId, err := s.repo.GetStatusIdByStatusName(ctx, entities.OngoingStatusName)
	if err != nil {
		return "", fmt.Errorf("failed to get status by status name: %w", err)
	}

	// may join lesson if only was ongoing status
	if lesson.StatusId != ongoingStatusId {
		return "", serviceErrs.ErrorStatusNonOngoing
	}

	token, err := s.meetService.GenerateMeetingToken(s.meetService.NameRoomByLessonId(lessonId))
	if err != nil {
		return "", fmt.Errorf("failed to generate meeting token: %w", err)
	}

	return token, nil
}
