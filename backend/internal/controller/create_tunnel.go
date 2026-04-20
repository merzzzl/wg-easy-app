package controller

import (
	"log/slog"
	"net/http"
)

func (c *Controller) CreateTunnel(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")

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

	writeJSON(w, http.StatusCreated, tunnelModel)
}
