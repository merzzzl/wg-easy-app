package tunnel

import (
	"context"
	"fmt"
	"log/slog"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) DeleteAllByUser(ctx context.Context, user *model.User) error {
	slog.Info("tunnel.delete_all_by_user called", "user_id", user.ID)

	tx, err := s.db.OpenTx(ctx)
	if err != nil {
		return fmt.Errorf("open transaction: %w", err)
	}

	committed := false

	defer func() {
		if !committed {
			if err := tx.Rollback(ctx); err != nil {
				slog.Warn("tunnel.delete_all_by_user rollback failed", "error", err)
			}
		}
	}()

	tunnels, err := tx.ListTunnelsByUserID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("list tunnels: %w", err)
	}

	for _, tunnel := range tunnels {
		if tunnel.WGClientID != "" {
			if err := s.wg.DeleteClient(ctx, tunnel.WGClientID); err != nil {
				slog.Error("tunnel.delete_all_by_user wg delete failed", "user_id", user.ID, "tunnel_id", tunnel.ID, "error", err)

				return fmt.Errorf("delete wg client: %w", err)
			}
		}

		if err := tx.DeleteTunnel(ctx, tunnel.ID); err != nil {
			return fmt.Errorf("delete tunnel record: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	committed = true

	slog.Info("tunnel.delete_all_by_user succeeded", "user_id", user.ID, "deleted_tunnels", len(tunnels))

	return nil
}
