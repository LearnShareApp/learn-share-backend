package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	internalErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
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
				return internalErrs.ErrorNonUniqueData
			}
		}

		return fmt.Errorf("failed to insert review: %w", err)
	}
	return nil
}

func (r *Repository) GetReviewsByTeacherId(ctx context.Context, id int) ([]*entities.Review, error) {
	const query = `
    SELECT
		r.review_id,
		r.teacher_id,
		r.student_id,
		r.category_id,
		r.skill_id,
		r.rate,
		r.comment,
		
		u.user_id,
		u.email,
		u.name,
		u.surname,
		u.avatar
	FROM reviews r
    INNER JOIN users u ON r.student_id = u.user_id
    WHERE r.teacher_id = $1`

	// temp struct for executed data
	type result struct {
		entities.Review
		entities.User
	}

	var rows []result
	err := r.db.SelectContext(ctx, &rows, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// not an error
			return make([]*entities.Review, 0), nil
		}
		return nil, fmt.Errorf("failed to extract reviews: %w", err)
	}

	// Map for grope results
	reviewsMap := make(map[int]*entities.Review)

	for _, row := range rows {
		review, exists := reviewsMap[row.Review.Id]
		if !exists {
			review = &row.Review
			if review.StudentData == nil {
				review.StudentData = &row.User
			}
			reviewsMap[row.Review.Id] = review
		}
	}

	reviews := make([]*entities.Review, 0, len(reviewsMap))
	for _, review := range reviewsMap {
		reviews = append(reviews, review)
	}

	return reviews, nil
}
