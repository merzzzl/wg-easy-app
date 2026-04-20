package controller

import "net/http"

func (c *Controller) ListTunnels(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")

		return
	}

	tunnels, err := c.tunnelService.ListByUserID(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list tunnels")

		return
	}

	writeJSON(w, http.StatusOK, tunnels)
}
