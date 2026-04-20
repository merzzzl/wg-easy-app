package auth

import (
	"context"
	"log"
)

func (s *Service) SendStartMessage(ctx context.Context, chatID int64) error {
	log.Printf("info auth.send_start_message called chat_id=%d has_mini_app_url=%t", chatID, s.config.MiniAppURL != "")

	text := "Добро пожаловать. Здесь вы можете управлять своими WireGuard-конфигами."
	if s.config.MiniAppURL == "" {
		return s.tg.SendMessage(ctx, chatID, text)
	}

	return s.tg.SendWebAppMessage(ctx, chatID, text, "Открыть приложение", s.config.MiniAppURL)
}
