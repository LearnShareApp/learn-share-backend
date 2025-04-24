package repository

import (
	"context"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

func (r *Repository) CreateComplaint(ctx context.Context, complaint *entities.Complaint) error {
	query, args, err := r.sqlBuilder.
		Insert("complaints").
		Columns("complainer_id", "reported_id", "reason", "description").
		Values(complaint.ComplainerID, complaint.ReportedID, complaint.Reason, complaint.Description).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert complaint: %w", err)
	}

	return nil
}

func (r *Repository) GetAllComplaints(ctx context.Context) ([]*entities.Complaint, error) {
	query, args, err := r.sqlBuilder.
		Select(
			"complaint_id",
			"complainer_id",
			"reported_id",
			"reason",
			"description",
			"created_at",
		).
		From("complaints").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var complaints []*entities.Complaint

	if err := r.db.SelectContext(ctx, &complaints, query, args...); err != nil {
		return nil, fmt.Errorf("failed to select complaints: %w", err)
	}

	return complaints, nil
}
