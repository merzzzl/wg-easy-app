package controller

import "net/http"

func (c *Controller) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", c.Health)
	mux.HandleFunc("POST /telegram/webhook", c.TelegramWebhook)

	api := http.NewServeMux()
	api.HandleFunc("GET /me", c.GetMe)
	api.HandleFunc("GET /tunnels", c.ListTunnels)
	api.HandleFunc("POST /tunnels", c.CreateTunnel)
	api.HandleFunc("DELETE /tunnels/{id}", c.DeleteTunnel)
	api.HandleFunc("GET /tunnels/{id}/qr", c.GetTunnelQR)
	api.HandleFunc("GET /tunnels/{id}/config", c.SendTunnelConfig)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", c.AuthMiddleware(api)))

	return c.RequestLogger(mux)
}
