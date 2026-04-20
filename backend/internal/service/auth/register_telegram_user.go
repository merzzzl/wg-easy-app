package auth

import (
	"context"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) RegisterTelegramUser(ctx context.Context, telegramUser model.TelegramUser) (model.User, bool, error) {
	return s.UpsertTelegramUser(ctx, telegramUser)
}
