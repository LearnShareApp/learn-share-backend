package get_teachers

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	GetAllTeachersData(ctx context.Context) ([]entities.User, error)
}
