package auth

import (
	"context"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) UpsertTelegramUser(ctx context.Context, user model.TelegramUser) (model.User, bool, error) {
	return s.db.UpsertUser(ctx, model.UserUpsertParams(user))
}
