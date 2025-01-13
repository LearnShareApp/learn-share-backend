package repository

import (
	"context"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

func (r *Repository) GetCategories(ctx context.Context) ([]*entities.Category, error) {
	query := `SELECT category_id, name, min_age FROM public.categories`

	var categories []*entities.Category
	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find categories: %w", err)
	}

	return categories, nil
}
