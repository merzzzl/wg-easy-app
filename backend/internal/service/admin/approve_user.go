package admin

import (
	"context"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) ApproveUser(ctx context.Context, username string) (model.User, error) {
	return s.db.SetUserStatusByUsername(ctx, username, model.UserStatusApproved)
}
