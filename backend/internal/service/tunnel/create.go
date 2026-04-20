package tunnel

import (
	"context"
	"fmt"
	"log"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) Create(ctx context.Context, user *model.User) (model.Tunnel, error) {
	log.Printf("info tunnel.create called user_id=%d username=%s", user.ID, user.Username)

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

	response, err := s.wg.CreateClient(ctx, model.WGEasyCreateClientParams{Name: wgClientName})
	if err != nil {
		return model.Tunnel{}, fmt.Errorf("create wg-easy client: %w", err)
	}

	tunnel, err = tx.SetTunnelWGClientID(ctx, tunnel.ID, response.ClientID)
	if err != nil {
		_ = s.wg.DeleteClient(ctx, response.ClientID)

		return model.Tunnel{}, fmt.Errorf("set wg client id: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		_ = s.wg.DeleteClient(ctx, response.ClientID)

		return model.Tunnel{}, fmt.Errorf("commit transaction: %w", err)
	}

	committed = true

	log.Printf("info tunnel.create succeeded user_id=%d tunnel_id=%d wg_client_name=%s wg_client_id=%s", user.ID, tunnel.ID, tunnel.WGClientName, tunnel.WGClientID)

	return tunnel, nil
}
