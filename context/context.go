package context

import (
	"context"

	"soramon0/webapp/models"
)

const (
	userKey = privateKey("user")
)

type privateKey string

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	if tmp := ctx.Value(userKey); tmp != nil {
		if u, ok := tmp.(*models.User); ok {
			return u
		}
	}

	return nil
}
