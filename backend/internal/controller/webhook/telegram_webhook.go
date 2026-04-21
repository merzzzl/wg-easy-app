package webhook

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-telegram/bot/models"

	"wg-easy-app/backend/internal/model"
	"wg-easy-app/backend/internal/service/auth"
)

const adminCommandArgs = 2

func normalizeAdminUsername(username string) string {
	return strings.TrimPrefix(strings.TrimSpace(username), "@")
}

func (c *Controller) HandleUpdate(ctx context.Context, update *models.Update) error {
	if update == nil {
		return nil
	}

	message := update.Message
	if message == nil || strings.TrimSpace(message.Text) == "" {
		return nil
	}

	telegramUser, err := telegramUserFromMessage(message)
	if err != nil {
		if errors.Is(err, auth.ErrUsernameRequired) {
			return err
		}

		return err
	}

	if err := c.notificationService.BindAdminChat(ctx, telegramUser); err != nil {
		return fmt.Errorf("bind admin chat: %w", err)
	}

	text := strings.TrimSpace(message.Text)
	if c.notificationService.IsAdminUsername(telegramUser.Username) && strings.HasPrefix(text, "/") {
		if handled, err := c.handleAdminCommand(ctx, telegramUser, text); handled {
			if err != nil {
				return err
			}

			return nil
		}
	}

	if !strings.HasPrefix(text, "/start") {
		return nil
	}

	user, created, err := c.authService.RegisterTelegramUser(ctx, telegramUser)
	if err != nil {
		return fmt.Errorf("register telegram user: %w", err)
	}

	if created {
		if err := c.notificationService.NotifyRegistration(ctx, &user); err == nil {
			_, _ = c.authService.SetUserStatusByTelegramID(ctx, user.TelegramID, model.UserStatusWaitingApprove)
		}
	}

	if err := c.authService.SendStartMessage(ctx, telegramUser.ChatID); err != nil {
		return fmt.Errorf("send start message: %w", err)
	}

	return nil
}

func (c *Controller) handleAdminCommand(ctx context.Context, adminUser model.TelegramUser, text string) (bool, error) {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return false, nil
	}

	switch parts[0] {
	case "/help":
		return true, c.notificationService.SendAdminText(ctx, adminUser.ChatID, "*Admin commands*\n\n`/users_approved`\n`/users_waiting`\n`/approve @username`\n`/revoke @username`\n`/help`")
	case "/users_approved":
		users, err := c.adminService.ListApprovedUsers(ctx)
		if err != nil {
			return true, fmt.Errorf("list approved users: %w", err)
		}

		return true, c.notificationService.SendAdminList(ctx, adminUser.ChatID, "Approved users", users)
	case "/users_waiting":
		users, err := c.adminService.ListWaitingUsers(ctx)
		if err != nil {
			return true, fmt.Errorf("list waiting users: %w", err)
		}

		return true, c.notificationService.SendAdminList(ctx, adminUser.ChatID, "Waiting approval", users)
	case "/approve":
		if len(parts) < adminCommandArgs {
			return true, c.notificationService.SendAdminText(ctx, adminUser.ChatID, "Usage: `/approve @username`")
		}

		username := normalizeAdminUsername(parts[1])
		if username == "" {
			return true, c.notificationService.SendAdminText(ctx, adminUser.ChatID, "Invalid username")
		}

		user, err := c.adminService.ApproveUser(ctx, username)
		if err != nil {
			return true, fmt.Errorf("approve user: %w", err)
		}

		return true, c.notificationService.SendAdminText(ctx, adminUser.ChatID, "*User approved*\n\n@"+user.Username)
	case "/revoke":
		if len(parts) < adminCommandArgs {
			return true, c.notificationService.SendAdminText(ctx, adminUser.ChatID, "Usage: `/revoke @username`")
		}

		username := normalizeAdminUsername(parts[1])
		if username == "" {
			return true, c.notificationService.SendAdminText(ctx, adminUser.ChatID, "Invalid username")
		}

		user, deletedTunnels, err := c.adminService.RevokeUser(ctx, username)
		if err != nil {
			return true, fmt.Errorf("revoke user: %w", err)
		}

		return true, c.notificationService.SendAdminText(ctx, adminUser.ChatID, fmt.Sprintf("*User revoked*\n\n@%s\nDeleted tunnels: `%d`\nNew status: `waiting_approve`", user.Username, deletedTunnels))
	default:
		return false, nil
	}
}
