package controller

import (
	"errors"
	"log"
	"net/http"

	"wg-easy-app/backend/internal/service/auth"
)

func (c *Controller) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		initData := r.Header.Get("Tg-Token")
		if initData == "" {
			log.Printf("auth_failed path=%s reason=missing_tg_token remote=%s", r.URL.Path, r.RemoteAddr)
			writeError(w, http.StatusUnauthorized, "missing tg-token header")

			return
		}

		user, created, err := c.authService.Authenticate(r.Context(), initData)
		if err != nil {
			status := http.StatusUnauthorized
			if errors.Is(err, auth.ErrUsernameRequired) {
				status = http.StatusBadRequest
			}

			log.Printf("auth_failed path=%s status=%d reason=%v remote=%s", r.URL.Path, status, err, r.RemoteAddr)

			writeError(w, status, err.Error())

			return
		}

		if created {
			_ = c.notificationService.NotifyRegistration(r.Context(), &user)
		}

		log.Printf("auth_ok path=%s telegram_id=%d username=%s remote=%s", r.URL.Path, user.TelegramID, user.Username, r.RemoteAddr)

		next.ServeHTTP(w, r.WithContext(withCurrentUser(r.Context(), &user)))
	})
}
