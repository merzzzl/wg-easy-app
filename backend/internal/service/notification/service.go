package notification

import (
	"fmt"

	"wg-easy-app/backend/internal/config"
	"wg-easy-app/backend/internal/repository/telegram"
)

type Service struct {
	adminUsername string
	tg            *telegram.Repository
}

func New(cfg *config.Config, tg *telegram.Repository) *Service {
	return &Service{
		adminUsername: cfg.AdminUsername,
		tg:            tg,
	}
}

func formatUsername(username string) string {
	if username == "" {
		return "-"
	}

	return "@" + username
}

func formatTunnelLine(wgClientName string) string {
	if wgClientName == "" {
		return "-"
	}

	return wgClientName
}

func actionText(action, username string, telegramID, tunnelID int64, wgClientName string) string {
	message := fmt.Sprintf("%s\nusername: %s\ntelegram_id: %d", action, formatUsername(username), telegramID)
	if tunnelID == 0 {
		return message
	}

	return fmt.Sprintf("%s\ntunnel_id: %d\nwg_client_name: %s", message, tunnelID, formatTunnelLine(wgClientName))
}

func registrationText(username string, telegramID int64, status string) string {
	return fmt.Sprintf("Новая регистрация, требуется подтверждение\nusername: %s\ntelegram_id: %d\nstatus: %s", formatUsername(username), telegramID, status)
}
