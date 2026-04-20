package telegram

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (r *Repository) SendWebAppMessage(ctx context.Context, chatID any, text, buttonText, appURL string) error {
	_, err := r.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{{
				{
					Text: buttonText,
					WebApp: &models.WebAppInfo{
						URL: appURL,
					},
				},
			}},
		},
	})
	if err != nil {
		return fmt.Errorf("send web app message: %w", err)
	}

	return nil
}
