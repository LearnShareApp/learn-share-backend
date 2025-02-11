package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	internalErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
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
	INSERT INTO users (email, password, name, surname, birthdate, avatar) 
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING user_id
	`

	var userID int
	if err := r.db.QueryRowContext(ctx, query, user.Email, user.Password, user.Name, user.Surname, user.Birthdate, user.Avatar).Scan(&userID); err != nil {
		return 0, err
	}
	return userID, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	const query = `SELECT user_id, email, password, name, surname, birthdate FROM public.users WHERE email = $1`

	var user entities.User
	err := r.db.GetContext(ctx, &user, query, email)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, internalErrs.ErrorSelectEmpty
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return &user, nil
}

func (r *Repository) GetUserById(ctx context.Context, id int) (*entities.User, error) {
	const query = `SELECT user_id, email, password, name, surname, registration_date, birthdate, avatar FROM public.users WHERE user_id = $1`

	var user entities.User
	err := r.db.GetContext(ctx, &user, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrs.ErrorSelectEmpty
		}

		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return &user, nil
}

func (r *Repository) GetUserStatByUserId(ctx context.Context, id int) (*entities.StudentStatistic, error) {
	const query = `
    SELECT 
        COUNT(DISTINCT CASE WHEN s.name = $1 THEN l.lesson_id END) as count_of_finished_lesson,
        COUNT(DISTINCT CASE WHEN s.name = $2 THEN l.lesson_id END) as count_of_verification_lesson,
        COUNT(DISTINCT CASE WHEN s.name = $3 THEN l.lesson_id END) as count_of_waiting_lesson,
        COUNT(DISTINCT CASE WHEN s.name = $1 THEN l.teacher_id END) as count_of_teachers
    FROM lessons l
    INNER JOIN statuses s ON s.status_id = l.status_id
    WHERE l.student_id = $4
    `

	var stat entities.StudentStatistic
	err := r.db.GetContext(ctx, &stat, query,
		entities.FinishedStatusName,     // $1
		entities.VerificationStatusName, // $2
		entities.WaitingStatusName,      // $3
		id,                              // $4
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrs.ErrorSelectEmpty
		}

		return nil, fmt.Errorf("failed to find user's statistic by user id: %w", err)
	}

	return &stat, nil
}

func (r *Repository) UpdateUser(ctx context.Context, userId int, user *entities.User) error {
	const query = `
	UPDATE users SET (name, surname, password, birthdate, avatar) = ($2, $3, $4, $5, $6) WHERE user_id = $1
	`

	if _, err := r.db.ExecContext(ctx, query, userId,
		user.Name,
		user.Surname,
		user.Password,
		user.Birthdate,
		user.Avatar); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
