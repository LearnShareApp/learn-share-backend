package start_lesson

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	IsUserExistsById(ctx context.Context, id int) (bool, error)
	IsTeacherExistsByUserId(ctx context.Context, id int) (bool, error)
	IsLessonExistsById(ctx context.Context, id int) (bool, error)
	GetTeacherByUserId(ctx context.Context, id int) (*entities.Teacher, error)
	GetLessonById(ctx context.Context, id int) (*entities.Lesson, error)
	GetStatusIdByStatusName(ctx context.Context, name string) (int, error)
	EditStatusAndSaveTokenInLesson(ctx context.Context, lessonId int, statusId int, token string) error
}
