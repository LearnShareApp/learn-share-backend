package get_times

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	IsUserExistsById(ctx context.Context, id int) (bool, error)
	IsTeacherExistsByUserId(ctx context.Context, id int) (bool, error)
	GetTeacherByUserId(ctx context.Context, userId int) (*entities.Teacher, error)
	GetScheduleTimesByTeacherId(ctx context.Context, Id int) ([]*entities.ScheduleTime, error)
}
