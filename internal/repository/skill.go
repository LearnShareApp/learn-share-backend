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

func (r *Repository) CreateSkill(ctx context.Context, skill *entities.Skill) error {
	const query = `
	INSERT INTO skills (teacher_id, category_id, video_card_link, about) 
	VALUES ($1, $2, $3, $4)
	`

	if _, err := r.db.ExecContext(ctx, query, skill.TeacherId, skill.CategoryId, skill.VideoCardLink, skill.About); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// Код ошибки 23505 означает unique_violation
			if pqErr.Code == "23505" {
				return internalErrs.ErrorSkillRegistered
			}
		}

		return fmt.Errorf("failed to insert skill: %w", err)
	}
	return nil
}

func (r *Repository) GetSkillsByTeacherId(ctx context.Context, id int) ([]*entities.Skill, error) {
	const query = `
	SELECT 
		s.skill_id, 
		s.teacher_id, 
		s.category_id, 
		s.video_card_link, 
		s.about, 
		s.rate, 
		s.is_active,
		c.name as category_name
	FROM public.skills s
	INNER JOIN public.categories c ON s.category_id = c.category_id
	WHERE teacher_id = $1`

	var skills []*entities.Skill
	err := r.db.SelectContext(ctx, &skills, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, internalErrs.ErrorSelectEmpty
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find skills: %w", err)
	}

	return skills, nil
}
