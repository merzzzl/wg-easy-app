package tunnel

import (
	"context"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) ListByUserID(ctx context.Context, userID int64) ([]model.Tunnel, error) {
	return s.db.ListTunnelsByUserID(ctx, userID)
}
