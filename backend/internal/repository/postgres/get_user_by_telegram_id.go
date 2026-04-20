package postgres

import (
	"context"
	"fmt"

	"wg-easy-app/backend/internal/model"
)

func (r *Repository) GetUserByTelegramID(ctx context.Context, telegramID int64) (model.User, error) {
	const query = `
		SELECT
			id,
			telegram_id,
			username,
			language_code,
			chat_id,
			status,
			created_at,
			updated_at
		FROM users
		WHERE telegram_id = ?`

	var user model.User

	err := r.conn.QueryRowContext(ctx, query, telegramID).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.LanguageCode,
		&user.ChatID,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return model.User{}, mapNotFound("user", fmt.Errorf("get user by telegram id: %w", err))
	}

	return user, nil
}
