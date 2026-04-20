package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-telegram/bot/models"
	"wg-easy-app/backend/internal/service/auth"
)

func (c *Controller) TelegramWebhook(w http.ResponseWriter, r *http.Request) {
	var update models.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		writeError(w, http.StatusBadRequest, "invalid telegram update")

		return
	}

	message := update.Message
	if message == nil || !strings.HasPrefix(strings.TrimSpace(message.Text), "/start") {
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})

		return
	}

	telegramUser, err := telegramUserFromMessage(message)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, auth.ErrUsernameRequired) {
			status = http.StatusBadRequest
		}

		writeError(w, status, err.Error())

		return
	}

	user, created, err := c.authService.RegisterTelegramUser(r.Context(), telegramUser)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to register telegram user")

		return
	}

	if created {
		_ = c.notificationService.NotifyRegistration(r.Context(), &user)
	}

	if err := c.authService.SendStartMessage(r.Context(), telegramUser.ChatID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to send start message")

		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
