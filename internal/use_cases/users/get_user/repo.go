package get_user

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	GetUserById(ctx context.Context, id int) (*entities.User, error)
	IsTeacherExistsByUserId(ctx context.Context, id int) (bool, error)
}
