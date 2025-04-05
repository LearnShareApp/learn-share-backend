package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"

	"github.com/Masterminds/squirrel"
)

func (r *Repository) GetCategories(ctx context.Context) ([]*entities.Category, error) {
	query, args, err := r.sqlBuilder.Select("category_id", "name", "min_age").
		From("categories").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var categories []*entities.Category

	err = r.db.SelectContext(ctx, &categories, query, args...)
	if err != nil {
		// empty categories isn't error
		if errors.Is(err, sql.ErrNoRows) {
			return categories, nil
		}

		return nil, fmt.Errorf("failed to find categories: %w", err)
	}

	return categories, nil
}

func (r *Repository) IsCategoryExistsByID(ctx context.Context, id int) (bool, error) {
	query, args, err := r.sqlBuilder.
		Select("1").
		From("categories").
		Where(squirrel.Eq{"category_id": id}).
		Prefix("SELECT EXISTS (").
		Suffix(")").
		ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	var exists bool

	err = r.db.GetContext(ctx, &exists, query, args...)
	if err != nil {
		return false, fmt.Errorf("failed to check category existence: %w", err)
	}

	return exists, nil
}
