package repository

import (
	"context"
	"fmt"
)

func (r *Repository) GetStatusIDByStatusName(ctx context.Context, name string) (int, error) {
	const query = `
	SELECT status_id FROM statuses WHERE name = $1
	`

	var statusId int
	if err := r.db.GetContext(ctx, &statusId, query, name); err != nil {
		return 0, fmt.Errorf("failed to get statusId by status name: %w", err)
	}

	return statusId, nil
}
