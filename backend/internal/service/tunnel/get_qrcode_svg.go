package tunnel

import (
	"context"
	"fmt"
	"log/slog"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) GetQRCodeSVG(ctx context.Context, user *model.User, tunnelID int64) (string, error) {
	slog.Info("tunnel.get_qrcode_svg called", "user_id", user.ID, "tunnel_id", tunnelID)

	tunnels, err := s.db.ListTunnelsByUserID(ctx, user.ID)
	if err != nil {
		return "", fmt.Errorf("list tunnels: %w", err)
	}

	tunnel, err := findTunnelByID(tunnels, tunnelID)
	if err != nil {
		return "", err
	}

	qrcode, err := s.wg.GetClientQRCodeSVG(ctx, tunnel.WGClientID)
	if err != nil {
		return "", fmt.Errorf("get qrcode svg: %w", err)
	}

	return qrcode, nil
}
