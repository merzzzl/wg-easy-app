package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"wg-easy-app/backend/internal/model"
)

func (s *Service) ValidateInitData(initData string) (model.TelegramUser, error) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return model.TelegramUser{}, fmt.Errorf("%w: parse query: %w", ErrInvalidInitData, err)
	}

	hash := values.Get("hash")
	if hash == "" {
		return model.TelegramUser{}, fmt.Errorf("%w: missing hash", ErrInvalidInitData)
	}

	dataCheckString := makeDataCheckString(values)

	secretKeyMac := hmac.New(sha256.New, []byte("WebAppData"))
	if _, err := secretKeyMac.Write([]byte(s.config.MainBotToken)); err != nil {
		return model.TelegramUser{}, fmt.Errorf("%w: write bot token: %w", ErrInvalidInitData, err)
	}

	secret := secretKeyMac.Sum(nil)

	mac := hmac.New(sha256.New, secret)
	if _, err := mac.Write([]byte(dataCheckString)); err != nil {
		return model.TelegramUser{}, fmt.Errorf("%w: write signature payload: %w", ErrInvalidInitData, err)
	}

	expectedHash := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expectedHash), []byte(hash)) {
		return model.TelegramUser{}, fmt.Errorf("%w: hash mismatch", ErrInvalidInitData)
	}

	userJSON := values.Get("user")
	if userJSON == "" {
		return model.TelegramUser{}, fmt.Errorf("%w: missing user payload", ErrInvalidInitData)
	}

	var payload struct {
		ID           int64  `json:"id"`
		Username     string `json:"username"`
		LanguageCode string `json:"language_code"`
	}

	if err := json.Unmarshal([]byte(userJSON), &payload); err != nil {
		return model.TelegramUser{}, fmt.Errorf("%w: decode user payload: %w", ErrInvalidInitData, err)
	}

	if payload.ID == 0 {
		return model.TelegramUser{}, fmt.Errorf("%w: missing telegram id", ErrInvalidInitData)
	}

	if payload.Username == "" {
		return model.TelegramUser{}, ErrUsernameRequired
	}

	return model.TelegramUser{
		TelegramID:   payload.ID,
		Username:     payload.Username,
		LanguageCode: payload.LanguageCode,
	}, nil
}

func makeDataCheckString(values url.Values) string {
	parts := make([]string, 0, len(values))
	for key, items := range values {
		if key == "hash" || len(items) == 0 {
			continue
		}

		parts = append(parts, key+"="+items[0])
	}

	sort.Strings(parts)

	return strings.Join(parts, "\n")
}
