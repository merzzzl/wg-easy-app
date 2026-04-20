package auth

import "context"

func (s *Service) SendStartMessage(ctx context.Context, chatID int64) error {
	text := "Добро пожаловать. Здесь вы можете управлять своими WireGuard-конфигами."
	if s.config.MiniAppURL == "" {
		return s.tg.SendMessage(ctx, chatID, text)
	}

	return s.tg.SendWebAppMessage(ctx, chatID, text, "Открыть приложение", s.config.MiniAppURL)
}
