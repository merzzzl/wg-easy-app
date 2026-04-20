package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-telegram/bot/models"
	"wg-easy-app/backend/internal/model"
	"wg-easy-app/backend/internal/service/auth"
	"wg-easy-app/backend/internal/service/notification"
	"wg-easy-app/backend/internal/service/tunnel"
)

type Controller struct {
	authService         *auth.Service
	tunnelService       *tunnel.Service
	notificationService *notification.Service
}

var (
	ErrInvalidTunnelID   = errors.New("invalid tunnel id")
	ErrTelegramUserEmpty = errors.New("telegram user not found")
)

func New(authService *auth.Service, tunnelService *tunnel.Service, notificationService *notification.Service) *Controller {
	return &Controller{
		authService:         authService,
		tunnelService:       tunnelService,
		notificationService: notificationService,
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

func parseTunnelID(r *http.Request) (int64, error) {
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
		return http.StatusNotFound, "tunnel not found"
	case errors.Is(err, tunnel.ErrTunnelLimitExceeded):
		return http.StatusConflict, "tunnel limit exceeded"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

func telegramUserFromMessage(message *models.Message) (model.TelegramUser, error) {
	if message == nil || message.From == nil {
		return model.TelegramUser{}, ErrTelegramUserEmpty
	}

	if strings.TrimSpace(message.From.Username) == "" {
		return model.TelegramUser{}, auth.ErrUsernameRequired
	}

	return model.TelegramUser{
		TelegramID:   message.From.ID,
		Username:     message.From.Username,
		LanguageCode: message.From.LanguageCode,
		ChatID:       message.Chat.ID,
	}, nil
}
