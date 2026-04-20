package notification

import (
	"context"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) NotifyRegistration(ctx context.Context, user *model.User) error {
	return s.tg.SendMessage(ctx, s.adminUsername, registrationText(user.Username, user.TelegramID, string(user.Status)))
}
