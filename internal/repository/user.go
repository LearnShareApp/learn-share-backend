package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

func (r *Repository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	const req = `SELECT EXISTS(SELECT 1 FROM public.users WHERE email = $1)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, req, email)

	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *entities.User) (int64, error) {
	const req = `
	INSERT INTO users (email, password) 
	VALUES ($1, $2)
	RETURNING id
	`

	var userID int64
	if err := r.db.QueryRowContext(ctx, req, user.Email, user.Password).Scan(&userID); err != nil {
		return 0, err
	}
	return userID, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `SELECT id, email, password FROM public.users WHERE email = $1`

	var user entities.User
	err := r.db.GetContext(ctx, &user, query, email)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return &user, nil
}
