package admin

import (
	"context"
	"fmt"
	"log/slog"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) RevokeUser(ctx context.Context, telegramID int64) (model.User, int, error) {
	tx, err := s.db.OpenTx(ctx)
	if err != nil {
		return model.User{}, 0, fmt.Errorf("open transaction: %w", err)
	}

	committed := false

	defer func() {
		if !committed {
			if err := tx.Rollback(ctx); err != nil {
				slog.Warn("admin.revoke_user rollback failed", "error", err)
			}
		}
	}()

	user, err := tx.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		return model.User{}, 0, fmt.Errorf("get user by telegram id: %w", err)
	}

	tunnels, err := tx.ListTunnelsByUserID(ctx, user.ID)
	if err != nil {
		return model.User{}, 0, fmt.Errorf("list tunnels by user: %w", err)
	}

	for _, tunnel := range tunnels {
		if tunnel.WGClientID != "" {
			if err := s.wg.DeleteClient(ctx, tunnel.WGClientID); err != nil {
				slog.Error("admin.revoke_user wg delete failed", "telegram_id", telegramID, "tunnel_id", tunnel.ID, "error", err)

				return model.User{}, 0, fmt.Errorf("delete wg client: %w", err)
			}
		}

		if err := tx.DeleteTunnel(ctx, tunnel.ID); err != nil {
			return model.User{}, 0, fmt.Errorf("delete tunnel record: %w", err)
		}
	}

	user, err = tx.SetUserStatusByTelegramID(ctx, telegramID, model.UserStatusWaitingApprove)
	if err != nil {
		return model.User{}, 0, fmt.Errorf("set waiting_approve status: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.User{}, 0, fmt.Errorf("commit transaction: %w", err)
	}

	committed = true

	return user, len(tunnels), nil
}
