package http

import (
	"encoding/json"
	"errors"
	"io/fs"
	nethttp "net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"wg-easy-app/backend/internal/service/notification"
	"wg-easy-app/backend/internal/service/tunnel"
)

type Controller struct {
	tunnelService       *tunnel.Service
	notificationService *notification.Service
}

var ErrInvalidTunnelID = errors.New("invalid tunnel id")

func New(tunnelService *tunnel.Service, notificationService *notification.Service) *Controller {
	return &Controller{
		tunnelService:       tunnelService,
		notificationService: notificationService,
	}
}

func (c *Controller) Routes(authMiddleware func(nethttp.Handler) nethttp.Handler) nethttp.Handler {
	mux := nethttp.NewServeMux()
	mux.HandleFunc("GET /healthz", c.Health)

	api := nethttp.NewServeMux()
	api.HandleFunc("GET /me", c.Me)
	api.HandleFunc("GET /tunnels", c.ListTunnels)
	api.HandleFunc("POST /tunnels", c.CreateTunnel)
	api.HandleFunc("DELETE /tunnels/{id}", c.DeleteTunnel)
	api.HandleFunc("GET /tunnels/{id}/qr", c.TunnelQR)
	api.HandleFunc("GET /tunnels/{id}/config", c.SendTunnelConfig)

	mux.Handle("/api/v1/", nethttp.StripPrefix("/api/v1", authMiddleware(api)))

	return mux
}

func Static(staticDir string, api nethttp.Handler) nethttp.Handler {
	if strings.TrimSpace(staticDir) == "" {
		return api
	}

	if _, err := os.Stat(filepath.Join(staticDir, "index.html")); err != nil {
		return api
	}

	fileServer := nethttp.FileServer(nethttp.Dir(staticDir))
	indexPath := filepath.Join(staticDir, "index.html")
	staticFS := os.DirFS(staticDir)

	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		requestPath := r.URL.Path
		if strings.HasPrefix(requestPath, "/api/") || strings.HasPrefix(requestPath, "/telegram/") || requestPath == "/healthz" {
			api.ServeHTTP(w, r)

			return
		}

		if requestPath == "/" {
			nethttp.ServeFile(w, r, indexPath)

			return
		}

		cleanPath := filepath.Clean(strings.TrimPrefix(requestPath, "/"))
		if !fs.ValidPath(cleanPath) {
			nethttp.ServeFile(w, r, indexPath)

			return
		}

		info, err := fs.Stat(staticFS, cleanPath)
		if err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)

			return
		}

		nethttp.ServeFile(w, r, indexPath)
	})
}

func writeJSON(w nethttp.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		nethttp.Error(w, nethttp.StatusText(nethttp.StatusInternalServerError), nethttp.StatusInternalServerError)
	}
}

func writeError(w nethttp.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func parseTunnelID(r *nethttp.Request) (int64, error) {
	value := r.PathValue("id")

	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, ErrInvalidTunnelID
	}

	return id, nil
}

func mapTunnelError(err error) (int, string) {
	switch {
	case errors.Is(err, tunnel.ErrTunnelNotFound):
		return nethttp.StatusNotFound, "tunnel not found"
	case errors.Is(err, tunnel.ErrTunnelLimitExceeded):
		return nethttp.StatusConflict, "tunnel limit exceeded"
	case errors.Is(err, tunnel.ErrUserNotApproved):
		return nethttp.StatusForbidden, "user is pending approval"
	default:
		return nethttp.StatusInternalServerError, "internal server error"
	}
}
