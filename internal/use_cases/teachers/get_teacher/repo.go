package get_teacher

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	IsTeacherExistsByUserId(ctx context.Context, id int) (bool, error)
	GetUserById(ctx context.Context, id int) (*entities.User, error)
}
