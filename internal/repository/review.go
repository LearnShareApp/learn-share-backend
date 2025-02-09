package repository

import (
	"context"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/lib/pq"
)

func (r *Repository) CreateReview(ctx context.Context, review *entities.Review) error {
	const query = `
	INSERT INTO reviews (teacher_id, student_id, category_id, skill_id, rate, comment)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	if _, err := r.db.ExecContext(ctx, query,
		review.TeacherId,
		review.StudentId,
		review.CategoryId,
		review.SkillId,
		review.Rate,
		review.Comment); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// Код ошибки 23505 означает unique_violation
			if pqErr.Code == "23505" {
				return errors.ErrorNonUniqueData
			}
		}

		return fmt.Errorf("failed to insert review: %w", err)
	}
	return nil
}
