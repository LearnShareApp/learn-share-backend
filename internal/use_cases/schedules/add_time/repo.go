package add_time

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"time"
)

type repo interface {
	IsUserExistsById(ctx context.Context, id int) (bool, error)
	IsTeacherExistsByUserId(ctx context.Context, id int) (bool, error)
	GetTeacherByUserId(ctx context.Context, userId int) (*entities.Teacher, error)
	IsScheduleTimeExistsByTeacherIdAndDatetime(ctx context.Context, id int, datetime time.Time) (bool, error)
	CreateScheduleTime(ctx context.Context, teacherId int, datetime time.Time) error
}
