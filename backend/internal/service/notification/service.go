package notification

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"wg-easy-app/backend/internal/config"
	"wg-easy-app/backend/internal/model"
	"wg-easy-app/backend/internal/repository/postgres"
	"wg-easy-app/backend/internal/repository/telegram"
)

type Service struct {
	adminUsername string
	db            *postgres.Repository
	tg            *telegram.Repository
}

func New(cfg *config.Config, db *postgres.Repository, tg *telegram.Repository) *Service {
	return &Service{
		adminUsername: cfg.AdminUsername,
		db:            db,
		tg:            tg,
	}
}

func normalizeUsername(username string) string {
	return strings.TrimPrefix(strings.TrimSpace(strings.ToLower(username)), "@")
}

func (s *Service) IsAdminUsername(username string) bool {
	return normalizeUsername(username) == normalizeUsername(s.adminUsername)
}

func (s *Service) BindAdminChat(ctx context.Context, telegramUser model.TelegramUser) error {
	if normalizeUsername(telegramUser.Username) != normalizeUsername(s.adminUsername) {
		return nil
	}

	slog.Info("notification.bind_admin_chat called", "telegram_id", telegramUser.TelegramID, "chat_id", telegramUser.ChatID)

	return s.db.SetAdminChatID(ctx, telegramUser.ChatID)
}

func (s *Service) SendAdminText(ctx context.Context, chatID int64, text string) error {
	return s.tg.SendMarkdownMessage(ctx, chatID, text)
}

func (s *Service) SendAdminList(ctx context.Context, chatID int64, title string, users []model.User) error {
	if len(users) == 0 {
		return s.tg.SendMarkdownMessage(ctx, chatID, markdownTitle(title)+"\n\n_no users found_")
	}

	var builder strings.Builder

	_, _ = builder.WriteString(markdownTitle(title+"\n"))

	for _, user := range users {
		_, _ = fmt.Fprintf(&builder, "\n@%s", escapeMarkdown(user.Username))
	}

	return s.tg.SendMarkdownMessage(ctx, chatID, builder.String())
}

func (s *Service) sendAdminMessage(ctx context.Context, text string) error {
	chatID, err := s.db.GetAdminChatID(ctx)
	if err != nil {
		return fmt.Errorf("get admin chat id: %w", err)
	}

	return s.tg.SendMarkdownMessage(ctx, chatID, text)
}

func formatTunnelLine(wgClientName string) string {
	if wgClientName == "" {
		return "-"
	}

	return wgClientName
}

func actionText(action, username string, tunnelID int64, wgClientName string) string {
	message := fmt.Sprintf("%s\n\nUser: %s\nUsername: @%s", markdownTitle(action), escapeMarkdown(username), escapeMarkdown(username))
	if tunnelID == 0 {
		return message
	}

	return fmt.Sprintf("%s\nTunnel ID: `%d`\nTunnel name: `%s`", message, tunnelID, escapeMarkdown(formatTunnelLine(wgClientName)))
}

func registrationText(username, status string) string {
	return fmt.Sprintf("%s\n\nUser: @%s\nStatus: `%s`", markdownTitle("New registration request"), escapeMarkdown(username), escapeMarkdown(status))
}

func markdownTitle(title string) string {
	return "*" + escapeMarkdown(title) + "*"
}

func escapeMarkdown(value string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)

	return replacer.Replace(value)
}
