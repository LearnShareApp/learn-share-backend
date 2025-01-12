package get_categories

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	GetCategories(ctx context.Context) ([]*entities.Category, error)
}
