package get_times

import (
	"context"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	GetTeacherById(ctx context.Context, teacherID int) (*entities.Teacher, error)
	GetTeacherByUserId(ctx context.Context, userID int) (*entities.Teacher, error)
	GetScheduleTimesByTeacherId(ctx context.Context, id int) ([]*entities.ScheduleTime, error)
}
