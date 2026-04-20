package tunnel

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"wg-easy-app/backend/internal/config"
	"wg-easy-app/backend/internal/model"
	"wg-easy-app/backend/internal/repository/postgres"
	"wg-easy-app/backend/internal/repository/telegram"
	"wg-easy-app/backend/internal/repository/wgeasy"
)

var (
	ErrTunnelLimitExceeded = errors.New("tunnel: tunnel limit exceeded")
	ErrTunnelNotFound      = errors.New("tunnel: tunnel not found")
)

type Service struct {
	maxTunnels int
	db         *postgres.Repository
	tg         *telegram.Repository
	wg         *wgeasy.Repository
}

func New(cfg *config.Config, db *postgres.Repository, tg *telegram.Repository, wg *wgeasy.Repository) *Service {
	return &Service{
		maxTunnels: cfg.MaxTunnels,
		db:         db,
		tg:         tg,
		wg:         wg,
	}
}

func buildWGClientName(username string, tunnelID int64) string {
	name := strings.TrimSpace(username)
	if name == "" {
		return strconv.FormatInt(tunnelID, 10)
	}

	return fmt.Sprintf("%s-%d", name, tunnelID)
}

func findTunnelByID(tunnels []model.Tunnel, tunnelID int64) (model.Tunnel, error) {
	for _, tunnel := range tunnels {
		if tunnel.ID == tunnelID {
			return tunnel, nil
		}
	}

	return model.Tunnel{}, ErrTunnelNotFound
}
