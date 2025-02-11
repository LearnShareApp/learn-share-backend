package get_teacher_lessons

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	IsUserExistsById(ctx context.Context, id int) (bool, error)
	GetTeacherByUserId(ctx context.Context, id int) (*entities.Teacher, error)
	GetTeacherLessonsByTeacherId(ctx context.Context, id int) ([]*entities.Lesson, error)
}
