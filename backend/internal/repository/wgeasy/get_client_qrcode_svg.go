package wgeasy

import (
	"context"
	"fmt"
	"net/url"
)

func (r *Repository) GetClientQRCodeSVG(ctx context.Context, clientID string) (string, error) {
	body, err := r.doRaw(ctx, "GET", "/api/client/"+url.PathEscape(clientID)+"/qrcode.svg", nil)
	if err != nil {
		return "", fmt.Errorf("get client qrcode svg: %w", err)
	}

	return string(body), nil
}
