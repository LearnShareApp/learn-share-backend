package lesson

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

// PlanLesson set lesson in planned state.
func (s *LessonService) PlanLesson(ctx context.Context, userID int, lessonID int) error {
	return s.changeLessonStateAsTeacher(ctx, userID, lessonID, entities.Planned)
}
