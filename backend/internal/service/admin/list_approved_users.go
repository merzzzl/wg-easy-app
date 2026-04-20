package admin

import (
	"context"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) ListApprovedUsers(ctx context.Context) ([]model.User, error) {
	return s.db.ListUsersByStatus(ctx, model.UserStatusApproved)
}
