package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repository) IsTeacherExistsByUserId(ctx context.Context, id int) (bool, error) {
	const req = `SELECT EXISTS(SELECT 1 FROM public.teachers WHERE teacher_id = $1)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, req, id)

	if err != nil {
		return false, fmt.Errorf("failed to check teacher existence: %w", err)
	}

	return exists, nil
}

func (r *Repository) CreateTeacher(ctx context.Context, userId int) error {
	const req = `
	INSERT INTO teachers (user_id) 
	VALUES ($1)
	`

	if _, err := r.db.ExecContext(ctx, req, userId); err != nil {
		return fmt.Errorf("failed to insert teacher: %w", err)
	}
	return nil
}

func (r *Repository) CreateTeacherIfNotExists(ctx context.Context, userId int) (int, error) {
	const (
		selectQuery = `
		SELECT teacher_id FROM teachers WHERE user_id = $1
		`

		insertQuery = `
		INSERT INTO teachers (user_id) 
		VALUES ($1)
		RETURNING teacher_id
		`
	)

	var teacherId int

	err := r.db.GetContext(ctx, &teacherId, selectQuery, userId)
	if err == nil {
		return teacherId, nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		// Если ошибка не sql.ErrNoRows, возвращаем её
		return 0, fmt.Errorf("failed to select teacher: %w", err)
	}

	if err := r.db.QueryRowContext(ctx, insertQuery, userId).Scan(&teacherId); err != nil {
		return 0, fmt.Errorf("failed to insert teacher: %w", err)
	}
	return teacherId, nil
}
