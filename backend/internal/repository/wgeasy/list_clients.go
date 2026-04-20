package wgeasy

import (
	"context"
	"fmt"

	"wg-easy-app/backend/internal/model"
)

func (r *Repository) ListClients(ctx context.Context) ([]model.WGEasyClient, error) {
	var items []model.WGEasyClient
	if err := r.doJSON(ctx, "GET", "/api/client", nil, &items); err != nil {
		return nil, fmt.Errorf("list wg-easy clients: %w", err)
	}

	return items, nil
}
