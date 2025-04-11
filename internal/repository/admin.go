package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	internalErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/Masterminds/squirrel"
)

func (r *Repository) IsUserAdminByID(ctx context.Context, id int) (bool, error) {
	query, args, err := r.sqlBuilder.
		Select("is_admin").
		From("users").
		Where(squirrel.Eq{"user_id": id}).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	var isAdmin bool

	err = r.db.GetContext(ctx, &isAdmin, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, internalErrs.ErrorSelectEmpty
		}

		return false, fmt.Errorf("failed to check is user an admin: %w", err)
	}

	return isAdmin, nil
}
