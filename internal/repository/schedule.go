package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	internalErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"

	"github.com/jmoiron/sqlx"
)

func (r *Repository) IsScheduleTimeExistsByTeacherIDAndDatetime(ctx context.Context, id int, datetime time.Time) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM schedule_times WHERE teacher_id = $1 AND datetime = $2)`

	var exists bool

	err := r.db.GetContext(ctx, &exists, query, id, datetime)
	if err != nil {
		return false, fmt.Errorf("failed to check schadule_time existence by teacherID and time: %w", err)
	}

	return exists, nil
}

func (r *Repository) IsScheduleTimeExistsById(ctx context.Context, id int) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM schedule_times WHERE schedule_time_id = $1)`

	var exists bool

	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check schadule_time existence by scheduleTimeID: %w", err)
	}

	return exists, nil
}

func (r *Repository) CreateScheduleTime(ctx context.Context, teacherId int, datetime time.Time) error {
	const query = `
	INSERT INTO schedule_times (teacher_id, datetime) 
	VALUES ($1, $2)
	`

	if _, err := r.db.ExecContext(ctx, query, teacherId, datetime); err != nil {
		return fmt.Errorf("failed to insert schedule time: %w", err)
	}

	return nil
}

func (r *Repository) GetScheduleTimeByID(ctx context.Context, id int) (*entities.ScheduleTime, error) {
	const query = `SELECT schedule_time_id, teacher_id, datetime, is_available FROM schedule_times WHERE schedule_time_id = $1`

	var scheduleTime entities.ScheduleTime

	err := r.db.GetContext(ctx, &scheduleTime, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrs.ErrorSelectEmpty
		}

		return nil, fmt.Errorf("failed to get schadule_time existence by scheduleTimeID: %w", err)
	}

	return &scheduleTime, nil
}

func (r *Repository) GetScheduleTimesByTeacherID(ctx context.Context, id int) ([]*entities.ScheduleTime, error) {
	const query = `
		SELECT schedule_time_id, teacher_id, datetime, is_available FROM schedule_times 
		WHERE teacher_id = $1 AND 
		      datetime >= NOW()
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

func bookScheduleTime(ctx context.Context, tx *sqlx.Tx, id int) error {
	const query = `
	UPDATE schedule_times SET is_available = false WHERE schedule_time_id = $1
	`

	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to update schedule_times table: %w", err)
	}

	return nil
}
