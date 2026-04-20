package auth

import (
	"context"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) SetUserStatusByTelegramID(ctx context.Context, telegramID int64, status model.UserStatus) (model.User, error) {
	return s.db.SetUserStatusByTelegramID(ctx, telegramID, status)
}
