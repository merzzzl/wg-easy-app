package postgres

import (
	"context"
	"fmt"
)

func (r *Repository) DeleteTunnel(ctx context.Context, tunnelID int64) error {
	const query = `DELETE FROM tunnels WHERE id = ?`

	if _, err := r.conn.ExecContext(ctx, query, tunnelID); err != nil {
		return fmt.Errorf("delete tunnel: %w", err)
	}

	return nil
}
