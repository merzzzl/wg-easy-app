package controller

import "net/http"

func (c *Controller) DeleteTunnel(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")

		return
	}

	tunnelID, err := parseTunnelID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())

		return
	}

	tunnelModel, err := c.tunnelService.Delete(r.Context(), user, tunnelID)
	if err != nil {
		status, message := mapTunnelError(err)
		writeError(w, status, message)

		return
	}

	_ = c.notificationService.NotifyTunnelDeleted(r.Context(), user, tunnelModel)

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
