package tunnel

import (
	"context"
	"fmt"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) Delete(ctx context.Context, user *model.User, tunnelID int64) (model.Tunnel, error) {
	tunnels, err := s.db.ListTunnelsByUserID(ctx, user.ID)
	if err != nil {
		return model.Tunnel{}, fmt.Errorf("list tunnels: %w", err)
	}

	tunnel, err := findTunnelByID(tunnels, tunnelID)
	if err != nil {
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

	if err := tx.DeleteTunnel(ctx, tunnel.ID); err != nil {
		return model.Tunnel{}, fmt.Errorf("delete tunnel record: %w", err)
	}

	if err := s.wg.DeleteClient(ctx, tunnel.WGClientID); err != nil {
		return model.Tunnel{}, fmt.Errorf("delete wg-easy client: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Tunnel{}, fmt.Errorf("commit transaction: %w", err)
	}

	committed = true

	return tunnel, nil
}
