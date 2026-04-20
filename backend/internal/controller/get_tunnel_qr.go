package controller

import "net/http"

func (c *Controller) GetTunnelQR(w http.ResponseWriter, r *http.Request) {
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

	qr, err := c.tunnelService.GetQRCodeSVG(r.Context(), user, tunnelID)
	if err != nil {
		status, message := mapTunnelError(err)
		writeError(w, status, message)

		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"svg": qr})
}
