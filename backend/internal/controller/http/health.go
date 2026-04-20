package http

import nethttp "net/http"

func (c *Controller) Health(w nethttp.ResponseWriter, _ *nethttp.Request) {
	_ = c

	writeJSON(w, nethttp.StatusOK, map[string]string{"status": "ok"})
}
