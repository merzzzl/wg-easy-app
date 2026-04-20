package controller

import (
	"errors"
	"log/slog"
	"net/http"

	"wg-easy-app/backend/internal/service/auth"
)

func (c *Controller) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		initData := r.Header.Get("Tg-Token")
		if initData == "" {
			slog.Error("auth_failed", "path", r.URL.Path, "reason", "missing_tg_token", "remote", r.RemoteAddr)
			writeError(w, http.StatusUnauthorized, "missing tg-token header")

			return
		}

		user, created, err := c.authService.Authenticate(r.Context(), initData)
		if err != nil {
			status := http.StatusUnauthorized
			if errors.Is(err, auth.ErrUsernameRequired) {
				status = http.StatusBadRequest
			}

			slog.Error("auth_failed", "path", r.URL.Path, "status", status, "reason", err, "remote", r.RemoteAddr)

			writeError(w, status, err.Error())

			return
		}

		if created {
			_ = c.notificationService.NotifyRegistration(r.Context(), &user)
		}

		slog.Info("auth_ok", "path", r.URL.Path, "telegram_id", user.TelegramID, "username", user.Username, "remote", r.RemoteAddr)

		next.ServeHTTP(w, r.WithContext(withCurrentUser(r.Context(), &user)))
	})
}
