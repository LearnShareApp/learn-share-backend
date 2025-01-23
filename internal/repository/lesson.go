package repository

import (
	"context"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	internalErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/lib/pq"
)

func (r *Repository) CreateUnconfirmedLesson(ctx context.Context, lesson *entities.Lesson) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	defer tx.Rollback()

	const query = `
	INSERT INTO lessons (student_id, teacher_id, category_id, schedule_time_id) 
	VALUES ($1, $2, $3, $4)
	`

	// insert lesson
	if _, err := tx.ExecContext(ctx, query,
		lesson.StudentId,
		lesson.TeacherId,
		lesson.CategoryId,
		lesson.ScheduleTimeId); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// Код ошибки 23505 означает unique_violation
			if pqErr.Code == "23505" {
				return internalErrs.ErrorNonUniqueData
			}
		}
		return fmt.Errorf("failed to insert lesson: %w", err)
	}

	// book time
	if err = bookScheduleTime(ctx, tx, lesson.ScheduleTimeId); err != nil {
		return fmt.Errorf("failed to book schedule time: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
