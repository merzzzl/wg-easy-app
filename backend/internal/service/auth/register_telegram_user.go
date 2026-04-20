package auth

import (
	"context"
	"log"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) RegisterTelegramUser(ctx context.Context, telegramUser model.TelegramUser) (model.User, bool, error) {
	log.Printf("info auth.register_telegram_user called telegram_id=%d username=%s chat_id=%d", telegramUser.TelegramID, telegramUser.Username, telegramUser.ChatID)

	return s.UpsertTelegramUser(ctx, telegramUser)
}
