package notification

import (
	"context"
	"log"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) NotifyRegistration(ctx context.Context, user *model.User) error {
	log.Printf("info notification.notify_registration called telegram_id=%d username=%s admin=%s", user.TelegramID, user.Username, s.adminUsername)

	return s.tg.SendMessage(ctx, s.adminUsername, registrationText(user.Username, user.TelegramID, string(user.Status)))
}
