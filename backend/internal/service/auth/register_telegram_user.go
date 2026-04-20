package auth

import (
	"context"
	"log/slog"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) RegisterTelegramUser(ctx context.Context, telegramUser model.TelegramUser) (model.User, bool, error) {
	slog.Info("auth.register_telegram_user called", "telegram_id", telegramUser.TelegramID, "username", telegramUser.Username, "chat_id", telegramUser.ChatID)

	return s.UpsertTelegramUser(ctx, telegramUser)
}
