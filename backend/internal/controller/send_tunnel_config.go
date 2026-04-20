package controller

import "net/http"

func (c *Controller) SendTunnelConfig(w http.ResponseWriter, r *http.Request) {
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

	if err := c.tunnelService.SendConfig(r.Context(), user, tunnelID); err != nil {
		status, message := mapTunnelError(err)
		writeError(w, status, message)

		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
