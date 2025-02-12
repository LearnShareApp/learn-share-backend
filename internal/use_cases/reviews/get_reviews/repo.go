package get_reviews

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	IsTeacherExistsById(ctx context.Context, id int) (bool, error)
	GetReviewsByTeacherId(ctx context.Context, id int) ([]*entities.Review, error)
}
