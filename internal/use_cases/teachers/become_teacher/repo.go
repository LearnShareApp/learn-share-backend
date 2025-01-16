package become_teacher

import (
	"context"
)

type repo interface {
	IsUserExistsById(ctx context.Context, id int) (bool, error)
	IsTeacherExistsByUserId(ctx context.Context, id int) (bool, error)
	CreateTeacher(ctx context.Context, userId int) error
}
