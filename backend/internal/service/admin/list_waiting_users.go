package admin

import (
	"context"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) ListWaitingUsers(ctx context.Context) ([]model.User, error) {
	return s.db.ListUsersByStatuses(ctx, model.UserStatusPending, model.UserStatusWaitingApprove)
}
