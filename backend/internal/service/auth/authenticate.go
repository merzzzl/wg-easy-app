package auth

import (
	"context"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) Authenticate(ctx context.Context, initData string) (model.User, bool, error) {
	telegramUser, err := s.ValidateInitData(initData)
	if err != nil {
		return model.User{}, false, err
	}

	user, created, err := s.UpsertTelegramUser(ctx, telegramUser)
	if err != nil {
		return model.User{}, false, err
	}

	return user, created, nil
}
