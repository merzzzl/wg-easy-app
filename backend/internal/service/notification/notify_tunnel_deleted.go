package notification

import (
	"context"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) NotifyTunnelDeleted(ctx context.Context, user *model.User, tunnel model.Tunnel) error {
	return s.tg.SendMessage(ctx, s.adminUsername, actionText("Удален туннель", user.Username, user.TelegramID, tunnel.ID, tunnel.WGClientName))
}
