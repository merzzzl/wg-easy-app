package notification

import (
	"context"
	"log/slog"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) NotifyRegistration(ctx context.Context, user *model.User) error {
	slog.Info("notification.notify_registration called", "telegram_id", user.TelegramID, "username", user.Username, "admin_username", s.adminUsername)

	return s.sendAdminMessage(ctx, registrationText(user.Username, string(user.Status))+"\n\nApprove: `/approve @"+escapeMarkdown(user.Username)+"`\nRevoke: `/revoke @"+escapeMarkdown(user.Username)+"`")
}
