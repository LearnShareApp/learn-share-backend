package lesson

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

func (s *LessonService) FinishLesson(ctx context.Context, userID int, lessonID int) error {
	return s.changeLessonStateAsTeacher(ctx, userID, lessonID, entities.Finished)

}
