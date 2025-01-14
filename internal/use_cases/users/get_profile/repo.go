package get_profile

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	ExistsById(ctx context.Context, id int64) (bool, error)
	GetUserById(ctx context.Context, id int64) (*entities.User, error)
}
