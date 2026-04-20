package controller

import "net/http"

func (c *Controller) GetMe(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")

		return
	}

	usedTunnels := 0

	if user.IsApproved() {
		tunnels, err := c.tunnelService.ListByUserID(r.Context(), user)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to list tunnels")

			return
		}

		usedTunnels = len(tunnels)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user":         user,
		"max_tunnels":  c.tunnelService.MaxTunnels(),
		"used_tunnels": usedTunnels,
	})
}
