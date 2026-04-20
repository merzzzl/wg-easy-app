package admin

import (
	"context"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) ApproveUser(ctx context.Context, telegramID int64) (model.User, error) {
	return s.db.SetUserStatusByTelegramID(ctx, telegramID, model.UserStatusApproved)
}
