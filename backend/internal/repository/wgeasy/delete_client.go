package wgeasy

import (
	"context"
	"fmt"
	"net/url"
)

func (r *Repository) DeleteClient(ctx context.Context, clientID string) error {
	if err := r.doJSON(ctx, "DELETE", "/api/client/"+url.PathEscape(clientID), nil, nil); err != nil {
		return fmt.Errorf("delete wg-easy client: %w", err)
	}

	return nil
}
