package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"github.com/LearnShareApp/learn-share-backend/internal/errors"
)

func (r *Repository) IsUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM public.users WHERE email = $1)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, email)

	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

func (r *Repository) IsUserExistsById(ctx context.Context, id int) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM public.users WHERE user_id = $1)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, id)

	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *entities.User) (int, error) {
	const query = `
	INSERT INTO users (email, password, name, surname, birthdate) 
	VALUES ($1, $2, $3, $4, $5)
	RETURNING user_id
	`

	var userID int
	if err := r.db.QueryRowContext(ctx, query, user.Email, user.Password, user.Name, user.Surname, user.Birthdate).Scan(&userID); err != nil {
		return 0, err
	}
	return userID, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	const query = `SELECT user_id, email, password, name, surname, birthdate FROM public.users WHERE email = $1`

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

func (r *Repository) GetUserById(ctx context.Context, id int) (*entities.User, error) {
	const query = `SELECT user_id, email, password, name, surname, birthdate, registration_date FROM public.users WHERE user_id = $1`

	var user entities.User
	err := r.db.GetContext(ctx, &user, query, id)

	if err == sql.ErrNoRows {
		return nil, errors.ErrorSelectEmpty
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return &user, nil
}
