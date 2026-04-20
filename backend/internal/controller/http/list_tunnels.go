package http

import (
	nethttp "net/http"

	"wg-easy-app/backend/internal/middleware"
)

func (c *Controller) ListTunnels(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, ok := middleware.CurrentUser(r.Context())
	if !ok {
		writeError(w, nethttp.StatusUnauthorized, "unauthorized")

		return
	}

	tunnels, err := c.tunnelService.ListByUserID(r.Context(), user)
	if err != nil {
		status, message := mapTunnelError(err)
		writeError(w, status, message)

		return
	}

	writeJSON(w, nethttp.StatusOK, tunnels)
}
