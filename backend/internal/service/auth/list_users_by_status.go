package auth

import (
	"context"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) ListUsersByStatus(ctx context.Context, status model.UserStatus) ([]model.User, error) {
	return s.db.ListUsersByStatus(ctx, status)
}
