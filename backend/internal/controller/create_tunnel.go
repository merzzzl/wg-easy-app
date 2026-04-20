package controller

import "net/http"

func (c *Controller) CreateTunnel(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")

		return
	}

	tunnelModel, err := c.tunnelService.Create(r.Context(), user)
	if err != nil {
		status, message := mapTunnelError(err)
		writeError(w, status, message)

		return
	}

	_ = c.notificationService.NotifyTunnelCreated(r.Context(), user, tunnelModel)

	writeJSON(w, http.StatusCreated, tunnelModel)
}
