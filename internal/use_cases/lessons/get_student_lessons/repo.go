package get_student_lessons

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	IsUserExistsById(ctx context.Context, id int) (bool, error)
	GetStudentLessonsByUserId(ctx context.Context, id int) ([]*entities.Lesson, error)
}
