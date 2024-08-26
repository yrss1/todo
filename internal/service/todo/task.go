package todo

import (
	"context"
	"errors"
	"fmt"
	"github.com/yrss1/todo/internal/domain/task"
	"github.com/yrss1/todo/pkg/log"
	"github.com/yrss1/todo/pkg/store"
	"go.uber.org/zap"
)

func (s *Service) ListTasks(ctx context.Context) (res []task.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("ListTasks")

	data, err := s.taskRepository.List(ctx)
	if err != nil {
		logger.Error("failed to select", zap.Error(err))
		return
	}

	res = task.ParseFromEntities(data)

	return
}

func (s *Service) CreateTask(ctx context.Context, req task.Request) (res task.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("CreateTask")

	data := task.Entity{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		DueDate:     req.DueDate,
	}

	data.ID, err = s.taskRepository.Add(ctx, data)
	if err != nil {
		logger.Error("failed to create", zap.Error(err))
		return
	}

	res = task.ParseFromEntity(data)

	return
}

func (s *Service) GetTask(ctx context.Context, id string) (res task.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("GetTask").With(zap.String("id", id))

	data, err := s.taskRepository.Get(ctx, id)
	if err != nil {
		logger.Error("failed to get by id", zap.Error(err))
		return
	}

	res = task.ParseFromEntity(data)

	return
}

func (s *Service) UpdateTask(ctx context.Context, id string, req task.Request) (err error) {
	logger := log.LoggerFromContext(ctx).Named("UpdateTask").With(zap.String("id", id))

	data := task.Entity{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		DueDate:     req.DueDate,
	}

	err = s.taskRepository.Update(ctx, id, data)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		logger.Error("failed to update by id", zap.Error(err))
		return
	}

	return
}

func (s *Service) DeleteTask(ctx context.Context, id string) (err error) {
	logger := log.LoggerFromContext(ctx).Named("DeleteTask").With(zap.String("id", id))

	err = s.taskRepository.Delete(ctx, id)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		logger.Error("failed to delete by id", zap.Error(err))
		return
	}

	return
}

func (s *Service) SearchTask(ctx context.Context, req task.Request) (res []task.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("SearchTask")

	if req.Title != nil {
		logger = logger.With(zap.String("title", *(req.Title)))
	}
	if req.Status != nil {
		logger = logger.With(zap.String("status", fmt.Sprintf("%v", *(req.Status))))
	}

	searchData := task.Entity{
		Title:  req.Title,
		Status: req.Status,
	}
	data, err := s.taskRepository.Search(ctx, searchData)
	if err != nil {
		logger.Error("failed to search tasks", zap.Error(err))
		return
	}

	res = task.ParseFromEntities(data)

	return
}
