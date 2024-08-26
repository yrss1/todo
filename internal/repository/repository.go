package repository

import (
	"github.com/yrss1/todo/internal/domain/task"
	"github.com/yrss1/todo/internal/domain/user"
	"github.com/yrss1/todo/internal/repository/postgres"
	"github.com/yrss1/todo/pkg/store"
)

type Configuration func(r *Repository) error

type Repository struct {
	postgres store.SQLX

	User user.Repository
	Task task.Repository
}

func New(configs ...Configuration) (s *Repository, err error) {
	s = &Repository{}

	for _, cfg := range configs {
		if err = cfg(s); err != nil {
			return
		}
	}

	return
}

func WithPostgresStore(dbName string) Configuration {
	return func(r *Repository) (err error) {
		r.postgres, err = store.New(dbName)
		if err != nil {
			return
		}

		if err = store.Migrate(dbName); err != nil {
			return
		}

		r.User = postgres.NewUserRepository(r.postgres.Client)
		r.Task = postgres.NewTaskRepository(r.postgres.Client)

		return
	}
}
