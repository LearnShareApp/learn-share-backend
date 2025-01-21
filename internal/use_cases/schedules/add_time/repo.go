package add_time

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"time"
)

type repo interface {
	IsTeacherExistsByUserId(ctx context.Context, id int) (bool, error)
	GetTeacherByUserId(ctx context.Context, userId int) (*entities.Teacher, error)
	IsTimeExistsByTeacherIdAndDatetime(ctx context.Context, id int, datetime time.Time) (bool, error)
	CreateTime(ctx context.Context, teacherId int, datetime time.Time) error
}
