package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	internalErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

func (r *Repository) GetStateMachineItemByID(ctx context.Context, id int) (*entities.StateMachineItem, error) {
	query, args, err := r.sqlBuilder.
		Select(
			"i.item_id",
			"i.state_machine_id",
			"i.state_id",
			"s.name as state_name",
		).
		From("state_machines_items i").
		InnerJoin("states s ON i.state_id = s.state_id").
		Where(squirrel.Eq{"item_id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var stItem entities.StateMachineItem
	err = r.db.GetContext(ctx, &stItem, query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrs.ErrorSelectEmpty
		}

		return nil, fmt.Errorf("failed to find state machine item: %w", err)
	}

	return &stItem, nil
}

func (r *Repository) UpdateStateMachineItemState(ctx context.Context, stateMachineItemID, newStateID int) error {
	query, args, err := r.sqlBuilder.
		Update("state_machines_items").
		Set("state_id", newStateID).
		Where(squirrel.Eq{"item_id": stateMachineItemID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update state machine item: %w", err)
	}

	return nil
}

func (r *Repository) GetStateByID(ctx context.Context, id int) (*entities.State, error) {
	query, args, err := r.sqlBuilder.
		Select(
			"state_id",
			"name",
		).
		From("states").
		Where(squirrel.Eq{
			"state_id": id,
		}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var state entities.State
	err = r.db.GetContext(ctx, &state, query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrs.ErrorSelectEmpty
		}

		return nil, fmt.Errorf("failed to get state by id: %w", err)
	}

	return &state, nil
}

func (r *Repository) GetStateIDByName(ctx context.Context, name entities.StateName) (int, error) {
	query, args, err := r.sqlBuilder.
		Select(
			"state_id",
		).
		From("states").
		Where(squirrel.Eq{
			"name": name,
		}).
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	var id int
	err = r.db.GetContext(ctx, &id, query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, internalErrs.ErrorSelectEmpty
		}

		return 0, fmt.Errorf("failed to get state by name: %w", err)
	}

	return id, nil
}

func (r *Repository) CheckIsTransitionAvailable(ctx context.Context, stateMachineID, currentStateID, nextStateID int) (bool, error) {
	query, args, err := r.sqlBuilder.
		Select(
			"transition_id",
		).
		From("state_transitions").
		Where(squirrel.Eq{
			"state_machine_id": stateMachineID,
			"current_state_id": currentStateID,
			"next_state_id":    nextStateID,
		}).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	var id int
	err = r.db.GetContext(ctx, &id, query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("failed check transition available: %w", err)
	}

	return true, nil
}

func (r *Repository) getStateMachineByName(ctx context.Context, name entities.StateMachineName) (*entities.StateMachine, error) {
	query, args, err := r.sqlBuilder.
		Select(
			"state_machine_id",
			"name",
			"start_state_id",
		).
		From("state_machines").
		Where(squirrel.Eq{"name": name}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var stateMachine entities.StateMachine

	err = r.db.GetContext(ctx, &stateMachine, query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrs.ErrorSelectEmpty
		}

		return nil, fmt.Errorf("failed to extract state machine by name: %w", err)
	}

	return &stateMachine, nil
}

func (r *Repository) insertStateMachineItem(ctx context.Context, tx *sqlx.Tx, machine entities.StateMachine) (int, error) {
	query, args, err := r.sqlBuilder.
		Insert("state_machines_items").
		Columns("state_machine_id", "state_id").
		Values(machine.ID, machine.StartStateID).
		Suffix("RETURNING item_id").
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("failed to build insert query: %w", err)
	}

	var itemID int

	if err := tx.GetContext(ctx, &itemID, query, args...); err != nil {
		return 0, fmt.Errorf("failed to insert state machine item: %w", err)
	}

	return itemID, nil
}
