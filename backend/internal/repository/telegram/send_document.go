package telegram

import (
	"bytes"
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (r *Repository) SendDocument(ctx context.Context, chatID any, fileName, caption string, content []byte) error {
	_, err := r.bot.SendDocument(ctx, &bot.SendDocumentParams{
		ChatID:  chatID,
		Caption: caption,
		Document: &models.InputFileUpload{
			Filename: fileName,
			Data:     bytes.NewReader(content),
		},
	})
	if err != nil {
		return fmt.Errorf("send document: %w", err)
	}

	return nil
}
