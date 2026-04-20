package postgres

import (
	"context"
	"fmt"
	"strconv"
)

func (r *Repository) SetAdminChatID(ctx context.Context, chatID int64) error {
	const query = `
		INSERT INTO app_settings (key, value, updated_at)
		VALUES ('admin_chat_id', ?, CURRENT_TIMESTAMP)
		ON CONFLICT(key) DO UPDATE SET
			value = excluded.value,
			updated_at = CURRENT_TIMESTAMP`

	if _, err := r.conn.ExecContext(ctx, query, strconv.FormatInt(chatID, 10)); err != nil {
		return fmt.Errorf("set admin chat id: %w", err)
	}

	return nil
}
