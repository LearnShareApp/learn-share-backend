package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	internalErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/Masterminds/squirrel"
)

func (r *Repository) CreateSkill(ctx context.Context, skill *entities.Skill) error {
	query, args, err := r.sqlBuilder.
		Insert("skills").
		Columns("teacher_id", "category_id", "video_card_link", "about").
		Values(skill.TeacherID, skill.CategoryID, skill.VideoCardLink, skill.About).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// Код ошибки 23505 означает unique_violation
			if pqErr.Code == "23505" {
				return internalErrs.ErrorNonUniqueData
			}
		}

		return fmt.Errorf("failed to insert skill: %w", err)
	}

	return nil
}

func (r *Repository) IsSkillExistsByTeacherIDAndCategoryID(ctx context.Context, teacherId int, categoryId int) (bool, error) {
	query, args, err := r.sqlBuilder.
		Select("1").
		From("skills").
		Where(squirrel.Eq{
			"teacher_id":  teacherId,
			"category_id": categoryId,
		}).
		Prefix("SELECT EXISTS(").
		Suffix(")").
		ToSql()

	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	var exists bool
	err = r.db.GetContext(ctx, &exists, query, args...)
	if err != nil {
		return false, fmt.Errorf("failed to check skill existence by teacher id and category id: %w", err)
	}

	return exists, nil
}

func (r *Repository) GetSkillIdByTeacherIdAndCategoryId(ctx context.Context, teacherId int, categoryId int) (int, error) {
	query, args, err := r.sqlBuilder.
		Select("skill_id").
		From("skills").
		Where(squirrel.Eq{
			"teacher_id":  teacherId,
			"category_id": categoryId,
		}).
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	var id int
	err = r.db.GetContext(ctx, &id, query, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, internalErrs.ErrorSelectEmpty
	}

	if err != nil {
		return 0, fmt.Errorf("failed to find skills: %w", err)
	}

	return id, nil
}

func (r *Repository) GetSkillByID(ctx context.Context, id int) (*entities.Skill, error) {
	query, args, err := r.sqlBuilder.
		Select(
			"skill_id",
			"teacher_id",
			"category_id",
			"video_card_link",
			"about",
			"rate",
			"total_rate_score",
			"reviews_count",
			"is_active",
		).
		From("skills").
		Where(squirrel.Eq{"skill_id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var skill entities.Skill
	err = r.db.GetContext(ctx, &skill, query, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, internalErrs.ErrorSelectEmpty
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find skill: %w", err)
	}

	return &skill, nil
}

func (r *Repository) GetSkillsByTeacherID(ctx context.Context, teacherID int) ([]*entities.Skill, error) {
	query, args, err := r.sqlBuilder.
		Select(
			"s.skill_id",
			"s.teacher_id",
			"s.category_id",
			"s.video_card_link",
			"s.about",
			"s.rate",
			"s.total_rate_score",
			"s.reviews_count",
			"s.is_active",
			"c.name as category_name",
		).
		From("skills s").
		InnerJoin("categories c ON s.category_id = c.category_id").
		Where(squirrel.Eq{"s.teacher_id": teacherID}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var skills []*entities.Skill
	err = r.db.SelectContext(ctx, &skills, query, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, internalErrs.ErrorSelectEmpty
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find skills: %w", err)
	}

	return skills, nil
}

func (r *Repository) ActivateSkillByID(ctx context.Context, id int) error {
	query, args, err := r.sqlBuilder.
		Update("skills").
		Set("is_active", true).
		Where(squirrel.Eq{"skill_id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update is_active field for skill: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return internalErrs.ErrorSelectEmpty
	}

	return nil
}
