package todo

import (
	"context"
	"errors"
	"github.com/yrss1/todo/internal/domain/task"
	"github.com/yrss1/todo/pkg/log"
	"github.com/yrss1/todo/pkg/store"
	"go.uber.org/zap"
)

func (s *Service) ListTasks(ctx context.Context, userID string) (res []task.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("ListTasks")

	data, err := s.taskRepository.List(ctx, userID)
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
	}

	data.ID, err = s.taskRepository.Add(ctx, data)
	if err != nil {
		logger.Error("failed to create", zap.Error(err))
		return
	}

	res = task.ParseFromEntity(data)

	return
}

func (s *Service) GetTask(ctx context.Context, userID string, taskID string) (res task.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("GetTask").
		With(zap.String("userID", userID), zap.String("taskID", taskID))

	data, err := s.taskRepository.Get(ctx, userID, taskID)
	if err != nil {
		logger.Error("failed to get by id", zap.Error(err))
		return
	}

	res = task.ParseFromEntity(data)

	return
}

func (s *Service) UpdateTask(ctx context.Context, userID string, taskID string, req task.Request) (err error) {
	logger := log.LoggerFromContext(ctx).Named("UpdateTask").
		With(zap.String("userID", userID), zap.String("taskID", taskID))

	data := task.Entity{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}

	err = s.taskRepository.Update(ctx, userID, taskID, data)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		logger.Error("failed to update by id", zap.Error(err))
		return
	}

	return
}

func (s *Service) DeleteTask(ctx context.Context, userID string, taskID string) (err error) {
	logger := log.LoggerFromContext(ctx).Named("DeleteTask").
		With(zap.String("userID", userID), zap.String("taskID", taskID))

	err = s.taskRepository.Delete(ctx, userID, taskID)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		logger.Error("failed to delete by id", zap.Error(err))
		return
	}

	return
}
