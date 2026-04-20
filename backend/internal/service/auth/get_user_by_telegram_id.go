package auth

import (
	"context"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) GetUserByTelegramID(ctx context.Context, telegramID int64) (model.User, error) {
	return s.db.GetUserByTelegramID(ctx, telegramID)
}
