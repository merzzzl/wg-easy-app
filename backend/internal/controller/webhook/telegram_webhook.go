package webhook

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-telegram/bot/models"

	"wg-easy-app/backend/internal/model"
	"wg-easy-app/backend/internal/service/auth"
)

const adminCommandArgs = 2

func (c *Controller) TelegramWebhook(w http.ResponseWriter, r *http.Request) {
	var update models.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		writeError(w, http.StatusBadRequest, "invalid telegram update")

		return
	}

	message := update.Message
	if message == nil || strings.TrimSpace(message.Text) == "" {
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

	if err := c.notificationService.BindAdminChat(r.Context(), telegramUser); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to bind admin chat")

		return
	}

	text := strings.TrimSpace(message.Text)
	if c.notificationService.IsAdminUsername(telegramUser.Username) && strings.HasPrefix(text, "/") {
		if handled, err := c.handleAdminCommand(r, telegramUser, text); handled {
			if err != nil {
				writeError(w, http.StatusInternalServerError, err.Error())

				return
			}

			writeJSON(w, http.StatusOK, map[string]bool{"ok": true})

			return
		}
	}

	if !strings.HasPrefix(text, "/start") {
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})

		return
	}

	user, created, err := c.authService.RegisterTelegramUser(r.Context(), telegramUser)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to register telegram user")

		return
	}

	if created {
		if err := c.notificationService.NotifyRegistration(r.Context(), &user); err == nil {
			_, _ = c.authService.SetUserStatusByTelegramID(r.Context(), user.TelegramID, model.UserStatusWaitingApprove)
		}
	}

	if err := c.authService.SendStartMessage(r.Context(), telegramUser.ChatID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to send start message")

		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (c *Controller) handleAdminCommand(r *http.Request, adminUser model.TelegramUser, text string) (bool, error) {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return false, nil
	}

	switch parts[0] {
	case "/users_approved":
		users, err := c.adminService.ListApprovedUsers(r.Context())
		if err != nil {
			return true, fmt.Errorf("list approved users: %w", err)
		}

		return true, c.notificationService.SendAdminList(r.Context(), adminUser.ChatID, "Approved users", users)
	case "/users_waiting":
		users, err := c.adminService.ListWaitingUsers(r.Context())
		if err != nil {
			return true, fmt.Errorf("list waiting users: %w", err)
		}

		return true, c.notificationService.SendAdminList(r.Context(), adminUser.ChatID, "Waiting approval", users)
	case "/approve":
		if len(parts) < adminCommandArgs {
			return true, c.notificationService.SendAdminText(r.Context(), adminUser.ChatID, "Usage: /approve <telegram_id>")
		}

		telegramID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return true, c.notificationService.SendAdminText(r.Context(), adminUser.ChatID, "Invalid telegram_id")
		}

		user, err := c.adminService.ApproveUser(r.Context(), telegramID)
		if err != nil {
			return true, fmt.Errorf("approve user: %w", err)
		}

		return true, c.notificationService.SendAdminText(r.Context(), adminUser.ChatID, "Approved @"+user.Username)
	case "/revoke":
		if len(parts) < adminCommandArgs {
			return true, c.notificationService.SendAdminText(r.Context(), adminUser.ChatID, "Usage: /revoke <telegram_id>")
		}

		telegramID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return true, c.notificationService.SendAdminText(r.Context(), adminUser.ChatID, "Invalid telegram_id")
		}

		user, deletedTunnels, err := c.adminService.RevokeUser(r.Context(), telegramID)
		if err != nil {
			return true, fmt.Errorf("revoke user: %w", err)
		}

		return true, c.notificationService.SendAdminText(r.Context(), adminUser.ChatID, fmt.Sprintf("Revoked @%s and reset to waiting_approve. Deleted tunnels: %d", user.Username, deletedTunnels))
	default:
		return false, nil
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
