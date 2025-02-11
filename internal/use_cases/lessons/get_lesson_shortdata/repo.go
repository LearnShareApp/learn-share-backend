package get_lesson_shortdata

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	GetLessonById(ctx context.Context, id int) (*entities.Lesson, error)
	GetTeacherById(ctx context.Context, id int) (*entities.Teacher, error)
}
