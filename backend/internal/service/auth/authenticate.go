package auth

import (
	"context"
	"log/slog"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) Authenticate(ctx context.Context, initData string) (model.User, bool, error) {
	slog.Info("auth.authenticate called", "init_data_present", initData != "")

	telegramUser, err := s.ValidateInitData(initData)
	if err != nil {
		slog.Error("auth.authenticate validation failed", "error", err)

		return model.User{}, false, err
	}

	user, created, err := s.UpsertTelegramUser(ctx, telegramUser)
	if err != nil {
		slog.Error("auth.authenticate upsert failed", "telegram_id", telegramUser.TelegramID, "error", err)

		return model.User{}, false, err
	}

	slog.Info("auth.authenticate succeeded", "telegram_id", user.TelegramID, "created", created)

	return user, created, nil
}
