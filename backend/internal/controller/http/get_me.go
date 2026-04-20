package http

import (
	nethttp "net/http"

	"wg-easy-app/backend/internal/middleware"
)

func (c *Controller) Me(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, ok := middleware.CurrentUser(r.Context())
	if !ok {
		writeError(w, nethttp.StatusUnauthorized, "unauthorized")

		return
	}

	usedTunnels := 0

	if user.IsApproved() {
		tunnels, err := c.tunnelService.ListByUserID(r.Context(), user)
		if err != nil {
			writeError(w, nethttp.StatusInternalServerError, "failed to list tunnels")

			return
		}

		usedTunnels = len(tunnels)
	}

	writeJSON(w, nethttp.StatusOK, map[string]any{
		"user":         user,
		"max_tunnels":  c.tunnelService.MaxTunnels(),
		"used_tunnels": usedTunnels,
	})
}
