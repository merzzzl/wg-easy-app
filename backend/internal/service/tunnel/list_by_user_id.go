package tunnel

import (
	"context"
	"log/slog"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) ListByUserID(ctx context.Context, user *model.User) ([]model.Tunnel, error) {
	slog.Info("tunnel.list_by_user_id called", "user_id", user.ID)

	if err := ensureUserApproved(user); err != nil {
		slog.Warn("tunnel.list_by_user_id rejected for unapproved user", "user_id", user.ID, "status", user.Status)

		return nil, err
	}

	return s.db.ListTunnelsByUserID(ctx, user.ID)
}
