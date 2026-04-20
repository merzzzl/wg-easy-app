package tunnel

import (
	"context"
	"log"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) ListByUserID(ctx context.Context, userID int64) ([]model.Tunnel, error) {
	log.Printf("info tunnel.list_by_user_id called user_id=%d", userID)

	return s.db.ListTunnelsByUserID(ctx, userID)
}
