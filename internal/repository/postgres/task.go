package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/yrss1/todo/internal/domain/task"
	"github.com/yrss1/todo/pkg/store"
	"strings"
)

type TaskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) List(ctx context.Context, userID string) (dest []task.Entity, err error) {
	query := `
        SELECT id, title, description, status 
        FROM tasks
        WHERE user_id = $1 
        ORDER BY id`

	err = r.db.SelectContext(ctx, &dest, query, userID)
	return
}

func (r *TaskRepository) Add(ctx context.Context, data task.Entity) (id string, err error) {
	query := `
		INSERT INTO tasks (user_id, title, description, status) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id`

	args := []any{data.UserID, data.Title, data.Description, data.Status}

	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *TaskRepository) Get(ctx context.Context, userID string, taskID string) (dest task.Entity, err error) {
	query := `
	   SELECT id, title, description, status
	   FROM tasks
	   WHERE id = $1 AND user_id = $2`

	args := []any{taskID, userID}

	if err = r.db.GetContext(ctx, &dest, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}
	return
}

func (r *TaskRepository) Update(ctx context.Context, userID string, taskID string, data task.Entity) (err error) {
	sets, args := r.prepareArgs(data)

	if len(args) > 0 {
		args = append(args, taskID, userID)
		sets = append(sets, "updated_at=CURRENT_TIMESTAMP")

		query := fmt.Sprintf(
			"UPDATE tasks SET %s WHERE id=$%d AND user_id=$%d RETURNING id",
			strings.Join(sets, ", "),
			len(args)-1, // Позиция taskID
			len(args),   // Позиция userID
		)

		if err = r.db.QueryRowContext(ctx, query, args...).Scan(&taskID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = store.ErrorNotFound
			}
		}
	}

	return
}

func (r *TaskRepository) Delete(ctx context.Context, userID string, taskID string) (err error) {
	query := `
        DELETE FROM tasks
        WHERE id = $1 AND user_id = $2
        RETURNING id`

	args := []any{taskID, userID}

	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&taskID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *TaskRepository) prepareArgs(data task.Entity) (sets []string, args []any) {
	if data.Title != nil {
		args = append(args, data.Title)
		sets = append(sets, fmt.Sprintf("title=$%d", len(args)))
	}

	if data.Description != nil {
		args = append(args, data.Description)
		sets = append(sets, fmt.Sprintf("description=$%d", len(args)))
	}

	if data.Status != nil {
		args = append(args, data.Status)
		sets = append(sets, fmt.Sprintf("status=$%d", len(args)))
	}

	return
}
