package tunnel

import (
	"context"
	"log/slog"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) ListByUserID(ctx context.Context, userID int64) ([]model.Tunnel, error) {
	slog.Info("tunnel.list_by_user_id called", "user_id", userID)

	return s.db.ListTunnelsByUserID(ctx, userID)
}
