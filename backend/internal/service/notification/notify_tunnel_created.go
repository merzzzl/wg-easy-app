package notification

import (
	"context"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) NotifyTunnelCreated(ctx context.Context, user *model.User, tunnel model.Tunnel) error {
	return s.tg.SendMessage(ctx, s.adminUsername, actionText("Создан туннель", user.Username, user.TelegramID, tunnel.ID, tunnel.WGClientName))
}
