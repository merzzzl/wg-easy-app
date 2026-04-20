package postgres

import (
	"context"
	"fmt"

	"wg-easy-app/backend/internal/model"
)

func (r *Repository) SetTunnelWGClientID(ctx context.Context, tunnelID int64, wgClientID string) (model.Tunnel, error) {
	const query = `
		UPDATE tunnels
		SET wg_client_id = ?
		WHERE id = ?
		RETURNING
			id,
			user_id,
			wg_client_name,
			wg_client_id,
			created_at`

	var tunnel model.Tunnel

	err := r.conn.QueryRowContext(ctx, query, wgClientID, tunnelID).Scan(
		&tunnel.ID,
		&tunnel.UserID,
		&tunnel.WGClientName,
		&tunnel.WGClientID,
		&tunnel.CreatedAt,
	)
	if err != nil {
		return model.Tunnel{}, mapNotFound("tunnel", fmt.Errorf("set tunnel wg client id: %w", err))
	}

	return tunnel, nil
}
