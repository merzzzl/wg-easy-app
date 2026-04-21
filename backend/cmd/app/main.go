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
	"github.com/go-telegram/bot/models"
	_ "modernc.org/sqlite"

	"wg-easy-app/backend/internal/config"
	httpcontroller "wg-easy-app/backend/internal/controller/http"
	webhookcontroller "wg-easy-app/backend/internal/controller/webhook"
	"wg-easy-app/backend/internal/middleware"
	"wg-easy-app/backend/internal/migrations"
	postgresrepo "wg-easy-app/backend/internal/repository/postgres"
	telegramrepo "wg-easy-app/backend/internal/repository/telegram"
	wgeasyrepo "wg-easy-app/backend/internal/repository/wgeasy"
	adminservice "wg-easy-app/backend/internal/service/admin"
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

	var webhookController *webhookcontroller.Controller

	botClient, err := bot.New(
		cfg.MainBotToken,
		bot.WithAllowedUpdates(bot.AllowedUpdates{"message"}),
		bot.WithDefaultHandler(func(ctx context.Context, _ *bot.Bot, update *models.Update) {
			if webhookController == nil {
				log.Print("telegram update skipped: controller not initialized")

				return
			}

			if err := webhookController.HandleUpdate(ctx, update); err != nil {
				log.Printf("process telegram update: %v", err)
			}
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := botClient.DeleteWebhook(ctx, &bot.DeleteWebhookParams{
		DropPendingUpdates: false,
	}); err != nil {
		log.Fatalf("delete telegram webhook: %v", err)
	}

	log.Print("telegram long polling enabled")

	wgRepo, err := wgeasyrepo.New(cfg.WGEasyBaseURL, cfg.WGEasyUsername, cfg.WGEasyPassword, cfg.WGEasyInsecureTLS)
	if err != nil {
		log.Fatal(err)
	}

	dbRepo := postgresrepo.NewRepository(dbConn)
	tgRepo := telegramrepo.New(botClient)

	authService := authservice.New(cfg, dbRepo, tgRepo)
	adminService := adminservice.New(dbRepo, wgRepo)
	tunnelService := tunnelservice.New(cfg, dbRepo, tgRepo, wgRepo)
	notificationService := notificationservice.New(cfg, dbRepo, tgRepo)

	httpController := httpcontroller.New(tunnelService, notificationService)
	webhookController = webhookcontroller.New(authService, adminService, notificationService)

	mux := http.NewServeMux()
	mux.Handle("/", httpcontroller.Static("/app/static", httpController.Routes(middleware.Auth(authService, notificationService))))

	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	server := &http.Server{
		Addr:              addr,
		Handler:           middleware.RequestLogger(mux),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()

		_ = server.Shutdown(shutdownCtx)
	}()

	go botClient.Start(ctx)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
