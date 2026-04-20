package webhook

import (
	"errors"
	"strings"

	"github.com/go-telegram/bot/models"

	"wg-easy-app/backend/internal/model"
	"wg-easy-app/backend/internal/service/admin"
	"wg-easy-app/backend/internal/service/auth"
	"wg-easy-app/backend/internal/service/notification"
)

type Controller struct {
	authService         *auth.Service
	adminService        *admin.Service
	notificationService *notification.Service
}

var ErrTelegramUserEmpty = errors.New("telegram user not found")

func New(authService *auth.Service, adminService *admin.Service, notificationService *notification.Service) *Controller {
	return &Controller{
		authService:         authService,
		adminService:        adminService,
		notificationService: notificationService,
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
