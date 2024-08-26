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

func (r *TaskRepository) List(ctx context.Context) (dest []task.Entity, err error) {
	query := `
		SELECT id, user_id, title, description, status, due_date 
		FROM tasks
		ORDER BY id`

	err = r.db.SelectContext(ctx, &dest, query)

	return
}

func (r *TaskRepository) Add(ctx context.Context, data task.Entity) (id string, err error) {
	query := `
		INSERT INTO tasks (user_id, title, description, status, due_date) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id`

	args := []any{data.UserID, data.Title, data.Description, data.Status, data.DueDate}

	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *TaskRepository) Get(ctx context.Context, id string) (dest task.Entity, err error) {
	query := `
		SELECT id, user_id, title, description, status, due_date
		FROM tasks
		WHERE id=$1`

	args := []any{id}

	if err = r.db.GetContext(ctx, &dest, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *TaskRepository) Update(ctx context.Context, id string, data task.Entity) (err error) {
	sets, args := r.prepareArgs(data)

	if len(args) > 0 {
		args = append(args, id)
		sets = append(sets, "updated_at=CURRENT_TIMESTAMP")

		query := fmt.Sprintf("UPDATE tasks SET %s WHERE id=$%d RETURNING id", strings.Join(sets, ", "), len(args))

		if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = store.ErrorNotFound
			}
		}
	}

	return
}

func (r *TaskRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
		DELETE FROM tasks
		WHERE id=$1
		RETURNING id`

	args := []any{id}

	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *TaskRepository) Search(ctx context.Context, data task.Entity) (dest []task.Entity, err error) {
	query := "SELECT id, user_id, title, description, status, due_date FROM tasks WHERE 1=1"
	sets, args := r.prepareArgs(data)
	if len(sets) > 0 {
		query += " AND " + strings.Join(sets, " AND ")
	}

	err = r.db.SelectContext(ctx, &dest, query, args...)

	return
}

func (r *TaskRepository) prepareArgs(data task.Entity) (sets []string, args []any) {
	if data.UserID != nil {
		args = append(args, data.UserID)
		sets = append(sets, fmt.Sprintf("user_id=$%d", len(args)))
	}

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

	if data.DueDate != nil {
		args = append(args, data.DueDate)
		sets = append(sets, fmt.Sprintf("due_date=$%d", len(args)))
	}

	return
}
