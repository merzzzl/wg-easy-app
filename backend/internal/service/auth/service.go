package auth

import (
	"errors"

	"wg-easy-app/backend/internal/config"
	"wg-easy-app/backend/internal/repository/postgres"
	"wg-easy-app/backend/internal/repository/telegram"
)

var (
	ErrInvalidInitData  = errors.New("auth: invalid init data")
	ErrUsernameRequired = errors.New("auth: username is required")
)

type Service struct {
	config *config.Config
	db     *postgres.Repository
	tg     *telegram.Repository
}

func New(cfg *config.Config, db *postgres.Repository, tg *telegram.Repository) *Service {
	return &Service{
		config: cfg,
		db:     db,
		tg:     tg,
	}
}
