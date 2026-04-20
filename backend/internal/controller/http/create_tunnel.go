package http

import (
	"log/slog"
	nethttp "net/http"

	"wg-easy-app/backend/internal/middleware"
)

func (c *Controller) CreateTunnel(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, ok := middleware.CurrentUser(r.Context())
	if !ok {
		writeError(w, nethttp.StatusUnauthorized, "unauthorized")

		return
	}

	tunnelModel, err := c.tunnelService.Create(r.Context(), user)
	if err != nil {
		slog.Error("controller.create_tunnel failed", "user_id", user.ID, "error", err)

		status, message := mapTunnelError(err)
		writeError(w, status, message)

		return
	}

	_ = c.notificationService.NotifyTunnelCreated(r.Context(), user, tunnelModel)
	writeJSON(w, nethttp.StatusCreated, tunnelModel)
}
