package postgres

import (
	"context"
	"fmt"

	"wg-easy-app/backend/internal/model"
)

func (r *Repository) ListTunnelsByUserID(ctx context.Context, userID int64) ([]model.Tunnel, error) {
	const query = `
		SELECT
			id,
			user_id,
			wg_client_name,
			wg_client_id,
			created_at
		FROM tunnels
		WHERE user_id = ?
		ORDER BY id DESC`

	rows, err := r.conn.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list tunnels: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	items := make([]model.Tunnel, 0)

	for rows.Next() {
		var tunnel model.Tunnel

		if err := rows.Scan(
			&tunnel.ID,
			&tunnel.UserID,
			&tunnel.WGClientName,
			&tunnel.WGClientID,
			&tunnel.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan tunnel: %w", err)
		}

		items = append(items, tunnel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tunnels: %w", err)
	}

	return items, nil
}
