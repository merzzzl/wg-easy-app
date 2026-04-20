package auth

import (
	"context"
	"log"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) Authenticate(ctx context.Context, initData string) (model.User, bool, error) {
	log.Printf("info auth.authenticate called init_data_present=%t", initData != "")

	telegramUser, err := s.ValidateInitData(initData)
	if err != nil {
		log.Printf("info auth.authenticate validation_failed err=%v", err)

		return model.User{}, false, err
	}

	user, created, err := s.UpsertTelegramUser(ctx, telegramUser)
	if err != nil {
		log.Printf("info auth.authenticate upsert_failed telegram_id=%d err=%v", telegramUser.TelegramID, err)

		return model.User{}, false, err
	}

	log.Printf("info auth.authenticate succeeded telegram_id=%d created=%t", user.TelegramID, created)

	return user, created, nil
}
