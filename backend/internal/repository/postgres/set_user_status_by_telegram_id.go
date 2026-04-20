package postgres

import (
	"context"
	"fmt"

	"wg-easy-app/backend/internal/model"
)

func (r *Repository) SetUserStatusByTelegramID(ctx context.Context, telegramID int64, status model.UserStatus) (model.User, error) {
	const query = `
		UPDATE users
		SET
			status = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE telegram_id = ?
		RETURNING
			id,
			telegram_id,
			username,
			language_code,
			chat_id,
			status,
			created_at,
			updated_at`

	var user model.User
	if err := r.conn.QueryRowContext(ctx, query, status, telegramID).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.LanguageCode,
		&user.ChatID,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return model.User{}, mapNotFound("user", fmt.Errorf("set user status by telegram id: %w", err))
	}

	return user, nil
}
