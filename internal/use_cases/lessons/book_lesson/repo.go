package book_lesson

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	IsUserExistsById(ctx context.Context, id int) (bool, error)
	GetTeacherById(ctx context.Context, id int) (*entities.Teacher, error)
	IsCategoryExistsById(ctx context.Context, id int) (bool, error)
	IsSkillExistsByTeacherIdAndCategoryId(ctx context.Context, teacherId int, categoryId int) (bool, error)
	GetScheduleTimeById(ctx context.Context, id int) (*entities.ScheduleTime, error)
	CreateUnconfirmedLesson(ctx context.Context, lesson *entities.Lesson) error
}
