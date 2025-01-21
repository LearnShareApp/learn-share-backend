package repository

import (
	"context"
	"fmt"
	"time"
)

func (r *Repository) IsTimeExistsByTeacherIdAndDatetime(ctx context.Context, id int, datetime time.Time) (bool, error) {
	const req = `SELECT EXISTS(SELECT 1 FROM public.schedule_times WHERE teacher_id = $1 AND datetime = $2)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, req, id, datetime)

	if err != nil {
		return false, fmt.Errorf("failed to check schadule_time existence: %w", err)
	}

	return exists, nil
}

func (r *Repository) CreateTime(ctx context.Context, teacherId int, datetime time.Time) error {
	const req = `
	INSERT INTO schedule_times (teacher_id, datetime) 
	VALUES ($1, $2)
	`

	if _, err := r.db.ExecContext(ctx, req, teacherId, datetime); err != nil {
		return fmt.Errorf("failed to insert schedule time: %w", err)
	}
	return nil
}
