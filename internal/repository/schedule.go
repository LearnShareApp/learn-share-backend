package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"time"
)

func (r *Repository) IsTimeExistsByTeacherIdAndDatetime(ctx context.Context, id int, datetime time.Time) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM public.schedule_times WHERE teacher_id = $1 AND datetime = $2)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, id, datetime)

	if err != nil {
		return false, fmt.Errorf("failed to check schadule_time existence: %w", err)
	}

	return exists, nil
}

func (r *Repository) CreateTime(ctx context.Context, teacherId int, datetime time.Time) error {
	const query = `
	INSERT INTO schedule_times (teacher_id, datetime) 
	VALUES ($1, $2)
	`

	if _, err := r.db.ExecContext(ctx, query, teacherId, datetime); err != nil {
		return fmt.Errorf("failed to insert schedule time: %w", err)
	}
	return nil
}

func (r *Repository) GetAvailableScheduleTimesByTeacherId(ctx context.Context, id int) ([]*entities.ScheduleTime, error) {
	const query = `
		SELECT schedule_time_id, teacher_id, datetime FROM schedule_times WHERE teacher_id = $1 AND is_available
		`

	var times []*entities.ScheduleTime
	err := r.db.SelectContext(ctx, &times, query, id)
	if err != nil {
		// empty times isn't error
		if errors.Is(err, sql.ErrNoRows) {
			return times, nil
		}
		return nil, fmt.Errorf("failed to find schedule times: %w", err)
	}

	return times, nil
}
