package wgeasy

import (
	"context"
	"fmt"

	"wg-easy-app/backend/internal/model"
)

func (r *Repository) CreateClient(ctx context.Context, params model.WGEasyCreateClientParams) (model.WGEasyCreateClientResponse, error) {
	var response model.WGEasyCreateClientResponse
	if err := r.doJSON(ctx, "POST", "/api/client", params, &response); err != nil {
		return model.WGEasyCreateClientResponse{}, fmt.Errorf("create wg-easy client: %w", err)
	}

	return response, nil
}
