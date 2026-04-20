package controller

import (
	"errors"
	"net/http"

	"wg-easy-app/backend/internal/service/auth"
)

func (c *Controller) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		initData := r.Header.Get("Tg-Token")
		if initData == "" {
			writeError(w, http.StatusUnauthorized, "missing tg-token header")

			return
		}

		user, created, err := c.authService.Authenticate(r.Context(), initData)
		if err != nil {
			status := http.StatusUnauthorized
			if errors.Is(err, auth.ErrUsernameRequired) {
				status = http.StatusBadRequest
			}

			writeError(w, status, err.Error())

			return
		}

		if created {
			_ = c.notificationService.NotifyRegistration(r.Context(), &user)
		}

		next.ServeHTTP(w, r.WithContext(withCurrentUser(r.Context(), &user)))
	})
}
