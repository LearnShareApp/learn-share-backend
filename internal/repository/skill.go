package repository

import (
	"context"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/lib/pq"
)

func (r *Repository) CreateSkill(ctx context.Context, skill *entities.Skill) error {
	const req = `
	INSERT INTO skills (teacher_id, category_id, video_card_link, about) 
	VALUES ($1, $2, $3, $4)
	`

	if _, err := r.db.ExecContext(ctx, req, skill.TeacherId, skill.CategoryId, skill.VideoCardLink, skill.About); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// Код ошибки 23505 означает unique_violation
			if pqErr.Code == "23505" {
				return errors.ErrorSkillRegistered
			}
		}

		return fmt.Errorf("failed to insert skill: %w", err)
	}
	return nil
}
