package controller

import "net/http"

func (c *Controller) Health(w http.ResponseWriter, _ *http.Request) {
	_ = c

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
