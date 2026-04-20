package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-telegram/bot"
	_ "modernc.org/sqlite"
	"wg-easy-app/backend/internal/config"
	"wg-easy-app/backend/internal/controller"
	"wg-easy-app/backend/internal/migrations"
	postgresrepo "wg-easy-app/backend/internal/repository/postgres"
	telegramrepo "wg-easy-app/backend/internal/repository/telegram"
	wgeasyrepo "wg-easy-app/backend/internal/repository/wgeasy"
	authservice "wg-easy-app/backend/internal/service/auth"
	notificationservice "wg-easy-app/backend/internal/service/notification"
	tunnelservice "wg-easy-app/backend/internal/service/tunnel"
)

const (
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := sql.Open("sqlite", cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Printf("close sqlite connection: %v", err)
		}
	}()

	if _, err := dbConn.ExecContext(ctx, `PRAGMA foreign_keys = ON`); err != nil {
		log.Fatal(err)
	}

	if err := dbConn.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	if err := migrations.Run(ctx, dbConn); err != nil {
		log.Fatal(err)
	}

	botClient, err := bot.New(cfg.MainBotToken)
	if err != nil {
		log.Fatal(err)
	}

	wgRepo, err := wgeasyrepo.New(cfg.WGEasyBaseURL, cfg.WGEasyUsername, cfg.WGEasyPassword, cfg.WGEasyInsecureTLS)
	if err != nil {
		log.Fatal(err)
	}

	dbRepo := postgresrepo.NewRepository(dbConn)
	tgRepo := telegramrepo.New(botClient)

	authService := authservice.New(cfg, dbRepo, tgRepo)
	tunnelService := tunnelservice.New(cfg, dbRepo, tgRepo, wgRepo)
	notificationService := notificationservice.New(cfg, tgRepo)

	ctrl := controller.New(authService, tunnelService, notificationService)
	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	server := &http.Server{
		Addr:              addr,
		Handler:           controller.Static("/app/static", ctrl.Routes()),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()

		_ = server.Shutdown(shutdownCtx)
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
