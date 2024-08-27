package task

import "context"

type Repository interface {
	List(ctx context.Context, userID string) (dest []Entity, err error)
	Add(ctx context.Context, data Entity) (id string, err error)
	Get(ctx context.Context, userID string, taskID string) (dest Entity, err error)
	Update(ctx context.Context, userID string, taskID string, dest Entity) (err error)
	Delete(ctx context.Context, userID string, taskID string) (err error)
}
