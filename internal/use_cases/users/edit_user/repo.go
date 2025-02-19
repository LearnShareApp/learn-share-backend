package edit_user

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	GetUserById(ctx context.Context, id int) (*entities.User, error)
	UpdateUser(ctx context.Context, userId int, user *entities.User) error
}
