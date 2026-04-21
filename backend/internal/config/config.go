package config

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Host              string `env:"APP_HOST"                      envDefault:"0.0.0.0"`
	Port              int    `env:"APP_PORT"                      envDefault:"8080"`
	MainBotToken      string `env:"APP_MAIN_BOT_TOKEN,required"`
	AdminUsername     string `env:"APP_ADMIN_USERNAME,required"`
	MaxTunnels        int    `env:"APP_MAX_TUNNELS"               envDefault:"10"`
	DBURL             string `env:"APP_DB_URL,required"`
	WGEasyBaseURL     string `env:"APP_WG_EASY_BASE_URL,required"`
	WGEasyUsername    string `env:"APP_WG_EASY_USERNAME,required"`
	WGEasyPassword    string `env:"APP_WG_EASY_PASSWORD,required"`
	WGEasyInsecureTLS bool   `env:"APP_WG_EASY_INSECURE_TLS"      envDefault:"false"`
	MiniAppURL        string `env:"APP_MINI_APP_URL"              envDefault:""`
}

var (
	ErrMaxTunnelsInvalid = errors.New("APP_MAX_TUNNELS must be greater than zero")
	ErrPortInvalid       = errors.New("APP_PORT must be greater than zero")
	ErrAdminUsername     = errors.New("APP_ADMIN_USERNAME is required")
)

func Read() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}

	if cfg.MaxTunnels <= 0 {
		return nil, ErrMaxTunnelsInvalid
	}

	if cfg.Port <= 0 {
		return nil, ErrPortInvalid
	}

	if _, err := url.ParseRequestURI(cfg.WGEasyBaseURL); err != nil {
		return nil, fmt.Errorf("parse APP_WG_EASY_BASE_URL: %w", err)
	}

	if cfg.MiniAppURL != "" {
		if _, err := url.ParseRequestURI(cfg.MiniAppURL); err != nil {
			return nil, fmt.Errorf("parse APP_MINI_APP_URL: %w", err)
		}

		cfg.MiniAppURL = strings.TrimRight(cfg.MiniAppURL, "/")
	}

	if strings.TrimSpace(cfg.AdminUsername) == "" {
		return nil, ErrAdminUsername
	}

	if !strings.HasPrefix(cfg.AdminUsername, "@") {
		cfg.AdminUsername = "@" + cfg.AdminUsername
	}

	cfg.WGEasyBaseURL = strings.TrimRight(cfg.WGEasyBaseURL, "/")
	cfg.Host = strings.TrimSpace(cfg.Host)

	return &cfg, nil
}
