package http

import (
	nethttp "net/http"

	"wg-easy-app/backend/internal/middleware"
)

func (c *Controller) DeleteTunnel(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, ok := middleware.CurrentUser(r.Context())
	if !ok {
		writeError(w, nethttp.StatusUnauthorized, "unauthorized")

		return
	}

	tunnelID, err := parseTunnelID(r)
	if err != nil {
		writeError(w, nethttp.StatusBadRequest, err.Error())

		return
	}

	tunnelModel, err := c.tunnelService.Delete(r.Context(), user, tunnelID)
	if err != nil {
		status, message := mapTunnelError(err)
		writeError(w, status, message)

		return
	}

	_ = c.notificationService.NotifyTunnelDeleted(r.Context(), user, tunnelModel)

	writeJSON(w, nethttp.StatusOK, map[string]bool{"ok": true})
}
