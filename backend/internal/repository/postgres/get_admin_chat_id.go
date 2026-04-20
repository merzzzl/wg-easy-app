package postgres

import (
	"context"
	"fmt"
	"strconv"
)

func (r *Repository) GetAdminChatID(ctx context.Context) (int64, error) {
	const query = `SELECT value FROM app_settings WHERE key = 'admin_chat_id'`

	var value string
	if err := r.conn.QueryRowContext(ctx, query).Scan(&value); err != nil {
		return 0, mapNotFound("admin_chat_id", err)
	}

	chatID, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse admin chat id: %w", err)
	}

	return chatID, nil
}
