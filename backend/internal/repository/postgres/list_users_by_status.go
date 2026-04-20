package postgres

import (
	"context"
	"fmt"

	"wg-easy-app/backend/internal/model"
)

func (r *Repository) ListUsersByStatus(ctx context.Context, status model.UserStatus) ([]model.User, error) {
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
		WHERE status = ?
		ORDER BY id DESC`

	rows, err := r.conn.QueryContext(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("list users by status: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	items := make([]model.User, 0)

	for rows.Next() {
		var user model.User

		if err := rows.Scan(
			&user.ID,
			&user.TelegramID,
			&user.Username,
			&user.LanguageCode,
			&user.ChatID,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan user by status: %w", err)
		}

		items = append(items, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users by status: %w", err)
	}

	return items, nil
}
