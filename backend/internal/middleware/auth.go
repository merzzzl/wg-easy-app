package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	authservice "wg-easy-app/backend/internal/service/auth"
	"wg-easy-app/backend/internal/service/notification"
)

func Auth(authService *authservice.Service, notificationService *notification.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			initData := r.Header.Get("Tg-Token")
			if initData == "" {
				slog.Error("auth_failed", "path", r.URL.Path, "reason", "missing_tg_token", "remote", r.RemoteAddr)
				writeError(w, http.StatusUnauthorized, "missing tg-token header")

				return
			}

			user, created, err := authService.Authenticate(r.Context(), initData)
			if err != nil {
				status := http.StatusUnauthorized
				if errors.Is(err, authservice.ErrUsernameRequired) {
					status = http.StatusBadRequest
				}

				slog.Error("auth_failed", "path", r.URL.Path, "status", status, "reason", err, "remote", r.RemoteAddr)
				writeError(w, status, err.Error())

				return
			}

			if created {
				_ = notificationService.NotifyRegistration(r.Context(), &user)
			}

			slog.Info("auth_ok", "path", r.URL.Path, "telegram_id", user.TelegramID, "username", user.Username, "remote", r.RemoteAddr)

			next.ServeHTTP(w, r.WithContext(WithCurrentUser(r.Context(), &user)))
		})
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte("{\"error\":\"" + message + "\"}"))
}
