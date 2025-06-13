package lesson

import (
	"context"
	"fmt"
)

// generateLessonMeetingToken generate token for lesson
// only for local usage
func (s *LessonService) generateLessonMeetingToken(ctx context.Context, userID, lessonID int) (string, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user for token generation: %w", err)
	}

	token, err := s.meetCreator.GenerateMeetingToken(
		s.meetCreator.NameRoomByLessonID(lessonID),
		s.meetCreator.GetUserIdentityString(user.Name, user.Surname, user.ID))
	if err != nil {
		return "", fmt.Errorf("failed to generate meeting token: %w", err)
	}

	return token, nil
}
