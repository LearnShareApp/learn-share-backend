package lesson

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

// RejectLesson set lesson in planned state.
func (s *LessonService) RejectLesson(ctx context.Context, userID int, lessonID int) error {
	return s.changeLessonStateAsTeacher(ctx, userID, lessonID, entities.Rejected)
}
