package get_times

import (
	"context"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	GetTeacherById(ctx context.Context, teacherId int) (*entities.Teacher, error)
	GetTeacherByUserId(ctx context.Context, userId int) (*entities.Teacher, error)
	GetScheduleTimesByTeacherId(ctx context.Context, Id int) ([]*entities.ScheduleTime, error)
}
