package postgres

import (
	"context"
	"fmt"

	"wg-easy-app/backend/internal/model"
)

func (r *Repository) CreateTunnel(ctx context.Context, params model.CreateTunnelParams) (model.Tunnel, error) {
	const query = `
		INSERT INTO tunnels (
			user_id,
			wg_client_name,
			wg_client_id,
			created_at
		) VALUES (?, '', '', CURRENT_TIMESTAMP)
		RETURNING
			id,
			user_id,
			wg_client_name,
			wg_client_id,
			created_at`

	var tunnel model.Tunnel

	err := r.conn.QueryRowContext(ctx, query, params.UserID).Scan(
		&tunnel.ID,
		&tunnel.UserID,
		&tunnel.WGClientName,
		&tunnel.WGClientID,
		&tunnel.CreatedAt,
	)
	if err != nil {
		return model.Tunnel{}, fmt.Errorf("create tunnel: %w", err)
	}

	return tunnel, nil
}
