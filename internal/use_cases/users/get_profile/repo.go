package get_profile

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	IsUserExistsById(ctx context.Context, id int) (bool, error)
	GetUserById(ctx context.Context, id int) (*entities.User, error)
	IsTeacherExistsByUserId(ctx context.Context, id int) (bool, error)
}
