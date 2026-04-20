package controller

import (
	"context"

	"wg-easy-app/backend/internal/model"
)

type contextKey string

const currentUserKey contextKey = "current-user"

func withCurrentUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, currentUserKey, user)
}

func currentUser(ctx context.Context) (*model.User, bool) {
	user, ok := ctx.Value(currentUserKey).(*model.User)

	return user, ok
}
