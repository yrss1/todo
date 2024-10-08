package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/yrss1/todo/internal/domain/user"
	"github.com/yrss1/todo/pkg/store"
	"strings"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) List(ctx context.Context) (dest []user.Entity, err error) {
	query := `
		SELECT id, name, email 
		FROM users
		ORDER BY id`

	err = r.db.SelectContext(ctx, &dest, query)

	return
}

func (r *UserRepository) Add(ctx context.Context, data user.Entity) (id string, err error) {
	query := `
		INSERT INTO users (name, email, password) 
		VALUES ($1, $2, $3) 
		RETURNING id`

	args := []any{data.Name, data.Email, data.Password}

	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *UserRepository) Get(ctx context.Context, id string) (dest user.Entity, err error) {
	query := `
		SELECT id, name, email 
		FROM users
		WHERE id=$1`

	args := []any{id}

	if err = r.db.GetContext(ctx, &dest, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *UserRepository) Update(ctx context.Context, id string, data user.Entity) (err error) {
	sets, args := r.prepareArgs(data)

	if len(args) > 0 {
		args = append(args, id)
		sets = append(sets, "updated_at=CURRENT_TIMESTAMP")

		query := fmt.Sprintf("UPDATE users SET %s WHERE id=$%d RETURNING id", strings.Join(sets, ", "), len(args))

		if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = store.ErrorNotFound
			}
		}
	}

	return
}

func (r *UserRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
		DELETE FROM users
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

func (r *UserRepository) Search(ctx context.Context, data user.Entity) (dest []user.Entity, err error) {
	query := "SELECT id, name, email FROM users WHERE 1=1"

	sets, args := r.prepareArgs(data)
	if len(sets) > 0 {
		query += " AND " + strings.Join(sets, " AND ")
	}

	err = r.db.SelectContext(ctx, &dest, query, args...)

	return
}

func (r *UserRepository) prepareArgs(data user.Entity) (sets []string, args []any) {
	if data.Name != nil {
		args = append(args, data.Name)
		sets = append(sets, fmt.Sprintf("name=$%d", len(args)))
	}

	if data.Email != nil {
		args = append(args, data.Email)
		sets = append(sets, fmt.Sprintf("email=$%d", len(args)))
	}

	if data.Password != nil {
		args = append(args, data.Password)
		sets = append(sets, fmt.Sprintf("password=$%d", len(args)))
	}

	return
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (dest user.Entity, err error) {
	query := `SELECT id, name, email, password from users where email=$1`

	args := []any{email}

	if err = r.db.GetContext(ctx, &dest, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}
