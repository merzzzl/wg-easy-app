package telegram

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
)

func (r *Repository) SendMessage(ctx context.Context, chatID any, text string) error {
	if _, err := r.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (r *Repository) SendMarkdownMessage(ctx context.Context, chatID any, text string) error {
	if _, err := r.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "MarkdownV2",
	}); err != nil {
		return fmt.Errorf("send markdown message: %w", err)
	}

	return nil
}
