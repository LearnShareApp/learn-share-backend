package lesson

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

// StartLesson start lesson and returns meet token
func (s *LessonService) StartLesson(ctx context.Context, userID, lessonID int) (string, error) {
	err := s.changeLessonStateAsTeacher(ctx, userID, lessonID, entities.OngoingStatusName)
	if err != nil {
		return "", err
	}

	return s.generateLessonMeetingToken(ctx, userID, lessonID)
}
