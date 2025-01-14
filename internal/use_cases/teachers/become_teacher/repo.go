package become_teacher

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	IsUserExistsById(ctx context.Context, id int64) (bool, error)
	IsTeacherExistsByUserId(ctx context.Context, id int64) (bool, error)
	CreateTeacher(ctx context.Context, teacher *entities.Teacher) error
}