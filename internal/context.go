package internal

import (
	"context"
	"photos"
)

type contextKey uint

const (
	userCtxKey contextKey = iota
)

func ContextWithUser(ctx context.Context, user *photos.User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

func UserFromContext(ctx context.Context) *photos.User {

	userInf := ctx.Value(userCtxKey)
	if userInf == nil {
		return nil
	}

	if user, ok := userInf.(*photos.User); ok {
		return user
	}

	return nil

}
