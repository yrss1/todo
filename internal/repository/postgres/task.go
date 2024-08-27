package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/yrss1/todo/internal/domain/task"
	"github.com/yrss1/todo/pkg/store"
	"strconv"
	"strings"
)

type TaskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) List(ctx context.Context, userID, titleFilter, statusFilter, sortBy, sortOrder string, page, limit int) ([]task.Entity, error) {
	var (
		baseQuery strings.Builder
		args      []interface{}
		dest      []task.Entity
		paramIdx  = 1
		offset    = (page - 1) * limit
	)

	baseQuery.WriteString(`SELECT id, title, description, status FROM tasks WHERE user_id = $1`)
	args = append(args, userID)

	if titleFilter != "" {
		baseQuery.WriteString(fmt.Sprintf(` AND title ILIKE $%d`, paramIdx+1))
		args = append(args, "%"+titleFilter+"%")
		paramIdx++
	}

	if statusFilter != "" {
		baseQuery.WriteString(fmt.Sprintf(` AND status = $%d`, paramIdx+1))
		args = append(args, statusFilter)
		paramIdx++
	}

	if sortBy != "" {
		baseQuery.WriteString(fmt.Sprintf(` ORDER BY %s %s`, sortBy, sortOrder))
	} else {
		baseQuery.WriteString(` ORDER BY id`)
	}

	baseQuery.WriteString(fmt.Sprintf(` LIMIT $%d OFFSET $%d`, paramIdx+1, paramIdx+2))
	args = append(args, limit, offset)

	err := r.db.SelectContext(ctx, &dest, baseQuery.String(), args...)
	return dest, err
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

func (r *TaskRepository) buildQuery(userID, titleFilter, statusFilter, sortBy, sortOrder string) (string, []interface{}) {
	var queryBuilder strings.Builder
	queryBuilder.WriteString(`
        SELECT id, title, description, status
        FROM tasks
        WHERE user_id = $1`)

	args := []interface{}{userID}

	paramIdx := 2
	if titleFilter != "" {
		queryBuilder.WriteString(` AND title ILIKE $`)
		queryBuilder.WriteString(strconv.Itoa(paramIdx))
		args = append(args, "%"+titleFilter+"%")
		paramIdx++
	}

	if statusFilter != "" {
		queryBuilder.WriteString(` AND status = $`)
		queryBuilder.WriteString(strconv.Itoa(paramIdx))
		args = append(args, statusFilter)
		paramIdx++
	}

	if sortBy != "" {
		queryBuilder.WriteString(` ORDER BY `)
		queryBuilder.WriteString(sortBy)
		queryBuilder.WriteString(` `)
		queryBuilder.WriteString(sortOrder)
	} else {
		queryBuilder.WriteString(` ORDER BY id`)
	}

	return queryBuilder.String(), args
}
