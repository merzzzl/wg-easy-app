package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"wg-easy-app/backend/internal/model"
)

func (r *Repository) UpsertUser(ctx context.Context, params model.UserUpsertParams) (model.User, bool, error) {
	const selectQuery = `
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

	var existing model.User

	err := r.conn.QueryRowContext(ctx, selectQuery, params.TelegramID).Scan(
		&existing.ID,
		&existing.TelegramID,
		&existing.Username,
		&existing.LanguageCode,
		&existing.ChatID,
		&existing.Status,
		&existing.CreatedAt,
		&existing.UpdatedAt,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return model.User{}, false, fmt.Errorf("select existing user: %w", err)
	}

	if errors.Is(err, sql.ErrNoRows) {
		const insertQuery = `
			INSERT INTO users (
				telegram_id,
				username,
				language_code,
				chat_id,
				status,
				created_at,
				updated_at
			) VALUES (?, ?, ?, ?, 'pending', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			RETURNING
				id,
				telegram_id,
				username,
				language_code,
				chat_id,
				status,
				created_at,
				updated_at`

		var createdUser model.User
		if err := r.conn.QueryRowContext(ctx, insertQuery,
			params.TelegramID,
			params.Username,
			params.LanguageCode,
			params.ChatID,
		).Scan(
			&createdUser.ID,
			&createdUser.TelegramID,
			&createdUser.Username,
			&createdUser.LanguageCode,
			&createdUser.ChatID,
			&createdUser.Status,
			&createdUser.CreatedAt,
			&createdUser.UpdatedAt,
		); err != nil {
			return model.User{}, false, fmt.Errorf("insert user: %w", err)
		}

		return createdUser, true, nil
	}

	const updateQuery = `
		UPDATE users
		SET
			username = ?,
			language_code = ?,
			chat_id = CASE WHEN ? = 0 THEN chat_id ELSE ? END,
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

	var updatedUser model.User
	if err := r.conn.QueryRowContext(ctx, updateQuery,
		params.Username,
		params.LanguageCode,
		params.ChatID,
		params.ChatID,
		params.TelegramID,
	).Scan(
		&updatedUser.ID,
		&updatedUser.TelegramID,
		&updatedUser.Username,
		&updatedUser.LanguageCode,
		&updatedUser.ChatID,
		&updatedUser.Status,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	); err != nil {
		return model.User{}, false, fmt.Errorf("update user: %w", err)
	}

	return updatedUser, false, nil
}
