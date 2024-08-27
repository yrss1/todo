package auth

import "context"

type Repository interface {
	SignUp(ctx context.Context, data Entity) (id string, err error)
	//SignIn(ctx context.Context, id string) (dest Entity, err error)
}
