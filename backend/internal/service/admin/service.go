package admin

import (
	"wg-easy-app/backend/internal/repository/postgres"
	"wg-easy-app/backend/internal/repository/wgeasy"
)

type Service struct {
	db *postgres.Repository
	wg *wgeasy.Repository
}

func New(db *postgres.Repository, wg *wgeasy.Repository) *Service {
	return &Service{db: db, wg: wg}
}
