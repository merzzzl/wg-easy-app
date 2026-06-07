package notification

import (
	"context"
	"fmt"
	"log/slog"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) NotifyApprovedUsers(ctx context.Context, text string) (int, error) {
	users, err := s.db.ListUsersByStatus(ctx, model.UserStatusApproved)
	if err != nil {
		return 0, fmt.Errorf("list approved users: %w", err)
	}

	slog.Info("notification.notify_approved_users called", "approved_users", len(users))

	for i, user := range users {
		if err := s.tg.SendMessage(ctx, user.ChatID, text); err != nil {
			return i, fmt.Errorf("send notification to @%s: %w", user.Username, err)
		}
	}

	return len(users), nil
}
