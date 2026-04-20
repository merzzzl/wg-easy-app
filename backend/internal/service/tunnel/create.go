package tunnel

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) Create(ctx context.Context, user *model.User) (model.Tunnel, error) {
	slog.Info("tunnel.create called", "user_id", user.ID, "username", user.Username)

	if err := ensureUserApproved(user); err != nil {
		slog.Warn("tunnel.create rejected for unapproved user", "user_id", user.ID, "status", user.Status)

		return model.Tunnel{}, err
	}

	tx, err := s.db.OpenTx(ctx)
	if err != nil {
		return model.Tunnel{}, fmt.Errorf("open transaction: %w", err)
	}

	committed := false

	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	tunnels, err := tx.ListTunnelsByUserID(ctx, user.ID)
	if err != nil {
		return model.Tunnel{}, fmt.Errorf("list tunnels: %w", err)
	}

	if len(tunnels) >= s.maxTunnels {
		return model.Tunnel{}, ErrTunnelLimitExceeded
	}

	tunnel, err := tx.CreateTunnel(ctx, model.CreateTunnelParams{UserID: user.ID})
	if err != nil {
		return model.Tunnel{}, fmt.Errorf("create tunnel record: %w", err)
	}

	wgClientName := buildWGClientName(user.Username, tunnel.ID)

	tunnel, err = tx.SetTunnelWGClientName(ctx, tunnel.ID, wgClientName)
	if err != nil {
		return model.Tunnel{}, fmt.Errorf("set wg client name: %w", err)
	}

	response, err := s.wg.CreateClient(ctx, model.WGEasyCreateClientParams{Name: wgClientName, ExpiresAt: nil})
	if err != nil {
		slog.Error("tunnel.create wg-easy create failed", "user_id", user.ID, "wg_client_name", wgClientName, "error", err)

		return model.Tunnel{}, fmt.Errorf("create wg-easy client: %w", err)
	}

	wgClientID := strconv.FormatInt(response.ClientID, 10)

	tunnel, err = tx.SetTunnelWGClientID(ctx, tunnel.ID, wgClientID)
	if err != nil {
		_ = s.wg.DeleteClient(ctx, wgClientID)

		return model.Tunnel{}, fmt.Errorf("set wg client id: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		_ = s.wg.DeleteClient(ctx, wgClientID)

		return model.Tunnel{}, fmt.Errorf("commit transaction: %w", err)
	}

	committed = true

	slog.Info("tunnel.create succeeded", "user_id", user.ID, "tunnel_id", tunnel.ID, "wg_client_name", tunnel.WGClientName, "wg_client_id", tunnel.WGClientID)

	return tunnel, nil
}
