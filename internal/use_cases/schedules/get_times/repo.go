package get_times

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	IsTeacherExistsByUserId(ctx context.Context, id int) (bool, error)
	GetTeacherByUserId(ctx context.Context, userId int) (*entities.Teacher, error)
	GetAvailableScheduleTimesByTeacherId(ctx context.Context, Id int) ([]*entities.ScheduleTime, error)
}
