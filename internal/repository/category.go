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

func (r *Repository) IsCategoryExistsById(ctx context.Context, id int) (bool, error) {
	const req = `SELECT EXISTS(SELECT 1 FROM public.categories WHERE category_id = $1)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, req, id)

	if err != nil {
		return false, fmt.Errorf("failed to check category existence: %w", err)
	}

	return exists, nil
}
