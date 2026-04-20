package tunnel

import (
	"context"
	"fmt"
	"log"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) SendConfig(ctx context.Context, user *model.User, tunnelID int64) error {
	log.Printf("info tunnel.send_config called user_id=%d tunnel_id=%d chat_id=%d", user.ID, tunnelID, user.ChatID)

	tunnels, err := s.db.ListTunnelsByUserID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("list tunnels: %w", err)
	}

	tunnel, err := findTunnelByID(tunnels, tunnelID)
	if err != nil {
		return err
	}

	configText, err := s.wg.GetClientConfiguration(ctx, tunnel.WGClientID)
	if err != nil {
		return fmt.Errorf("get client configuration: %w", err)
	}

	fileName := tunnel.WGClientName + ".conf"
	if err := s.tg.SendDocument(ctx, user.ChatID, fileName, tunnel.WGClientName, []byte(configText)); err != nil {
		return fmt.Errorf("send config document: %w", err)
	}

	log.Printf("info tunnel.send_config succeeded user_id=%d tunnel_id=%d wg_client_name=%s", user.ID, tunnel.ID, tunnel.WGClientName)

	return nil
}
