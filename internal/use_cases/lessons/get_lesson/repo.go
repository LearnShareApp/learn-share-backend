package get_lesson

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	GetLessonById(ctx context.Context, id int) (*entities.Lesson, error)
	GetUserById(ctx context.Context, id int) (*entities.User, error)
}
