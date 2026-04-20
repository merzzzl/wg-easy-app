package notification

import (
	"context"
	"log/slog"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) NotifyTunnelCreated(ctx context.Context, user *model.User, tunnel model.Tunnel) error {
	slog.Info("notification.notify_tunnel_created called", "telegram_id", user.TelegramID, "tunnel_id", tunnel.ID, "admin_username", s.adminUsername)

	return s.sendAdminMessage(ctx, actionText("Создан туннель", user.Username, user.TelegramID, tunnel.ID, tunnel.WGClientName))
}
