package wgeasy

import (
	"context"
	"fmt"
	"net/url"
)

func (r *Repository) GetClientConfiguration(ctx context.Context, clientID string) (string, error) {
	body, err := r.doRaw(ctx, "GET", "/api/client/"+url.PathEscape(clientID)+"/configuration", nil)
	if err != nil {
		return "", fmt.Errorf("get client configuration: %w", err)
	}

	return string(body), nil
}
