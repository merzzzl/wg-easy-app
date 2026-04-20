package notification

import (
	"context"
	"log"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) NotifyTunnelDeleted(ctx context.Context, user *model.User, tunnel model.Tunnel) error {
	log.Printf("info notification.notify_tunnel_deleted called telegram_id=%d tunnel_id=%d admin=%s", user.TelegramID, tunnel.ID, s.adminUsername)

	return s.tg.SendMessage(ctx, s.adminUsername, actionText("Удален туннель", user.Username, user.TelegramID, tunnel.ID, tunnel.WGClientName))
}
