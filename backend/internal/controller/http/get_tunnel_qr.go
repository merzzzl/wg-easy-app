package http

import (
	nethttp "net/http"

	"wg-easy-app/backend/internal/middleware"
)

func (c *Controller) TunnelQR(w nethttp.ResponseWriter, r *nethttp.Request) {
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

	qr, err := c.tunnelService.GetQRCodeSVG(r.Context(), user, tunnelID)
	if err != nil {
		status, message := mapTunnelError(err)
		writeError(w, status, message)

		return
	}

	writeJSON(w, nethttp.StatusOK, map[string]string{"svg": qr})
}
