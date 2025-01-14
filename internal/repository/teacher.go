package repository

import (
	"context"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

func (r *Repository) IsTeacherExistsByUserId(ctx context.Context, id int64) (bool, error) {
	const req = `SELECT EXISTS(SELECT 1 FROM public.teachers WHERE teacher_id = $1)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, req, id)

	if err != nil {
		return false, fmt.Errorf("failed to check teacher existence: %w", err)
	}

	return exists, nil
}

func (r *Repository) CreateTeacher(ctx context.Context, teacher *entities.Teacher) error {
	const req = `
	INSERT INTO teachers (user_id) 
	VALUES ($1)
	`

	if teacher == nil {
		return fmt.Errorf("teacher is nil")
	}

	if _, err := r.db.ExecContext(ctx, req, teacher.UserId); err != nil {
		return fmt.Errorf("failed to insert teacher: %w", err)
	}
	return nil
}
