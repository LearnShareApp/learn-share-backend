package login

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
}
