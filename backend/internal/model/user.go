package model

import "time"

type UserStatus string

const (
	UserStatusPending  UserStatus = "pending"
	UserStatusApproved UserStatus = "approved"
)

type User struct {
	ID           int64      `json:"id"`
	TelegramID   int64      `json:"telegram_id"`
	Username     string     `json:"username"`
	LanguageCode string     `json:"language_code"`
	ChatID       int64      `json:"chat_id"`
	Status       UserStatus `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type UserUpsertParams struct {
	TelegramID   int64
	Username     string
	LanguageCode string
	ChatID       int64
}

func (u *User) IsApproved() bool {
	return u.Status == UserStatusApproved
}
