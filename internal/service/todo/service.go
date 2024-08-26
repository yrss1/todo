package todo

import (
	"github.com/yrss1/todo/internal/domain/task"
)

type Configuration func(s *Service) error

type Service struct {
	taskRepository task.Repository
}

func New(configs ...Configuration) (s *Service, err error) {
	s = &Service{}

	for _, cfg := range configs {
		if err = cfg(s); err != nil {
			return
		}
	}

	return
}

func WithTaskRepository(taskRepository task.Repository) Configuration {
	return func(s *Service) error {
		s.taskRepository = taskRepository
		return nil
	}
}
